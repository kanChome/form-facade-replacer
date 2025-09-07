// ffr.go: パッケージ横断の共通ロジックとディスパッチャの実装。
package ffr

import (
	"fmt"
	"os"
	"strings"
)

// バージョン情報（リリース時にldフラグで設定される）
// バージョン情報（ldflagsで上書きされる）
var (
	version   = "dev"
	buildDate = "unknown"
)

// --- Meta (usage/version) ---
// printUsage は CLI の使用方法を表示する（internal/ffr/cli.go から利用）。
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

// printVersion はバージョンとビルド時刻を表示する。
func printVersion() {
	fmt.Printf("Form Facade Replacer %s\n", version)
	fmt.Printf("Build Date: %s\n", buildDate)
	fmt.Println()
	fmt.Println("Laravel Form Facade を HTML に変換する高性能 Go ツール")
	fmt.Println("https://github.com/ryohirano/form-facade-replacer")
}

// --- Dynamic Attributes ---
// DynamicAttributePair は動的属性のキーと値のペアを表す。
type DynamicAttributePair struct {
	Key   string // 動的キー（例: $condition ? 'disabled' : ''）
	Value string // 動的値（例: $condition ? 'disabled' : null）
}

// detectDynamicAttributes は文字列から動的属性を抽出して配列で返す。
func detectDynamicAttributes(attrs string) []DynamicAttributePair {
	return extractDynamicAttributesBalanced(attrs)
}

// extractDynamicAttributesBalanced はカンマ分割時に括弧/クォートのバランスを考慮して抽出する。
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

// parseDynamicAttributePair は文字列から key=>value 形式の動的属性を解析する。
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

// isDynamicAttribute は与えられた文字列が動的属性の形かどうかを判定する。
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

// processDynamicAttributes は動的属性を HTML 属性文字列へ展開する。
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

// --- Dispatcher ---
// replaceFormPatterns は1ファイルの内容を各 replaceXXX に順次通し、結果を書き戻す。
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

// --- Hidden ---
// replaceFormHidden は Form::hidden(...) を <input type="hidden"> に置換する。
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

// processFormHidden は hidden の name/value/属性を整形して最終HTMLを生成する。
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

// --- Dynamic Input helper ---
// processFormInputDynamic は Form::input(type, ...) を共通の processFormInput に中継する。
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

// --- Color ---
// replaceFormColor は Form::color(...) を <input type="color"> に置換する。
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
