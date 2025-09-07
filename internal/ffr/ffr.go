package ffr

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// バージョン情報（リリース時にldフラグで設定される）
var (
	version   = "dev"
	buildDate = "unknown"
)

func printUsage() {
	fmt.Println("Laravel Form Facade から HTMLタグ置換スクリプト")
	fmt.Println("使用方法: go run form_facade_replacer.go <ファイルパス|ディレクトリパス>")
	fmt.Println()
	fmt.Println("引数:")
	fmt.Println(" ファイルパス 対象の.blade.phpファイル")
	fmt.Println(" ディレクトリパス 対象ディレクトリ（配下の.blade.phpファイルを再帰処理）")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println(" -h, --help このヘルプメッセージを表示")
	fmt.Println(" -v, --version バージョン情報を表示")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println(" go run form_facade_replacer.go resources/views/hoge")
	fmt.Println(" go run form_facade_replacer.go resources/views/hoge/fuga.blade.php")
}

func printVersion() {
	fmt.Printf("Form Facade Replacer %s\n", version)
	fmt.Printf("Build Date: %s\n", buildDate)
	fmt.Println()
	fmt.Println("Laravel Form Facade を HTML に変換する高性能 Go ツール")
	fmt.Println("https://github.com/ryohirano/form-facade-replacer")
}

// DynamicAttributePair 動的属性のキーと値のペア
type DynamicAttributePair struct {
	Key   string // 動的キー（例: $condition ? 'disabled' : ''）
	Value string // 動的値（例: $condition ? 'disabled' : null）
}

// detectDynamicAttributes 文字列から動的属性を検出し、DynamicAttributePairの配列を返す
func detectDynamicAttributes(attrs string) []DynamicAttributePair {
	return extractDynamicAttributesBalanced(attrs)
}

// extractDynamicAttributesBalanced バランスした括弧を考慮して動的属性を抽出
func extractDynamicAttributesBalanced(attrs string) []DynamicAttributePair {
	var pairs []DynamicAttributePair
	var current strings.Builder
	var parenCount, bracketCount, braceCount int
	var inQuotes bool
	var quoteChar rune
	var escapeNext bool

	i := 0
	for i < len(attrs) {
		char := rune(attrs[i])

		if escapeNext {
			current.WriteRune(char)
			escapeNext = false
			i++
			continue
		}

		if char == '\\' && inQuotes {
			current.WriteRune(char)
			escapeNext = true
			i++
			continue
		}

		if !inQuotes && (char == '"' || char == '\'') {
			inQuotes = true
			quoteChar = char
		} else if inQuotes && char == quoteChar {
			inQuotes = false
			quoteChar = 0
		}

		if !inQuotes {
			switch char {
			case '(':
				parenCount++
			case ')':
				parenCount--
			case '[':
				bracketCount++
			case ']':
				bracketCount--
			case '{':
				braceCount++
			case '}':
				braceCount--
			case ',':
				if parenCount == 0 && bracketCount == 0 && braceCount == 0 {
					// 区切り文字に到達
					attributeStr := strings.TrimSpace(current.String())
					if isDynamicAttribute(attributeStr) {
						pair := parseDynamicAttributePair(attributeStr)
						if pair.Key != "" && pair.Value != "" {
							pairs = append(pairs, pair)
						}
					}
					current.Reset()
					i++
					continue
				}
			}
		}

		current.WriteRune(char)
		i++
	}

	// 最後の属性を処理
	if current.Len() > 0 {
		attributeStr := strings.TrimSpace(current.String())
		if isDynamicAttribute(attributeStr) {
			pair := parseDynamicAttributePair(attributeStr)
			if pair.Key != "" && pair.Value != "" {
				pairs = append(pairs, pair)
			}
		}
	}

	return pairs
}

// parseDynamicAttributePair 動的属性のキーと値を抽出する関数
func parseDynamicAttributePair(input string) DynamicAttributePair {
	patterns := []string{
		// 1. 標準的なパターン: $変数 ? 'キー' : 'キー' => 値
		`(^\$\w+(?:\[[^\]]*\])*(?:->[a-zA-Z_]\w*\([^)]*\))?\s*\?\s*'[^']*'\s*:\s*'[^']*')\s*=>\s*(.+)`,
		// 2. 複雑な条件: (条件) ? 'キー' : 'キー' => 値
		`(^\(.*?\)\s*\?\s*'[^']*'\s*:\s*'[^']*')\s*=>\s*(.+)`,
		// 3. ネストした三項演算子: $変数 ? (ネスト) : 'キー' => 値
		`(^\$\w+(?:\[[^\]]*\])*\s*\?\s*\([^)]+\)\s*:\s*'[^']*')\s*=>\s*(.+)`,
		// 4. 複雑な条件式（&&、||演算子含む）: $変数->method() && $変数->method() ? 'キー' : 'キー' => 値
		`(^\$\w+(?:\[[^\]]*\])*(?:->[a-zA-Z_]\w*\([^)]*\))?(?:\s*&&\s*\$\w+(?:\[[^\]]*\])*(?:->[a-zA-Z_]\w*\([^)]*\))?)*\s*\?\s*'[^']*'\s*:\s*'[^']*')\s*=>\s*(.+)`,
		// 5. 複雑な条件式（||演算子含む）: $変数->method() || $変数->method() ? 'キー' : 'キー' => 値
		`(^\$\w+(?:\[[^\]]*\])*(?:->[a-zA-Z_]\w*\([^)]*\))?(?:\s*\|\|\s*\$\w+(?:\[[^\]]*\])*(?:->[a-zA-Z_]\w*\([^)]*\))?)*\s*\?\s*'[^']*'\s*:\s*'[^']*')\s*=>\s*(.+)`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		matches := re.FindStringSubmatch(input)

		if len(matches) >= 3 {
			key := strings.TrimSpace(matches[1])
			value := strings.TrimSpace(matches[2])
			return DynamicAttributePair{
				Key:   key,
				Value: value,
			}
		}
	}

	return DynamicAttributePair{}
}

// isDynamicAttribute 動的属性パターンを検出する関数
func isDynamicAttribute(input string) bool {
	// より包括的な動的属性パターンを検出する正規表現
	// 1. 単純な変数: $変数 ? 'キー' : 'キー' => 値
	// 2. 複雑な条件: (条件) ? 'キー' : 'キー' => 値
	// 3. ネストした三項演算子: $変数 ? (ネスト) : 'キー' => 値
	// 4. 複雑な条件式（&&、||演算子含む）
	patterns := []string{
		`^\$\w+(?:\[[^\]]*\])*(?:->[a-zA-Z_]\w*\([^)]*\))?\s*\?\s*`, // $変数 ?
		`^\(.*?\)\s*\?\s*`, // (条件) ?
		`^\$\w+(?:\[[^\]]*\])*\s*\?\s*\([^)]+\)\s*:\s*'[^']*'\s*=>`, // $変数 ? (ネスト) : 'キー' =>
		`^\$\w+(?:\[[^\]]*\])*\s*\?\s*'[^']*'\s*:\s*'[^']*'\s*=>`,   // $変数 ? 'キー' : 'キー' =>
		`^\$\w+(?:\[[^\]]*\])*(?:->[a-zA-Z_]\w*\([^)]*\))?(?:\s*&&\s*\$\w+(?:\[[^\]]*\])*(?:->[a-zA-Z_]\w*\([^)]*\))?)*\s*\?\s*'[^']*'\s*:\s*'[^']*'\s*=>`,   // $変数 && $変数 ?
		`^\$\w+(?:\[[^\]]*\])*(?:->[a-zA-Z_]\w*\([^)]*\))?(?:\s*\|\|\s*\$\w+(?:\[[^\]]*\])*(?:->[a-zA-Z_]\w*\([^)]*\))?)*\s*\?\s*'[^']*'\s*:\s*'[^']*'\s*=>`, // $変数 || $変数 ?
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		if re.MatchString(input) {
			return true
		}
	}

	return false
}

// processDynamicAttributes 動的属性を処理してHTML属性文字列を生成
func processDynamicAttributes(attrs string) string {
	var result strings.Builder

	pairs := detectDynamicAttributes(attrs)
	for _, pair := range pairs {
		// 値の特別処理
		value := pair.Value

		// null、true、false などのリテラル値は {{ }} で囲まない
		if value == "null" || value == "true" || value == "false" {
			result.WriteString(fmt.Sprintf(` {{ %s }}="%s"`, pair.Key, value))
		} else {
			// 変数や複雑な式は {{ }} で囲む
			result.WriteString(fmt.Sprintf(` {{ %s }}="{{ %s }}"`, pair.Key, value))
		}
	}

	return result.String()
}

func replaceFormPatterns(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	text := string(content)
	text = replaceFormOpen(text)
	text = replaceFormClose(text)
	text = replaceFormHidden(text)
	text = replaceFormButton(text)
	text = replaceFormTextarea(text)
	text = replaceFormLabel(text)
	text = replaceFormText(text)
	text = replaceFormInput(text)
	text = replaceFormNumber(text)
	text = replaceFormSelect(text)
	text = replaceFormCheckbox(text)
	text = replaceFormSubmit(text)
	text = replaceFormFile(text)
	text = replaceFormEmail(text)
	text = replaceFormPassword(text)
	text = replaceFormUrl(text)
	text = replaceFormTel(text)
	text = replaceFormSearch(text)
	text = replaceFormDate(text)
	text = replaceFormTime(text)
	text = replaceFormDatetime(text)
	text = replaceFormRange(text)
	text = replaceFormColor(text)
	text = replaceFormRadio(text)

	return os.WriteFile(filePath, []byte(text), 0644)
}

func replaceFormHidden(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::hidden\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::hidden\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			content := re.FindStringSubmatch(match)[1]
			params := extractParamsAdvanced(content)
			return processFormHidden(params)
		})
	}
	return text
}

func processFormHidden(params []string) string {
	if len(params) < 1 {
		return ""
	}

	// PHP文字列連結を含むフィールド名を適切に処理
	nameAttr := ProcessFieldName(params[0])
	value := ""
	if len(params) > 1 {
		value = params[1]
	}

	// 属性処理の統一
	attrProcessor := &AttributeProcessor{
		Order: []string{"id", "class"},
		Patterns: map[string]string{
			"id":    `'id'\s*=>\s*'([^']+)'`,
			"class": `'class'\s*=>\s*'([^']+)'`,
		},
	}

	extraAttrs := ""
	if len(params) > 2 {
		extraAttrs = attrProcessor.ProcessAttributes(params[2])
	}

	// hidden input専用の値フォーマット
	formattedValue := FormatHiddenValueAttribute(value, nameAttr)
	return fmt.Sprintf(`<input type="hidden" name="%s" value="%s"%s>`, nameAttr, formattedValue, extraAttrs)
}

func replaceFormButton(text string) string {
	singleParamPatterns := []string{
		`(?s)\{\{\s*Form::button\(\s*'([^']*)'\s*\)\s*\}\}`,
		`(?s)\{\!\!\s*Form::button\(\s*'([^']*)'\s*\)\s*\!\!\}`,
	}

	for _, pattern := range singleParamPatterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			matches := re.FindStringSubmatch(match)
			return processFormButton(matches[1], "")
		})
	}

	twoParamPatterns := []string{
		`(?s)\{\!\!\s*Form::button\(\s*(.*?)\s*,\s*\[\s*(.*?)\s*\]\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::button\(\s*(.*?)\s*,\s*\[\s*(.*?)\s*\]\s*\)\s*\}\}`,
	}

	for _, pattern := range twoParamPatterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			matches := re.FindStringSubmatch(match)
			return processFormButton(matches[1], matches[2])
		})
	}
	return text
}

func processFormButton(textParam, attrs string) string {
	if attrs == "" {
		return fmt.Sprintf(`<button>{!! %s !!}</button>`, textParam)
	}

	attrProcessor := &AttributeProcessor{
		Order: []string{"type", "onclick", "class", "id", "disabled"},
		Patterns: map[string]string{
			"type":     `'type'\s*=>\s*'([^']+)'`,
			"onclick":  `'onclick'\s*=>\s*'([^']+)'`,
			"class":    `'class'\s*=>\s*'([^']+)'`,
			"id":       `'id'\s*=>\s*'([^']+)'`,
			"disabled": `'disabled'\s*=>\s*'([^']+)'`,
		},
	}

	// 静的属性の処理
	extraAttrs := attrProcessor.ProcessAttributes(attrs)

	// data-属性の追加処理
	dataRe := regexCache.GetRegex(`'(data-[^']+)'\s*=>\s*'([^']+)'`)
	for _, match := range dataRe.FindAllStringSubmatch(attrs, -1) {
		extraAttrs += fmt.Sprintf(` %s="%s"`, match[1], match[2])
	}

	// 動的属性の処理
	dynamicAttrs := processDynamicAttributes(attrs)
	extraAttrs += dynamicAttrs

	return fmt.Sprintf(`<button%s>{!! %s !!}</button>`, extraAttrs, textParam)
}

func replaceFormTextarea(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::textarea\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::textarea\(\s*(.*?)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormTextarea(params)
			}
			return match
		})
	}

	return text
}

func processFormTextarea(params []string) string {
	if len(params) < 1 {
		return ""
	}
	name := ProcessFieldName(params[0])
	value := ""

	if len(params) > 1 {
		value = params[1]
	}

	attrProcessor := &AttributeProcessor{
		Order: []string{"cols", "rows", "placeholder", "class"},
		Patterns: map[string]string{
			"cols":        `'cols'\s*=>\s*(\d+)`,
			"rows":        `'rows'\s*=>\s*(?:'([^']+)'|(\d+))`,
			"placeholder": `'placeholder'\s*=>\s*'([^']+)'`,
			"class":       `'class'\s*=>\s*'([^']+)'`,
		},
	}

	extraAttrs := ""
	if len(params) > 2 {
		extraAttrs = attrProcessor.ProcessAttributes(params[2])
	}

	// 値なしまたは空値の場合の処理
	if len(params) < 2 || value == "" {
		return fmt.Sprintf(`<textarea name="%s"%s></textarea>`, name, extraAttrs)
	}

	formattedValue := FormatValueAttribute(value)
	if formattedValue == "" {
		return fmt.Sprintf(`<textarea name="%s"%s></textarea>`, name, extraAttrs)
	}
	return fmt.Sprintf(`<textarea name="%s"%s>%s</textarea>`, name, extraAttrs, formattedValue)
}

func replaceFormLabel(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::label\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::label\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormLabel(params)
			}
			return match
		})
	}
	return text
}

func processFormLabel(params []string) string {
	if len(params) < 1 {
		return ""
	}

	// PHP文字列連結を含むフィールド名を適切に処理
	name := ProcessFieldName(params[0])
	forAttr := name // デフォルトでは名前をfor属性に使用
	textParam := ""

	if len(params) == 1 {
		textParam = fmt.Sprintf("'%s'", name)
	} else {
		textParam = params[1]
	}

	// 属性処理の統一
	attrProcessor := &AttributeProcessor{
		Order: []string{"class", "id", "style"},
		Patterns: map[string]string{
			"class": `'class'\s*=>\s*'([^']+)'`,
			"id":    `'id'\s*=>\s*'([^']+)'`,
			"style": `'style'\s*=>\s*'([^']+)'`,
		},
	}

	extraAttrs := ""
	if len(params) > 2 {
		attrs := params[2]

		// for属性の特別処理
		forRe := regexCache.GetRegex(`'for'\s*=>\s*'([^']+)'`)
		if forRe.MatchString(attrs) {
			forAttr = forRe.FindStringSubmatch(attrs)[1]
		}

		extraAttrs = attrProcessor.ProcessAttributes(attrs)
	}

	return fmt.Sprintf(`<label for="%s"%s>{!! %s !!}</label>`, forAttr, extraAttrs, textParam)
}

func processFormInputDynamic(params []string) string {
	if len(params) < 2 {
		return ""
	}

	// 最初のパラメータからinput typeを取得
	inputType := strings.Trim(params[0], `'"`)

	// 残りのパラメータ（name, value, attributes）を processFormInput に渡す
	remainingParams := params[1:]

	return processFormInput(inputType, remainingParams)
}

func replaceFormSelect(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\{\s*Form::select\(\s*(.*?)\s*\)\s*\}\}`,
		`(?s)\{\!\!\s*Form::select\(\s*(.*?)\s*\)\s*\!\!\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormSelect(params)
			}
			return match
		})
	}
	return text
}

func processFormSelect(params []string) string {
	if len(params) < 2 {
		return ""
	}

	// PHP文字列連結を含むフィールド名を適切に処理
	name := ProcessFieldName(params[0])
	options := params[1]
	selected := ""
	if len(params) > 2 {
		selected = params[2]
	}

	attrProcessor := &AttributeProcessor{
		Order: []string{"class", "id", "onchange"},
		Patterns: map[string]string{
			"class":    `'class'\s*=>\s*'([^']+)'`,
			"id":       `'id'\s*=>\s*'([^']+)'`,
			"onchange": `'(?:onChange|onchange)'\s*=>\s*'([^']+)'`,
		},
	}

	extraAttrs := ""
	if len(params) > 3 {
		extraAttrs = attrProcessor.ProcessAttributes(params[3])
	}

	selectHTML := fmt.Sprintf(`<select name="%s"%s>
@foreach(%s as $key => $value)
<option value="{{ $key }}" @if($key == %s) selected @endif>{{ $value }}</option>
@endforeach
</select>`, name, extraAttrs, options, selected)
	return selectHTML
}

func replaceFormCheckbox(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\{\s*Form::checkbox\(\s*(.*?)\s*\)\s*\}\}`,
		`(?s)\{\!\!\s*Form::checkbox\(\s*(.*?)\s*\)\s*\!\!\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormCheckbox(params)
			}
			return match
		})
	}
	return text
}

func processFormCheckbox(params []string) string {
	if len(params) < 1 {
		return ""
	}

	// PHP文字列連結を含むフィールド名を適切に処理
	name := ProcessFieldName(params[0])
	value := ""
	if len(params) > 1 {
		value = strings.Trim(params[1], `'"`)
	}

	checked := ""
	if len(params) > 2 {
		checked = params[2]
	}

	// 属性処理の統一
	attrProcessor := &AttributeProcessor{
		Order: []string{"class", "id", "style", "disabled", "onClick", "onChange"},
		Patterns: map[string]string{
			"class":    `'class'\s*=>\s*(.+?)(?:\s*,|\s*\]|$)`,
			"id":       `'id'\s*=>\s*(.+?)(?:\s*,|\s*\]|$)`,
			"style":    `'style'\s*=>\s*(.+?)(?:\s*,|\s*\]|$)`,
			"disabled": `'disabled'\s*=>\s*(.+?)(?:\s*,|\s*\]|$)`,
			"onClick":  `'onClick'\s*=>\s*'([^']+)'`,
			"onChange": `'onChange'\s*=>\s*'([^']+)'`,
		},
	}

	extraAttrs := ""
	if len(params) > 3 {
		// data-属性の追加処理（最初に処理）
		dataRe := regexCache.GetRegex(`'(data-[^']+)'\s*=>\s*'([^']+)'`)
		for _, match := range dataRe.FindAllStringSubmatch(params[3], -1) {
			extraAttrs += fmt.Sprintf(` %s="%s"`, match[1], match[2])
		}

		// その他の属性の処理
		extraAttrs += attrProcessor.ProcessAttributes(params[3])
	}

	// HTML出力を生成
	var result string
	// 配列形式の名前かどうかで処理を分岐（Laravel の配列形式サポートのため）
	if strings.HasSuffix(name, "[]") {
		result = fmt.Sprintf(`<input type="checkbox" name="%s" value="{{ %s }}" @if(in_array(%s, (array)%s)) checked @endif%s>`,
			name, value, value, checked, extraAttrs)
	} else {
		result = fmt.Sprintf(`<input type="checkbox" name="%s" value="{{ %s }}" @if(%s) checked @endif%s>`, name, value, checked, extraAttrs)
	}

	// onclick/onchange属性のダブルクォート変換を最終出力に適用
	result = convertEventHandlerQuotesInHTML(result)

	return result
}

func replaceFormSubmit(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?i)は大文字小文字を無視、(?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?is)\{\!\!\s*Form::submit\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?is)\{\{\s*Form::submit\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormSubmit(params)
			}
			return match
		})
	}
	return text
}

func processFormSubmit(params []string) string {
	if len(params) < 1 {
		return ""
	}

	textParam := strings.Trim(params[0], `'"`)

	if textParam == "null" {
		textParam = ""
	}

	attrProcessor := &AttributeProcessor{
		Order: []string{"class", "id", "style", "onclick", "disabled"},
		Patterns: map[string]string{
			"class":    `'class'\s*=>\s*'([^']+)'`,
			"id":       `'id'\s*=>\s*'([^']+)'`,
			"style":    `'style'\s*=>\s*'([^']+)'`,
			"onclick":  `'onclick'\s*=>\s*'([^']+)'`,
			"disabled": `'disabled'\s*=>\s*'([^']*)'`,
		},
	}

	extraAttrs := ""
	if len(params) > 1 {
		extraAttrs = attrProcessor.ProcessAttributes(params[1])
	}

	return fmt.Sprintf(`<button type="submit"%s>%s</button>`, extraAttrs, textParam)
}

func replaceFormDate(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::date\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::date\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				params := extractParamsBalanced(paramStr)
				return processFormInput("date", params)
			}
			return match
		})
	}
	return text
}

func replaceFormTime(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::time\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::time\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				params := extractParamsBalanced(paramStr)
				return processFormInput("time", params)
			}
			return match
		})
	}
	return text
}

func replaceFormDatetime(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::datetime\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::datetime\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormInput("datetime-local", params)
			}
			return match
		})
	}
	return text
}

func replaceFormColor(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::color\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::color\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormInput("color", params)
			}
			return match
		})
	}
	return text
}

func replaceFormRadio(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::radio\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::radio\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				params := extractParamsBalanced(paramStr)
				return processFormRadio(params)
			}
			return match
		})
	}
	return text
}

func processFormRadio(params []string) string {
	if len(params) < 2 {
		return ""
	}

	// PHP文字列連結を含むフィールド名を適切に処理
	name := ProcessFieldName(params[0])

	// 値属性の適切なフォーマット（Radioボタンは常に元の形式を保持）
	rawValue := strings.TrimSpace(params[1])
	var value string
	if rawValue == "" || rawValue == "null" || rawValue == "''" || rawValue == `""` {
		value = ""
	} else {
		value = fmt.Sprintf("{{ %s }}", params[1])
	}

	checked := ""
	if len(params) > 2 {
		checked = params[2]
	}

	// 属性処理の統一（onchangeサポート追加）
	attrProcessor := &AttributeProcessor{
		Order: []string{"id", "class", "style", "onchange", "disabled"},
		Patterns: map[string]string{
			"id":       `'id'\s*=>\s*'([^']+)'`,
			"class":    `'class'\s*=>\s*'([^']+)'`,
			"style":    `'style'\s*=>\s*'([^']+)'`,
			"onchange": `'onchange'\s*=>\s*'([^']+)'`,
			"disabled": `'disabled'\s*=>\s*'([^']*)'`,
		},
	}

	extraAttrs := ""
	if len(params) > 3 {
		extraAttrs = attrProcessor.ProcessAttributes(params[3])
	}

	// ラジオボタンのチェック状態処理
	checkedAttr := ""
	if checked != "" && checked != "false" && checked != "null" {
		checkedAttr = fmt.Sprintf(" @if(%s) checked @endif", checked)
	}

	return fmt.Sprintf(`<input type="radio" name="%s" value="%s"%s%s>`, name, value, checkedAttr, extraAttrs)
}

func printSummary(config *ReplacementConfig) {
	fmt.Println()
	fmt.Println("=== 置換結果サマリー ===")
	if config.IsFile {
		fmt.Printf("対象ファイル: %s\n", config.TargetPath)
	} else {
		fmt.Printf("対象ディレクトリ: %s\n", config.TargetPath)
	}
	fmt.Printf("処理したファイル数: %d\n", config.FileCount)
	fmt.Println()
	if len(config.ProcessedFiles) > 0 {
		fmt.Println("=== 処理済みファイル ===")
		for _, file := range config.ProcessedFiles {
			fmt.Printf(" - %s\n", file)
		}
		fmt.Println()
	}
	var remainingFiles []string
	if config.IsFile {
		if hasFormFacade, _ := containsFormFacade(config.TargetPath); hasFormFacade {
			remainingFiles = []string{config.TargetPath}
		}
	} else {
		remainingFiles = findRemainingFormFacades(config.TargetPath)
	}
	if len(remainingFiles) > 0 {
		fmt.Println("=== Form facadeが残存するファイル ===")
		for _, file := range remainingFiles {
			fmt.Println(file)
		}
		fmt.Println()
		fmt.Println("=== 残存するForm facadeパターン ===")
		showRemainingPatterns(remainingFiles)
	} else {
		fmt.Println("Form facadeを含むファイルは見つかりませんでした（置換完了）")
	}
	fmt.Println()
	fmt.Println("置換処理が完了しました！")
}

func findRemainingFormFacades(targetDir string) []string {
	var remainingFiles []string
	filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".blade.php") {
			hasFormFacade, err := containsFormFacade(path)
			if err == nil && hasFormFacade {
				remainingFiles = append(remainingFiles, path)
			}
		}
		return nil
	})
	sort.Strings(remainingFiles)
	return remainingFiles
}

func showRemainingPatterns(files []string) {
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.Contains(line, "Form::") {
				fmt.Printf("%s:%d:%s\n", file, i+1, strings.TrimSpace(line))
			}
		}
	}
}
