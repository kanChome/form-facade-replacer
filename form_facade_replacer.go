package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// 正規表現キャッシュ
type RegexCache struct {
	mu    sync.RWMutex
	cache map[string]*regexp.Regexp
}

var regexCache = &RegexCache{
	cache: make(map[string]*regexp.Regexp),
}

// GetRegex 正規表現の取得（キャッシュあり）
func (rc *RegexCache) GetRegex(pattern string) *regexp.Regexp {
	rc.mu.RLock()
	if re, exists := rc.cache[pattern]; exists {
		rc.mu.RUnlock()
		return re
	}
	rc.mu.RUnlock()

	rc.mu.Lock()
	defer rc.mu.Unlock()

	// ダブルチェック（他のゴルーチンが既に作成している可能性）
	if re, exists := rc.cache[pattern]; exists {
		return re
	}

	re := regexp.MustCompile(pattern)
	rc.cache[pattern] = re
	return re
}

type ReplacementConfig struct {
	TargetPath     string
	IsFile         bool
	ProcessedFiles []string
	FileCount      int
}

func main() {
	config := &ReplacementConfig{
		TargetPath:     "",
		IsFile:         false,
		ProcessedFiles: make([]string, 0),
		FileCount:      0,
	}

	if len(os.Args) < 2 {
		fmt.Println("エラー: ファイルまたはディレクトリを指定してください。")
		printUsage()
		os.Exit(1)
	}

	arg := os.Args[1]
	if arg == "--help" || arg == "-h" {
		printUsage()
		return
	}

	config.TargetPath = arg

	info, err := os.Stat(config.TargetPath)
	if err != nil {
		log.Fatalf("エラー: '%s' が存在しません。", config.TargetPath)
	}

	config.IsFile = !info.IsDir()

	if config.IsFile {
		if !strings.HasSuffix(config.TargetPath, ".blade.php") {
			log.Fatalf("エラー: '%s' は.blade.phpファイルではありません。", config.TargetPath)
		}
		fmt.Printf("Form Facade置換を開始します (ファイル): %s\n", config.TargetPath)
	} else {
		fmt.Printf("Form Facade置換を開始します (ディレクトリ): %s\n", config.TargetPath)
	}

	err = processBladeFiles(config)
	if err != nil {
		log.Fatalf("ファイル処理中にエラーが発生しました: %v", err)
	}

	printSummary(config)
}

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
	fmt.Println()
	fmt.Println("例:")
	fmt.Println(" go run form_facade_replacer.go resources/views/hoge")
	fmt.Println(" go run form_facade_replacer.go resources/views/hoge/fuga.blade.php")
}

func processBladeFiles(config *ReplacementConfig) error {
	if config.IsFile {
		return processSingleFile(config, config.TargetPath)
	}

	return filepath.WalkDir(config.TargetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".blade.php") {
			return processSingleFile(config, path)
		}
		return nil
	})
}

func processSingleFile(config *ReplacementConfig, filePath string) error {
	hasFormFacade, err := containsFormFacade(filePath)
	if err != nil {
		return fmt.Errorf("ファイル %s のチェックに失敗しました: %v", filePath, err)
	}

	if hasFormFacade {
		fmt.Printf("処理中: %s\n", filePath)
		err := replaceFormPatterns(filePath)
		if err != nil {
			return fmt.Errorf("ファイル %s の処理に失敗しました: %v", filePath, err)
		}
		config.ProcessedFiles = append(config.ProcessedFiles, filePath)
		config.FileCount++
		fmt.Printf(" - 処理完了: %s\n", filePath)
	}
	return nil
}

func containsFormFacade(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Form::") {
			return true, nil
		}
	}
	return false, scanner.Err()
}

// 共通の正規表現パターンとヘルパー関数

// BladePattern 定数 - Bladeの基本パターン
const (
	BladeExclamationPattern = `\{\!\!\s*%s\s*\!\!\}`
	BladeCurlyPattern       = `\{\{\s*%s\s*\}\}`
)

// AttributeProcessor 属性処理の統一インターフェース
type AttributeProcessor struct {
	Order    []string
	Patterns map[string]string
}

// ProcessAttributes 属性を統一された順序で処理
func (ap *AttributeProcessor) ProcessAttributes(attrs string) string {
	var extraAttrs string
	for _, attr := range ap.Order {
		if pattern, exists := ap.Patterns[attr]; exists {
			if re := regexCache.GetRegex(pattern); re.MatchString(attrs) {
				matches := re.FindStringSubmatch(attrs)
				var val string
				if len(matches) > 2 && matches[2] != "" {
					val = matches[2] // 数値の場合
				} else {
					val = matches[1] // 文字列の場合
				}
				// disabled属性の特別処理
				if attr == "disabled" && (val == "" || val == "disabled") {
					extraAttrs += fmt.Sprintf(` %s`, attr)
				} else {
					extraAttrs += fmt.Sprintf(` %s="%s"`, attr, val)
				}
			}
		}
	}
	return extraAttrs
}

func DetectArrayHelper(value string) bool {
	return regexCache.GetRegex(`(?i)^(old|session|request|input)\s*\(`).MatchString(strings.TrimSpace(value))
}

// IsArrayFieldName 配列形式のフィールド名かどうかを判定
func IsArrayFieldName(fieldName string) bool {
	return regexCache.GetRegex(`\[.*\]`).MatchString(fieldName)
}

func FormatValueAttribute(value string) string {
	// 空値、null、空文字の場合は空文字を返す
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" || trimmedValue == "null" || trimmedValue == "''" || trimmedValue == `""` {
		return ""
	}

	// 通常のBlade出力を使用（Form::hiddenと一貫した動作）
	return fmt.Sprintf("{{ %s }}", value)
}

// FormatHiddenValueAttribute hidden input専用の値フォーマット
func FormatHiddenValueAttribute(value string, fieldName string) string {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" || trimmedValue == "null" || trimmedValue == "''" || trimmedValue == `""` {
		return ""
	}

	// 配列フィールドの場合は文字列結合形式を使用
	if IsArrayFieldName(fieldName) {
		return fmt.Sprintf("{{ is_array(%s) ? implode(',', %s) : %s }}", value, value, value)
	}

	// 通常フィールドの場合は通常のBlade出力
	return fmt.Sprintf("{{ %s }}", value)
}

func ProcessBladePatterns(text, formMethod string, processor func(string) string) string {
	patterns := []string{
		fmt.Sprintf(BladeExclamationPattern, formMethod),
		fmt.Sprintf(BladeCurlyPattern, formMethod),
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			matches := re.FindStringSubmatch(match)
			if len(matches) > 1 {
				return processor(matches[1])
			}
			return processor("")
		})
	}
	return text
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
	text = replaceFormNumber(text)
	text = replaceFormSelect(text)
	text = replaceFormCheckbox(text)
	text = replaceFormSubmit(text)
	text = replaceFormFile(text)

	return os.WriteFile(filePath, []byte(text), 0644)
}

func replaceFormOpen(text string) string {
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::open\(\s*\[\s*(.*?)\s*\]\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::open\(\s*\[\s*(.*?)\s*\]\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			return processFormOpen(re.FindStringSubmatch(match)[1])
		})
	}
	return text
}

func processFormOpen(content string) string {
	action := extractFormAction(content)
	method := extractFormMethod(content)
	extraAttrs := extractFormAttributes(content)

	return buildFormTag(action, method, extraAttrs)
}

// extractFormAction アクション属性の抽出（route、url の優先順位処理）
func extractFormAction(content string) string {
	// まず配列形式のroute（パラメータ付き）をチェック: 'route' => ['user.store', ['id' => 1]]
	paramRouteRe := regexCache.GetRegex(`'route'\s*=>\s*\[\s*'([^']+)'\s*,\s*(\[[^\]]*\])`)
	if paramMatches := paramRouteRe.FindStringSubmatch(content); len(paramMatches) > 2 {
		return fmt.Sprintf("{{ route('%s', %s) }}", paramMatches[1], paramMatches[2])
	}

	// 次に配列形式のroute（パラメータなし）をチェック: 'route' => ['user.index']
	arrayRouteRe := regexCache.GetRegex(`'route'\s*=>\s*\[\s*'([^']+)'\s*\]`)
	if arrayMatches := arrayRouteRe.FindStringSubmatch(content); len(arrayMatches) > 1 {
		return fmt.Sprintf("{{ route('%s') }}", arrayMatches[1])
	}

	// 最後に文字列形式のroute: 'route' => 'user.index'
	simpleRouteRe := regexCache.GetRegex(`'route'\s*=>\s*'([^']+)'`)
	if simpleMatches := simpleRouteRe.FindStringSubmatch(content); len(simpleMatches) > 1 {
		return fmt.Sprintf("{{ route('%s') }}", simpleMatches[1])
	}

	// route が見つからなかった場合のみ url 処理を実行（route が優先）
	return extractFormUrl(content)
}

// extractFormUrl URL属性の抽出
func extractFormUrl(content string) string {
	urlRe := regexCache.GetRegex(`'url'\s*=>\s*([^,\]]+)`)
	if matches := urlRe.FindStringSubmatch(content); len(matches) > 1 {
		urlVal := strings.TrimSpace(matches[1])
		if strings.HasPrefix(urlVal, "route(") {
			routeFuncRe := regexCache.GetRegex(`route\([^)]*(?:\([^)]*\)[^)]*)*\)`)
			if routeMatch := routeFuncRe.FindString(content); routeMatch != "" {
				return fmt.Sprintf("{{ %s }}", routeMatch)
			}
			return fmt.Sprintf("{{ %s }}", urlVal)
		}
		return urlVal
	}
	return ""
}

// extractFormMethod HTTP メソッドの抽出
func extractFormMethod(content string) string {
	methodRe := regexCache.GetRegex(`'method'\s*=>\s*'([^']+)'`)
	if methodRe.MatchString(content) {
		return methodRe.FindStringSubmatch(content)[1]
	}
	return "GET"
}

// extractFormAttributes フォーム属性の抽出
func extractFormAttributes(content string) string {
	attrProcessor := &AttributeProcessor{
		Order: []string{"class", "id", "target"},
		Patterns: map[string]string{
			"target": `'target'\s*=>\s*'([^']+)'`,
			"id":     `'id'\s*=>\s*'([^']+)'`,
			"class":  `'class'\s*=>\s*'([^']+)'`,
		},
	}
	return attrProcessor.ProcessAttributes(content)
}

// buildFormTag フォームタグの構築
func buildFormTag(action, method, extraAttrs string) string {
	if strings.ToUpper(method) == "GET" {
		return fmt.Sprintf(`<form action="%s" method="%s"%s>`, action, method, extraAttrs)
	}
	return fmt.Sprintf(`<form action="%s" method="%s"%s>
{{ csrf_field() }}`, action, method, extraAttrs)
}

func replaceFormClose(text string) string {
	return ProcessBladePatterns(text, `Form::close\(\)`, func(content string) string {
		return "</form>"
	})
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

	name := strings.Trim(params[0], `'"`)
	value := ""
	if len(params) > 1 {
		value = params[1]
	}

	nameAttr := name
	if strings.Contains(name, " . ") {
		patterns := []string{
			`^'([^']*)'\\s*\\.\\s*([^'\\s]+(?:\\[[^\\]]*\\]\\[[^\\]]*\\])?)\\s*\\.\\s*'([^']*)'$`,
			`^'([^']*)'\\s*\\.\\s*(.+?)\\s*\\.\\s*'([^']*)'$`,
		}

		for _, pattern := range patterns {
			re := regexCache.GetRegex(pattern)
			if matches := re.FindStringSubmatch(name); len(matches) == 4 {
				nameAttr = fmt.Sprintf("%s{{ %s }}%s", matches[1], matches[2], matches[3])
				break
			}
		}
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

	extraAttrs := attrProcessor.ProcessAttributes(attrs)

	// data-属性の追加処理
	dataRe := regexCache.GetRegex(`'(data-[^']+)'\s*=>\s*'([^']+)'`)
	for _, match := range dataRe.FindAllStringSubmatch(attrs, -1) {
		extraAttrs += fmt.Sprintf(` %s="%s"`, match[1], match[2])
	}

	return fmt.Sprintf(`<button%s>{!! %s !!}</button>`, extraAttrs, textParam)
}

func replaceFormTextarea(text string) string {
	patterns := []string{
		`\{\!\!\s*Form::textarea\(\s*([^}]+)\s*\)\s*\!\!\}`,
		`\{\{\s*Form::textarea\(\s*([^}]+)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			params := extractParams(re.FindStringSubmatch(match)[1])
			return processFormTextarea(params)
		})
	}

	return text
}

func processFormTextarea(params []string) string {
	if len(params) < 1 {
		return ""
	}
	name := strings.Trim(params[0], `'"`)
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
	patterns := []string{
		`\{\!\!\s*Form::label\(\s*([^}]+)\s*\)\s*\!\!\}`,
		`\{\{\s*Form::label\(\s*([^}]+)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			params := extractParams(re.FindStringSubmatch(match)[1])
			return processFormLabel(params)
		})
	}
	return text
}

func processFormLabel(params []string) string {
	if len(params) < 1 {
		return ""
	}

	name := strings.Trim(params[0], `'"`)
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

func replaceFormText(text string) string {
	patterns := []string{
		`\{\!\!\s*Form::text\(\s*([^}]+)\s*\)\s*\!\!\}`,
		`\{\{\s*Form::text\(\s*([^}]+)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			params := extractParams(re.FindStringSubmatch(match)[1])
			return processFormInput("text", params)
		})
	}
	return text
}

func replaceFormFile(text string) string {
	// より複雑なネストに対応したパターン
	// {!! Form::file(...) !!} と {{ Form::file(...) }} の両方に対応
	patterns := []string{
		`\{\!\!\s*Form::file\(\s*(.*?)\s*\)\s*\!\!\}`,
		`\{\{\s*Form::file\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				params := extractParamsBalanced(paramStr)
				return processFormFile(params)
			}
			return match
		})
	}
	return text
}

// extractParamsBalanced ネストした構造を考慮したパラメータ抽出
func extractParamsBalanced(paramsStr string) []string {
	var params []string
	var current strings.Builder
	var parenCount, bracketCount, braceCount int
	var inQuotes bool
	var quoteChar rune
	var escapeNext bool

	for _, char := range paramsStr {
		if escapeNext {
			current.WriteRune(char)
			escapeNext = false
			continue
		}

		if char == '\\' && inQuotes {
			current.WriteRune(char)
			escapeNext = true
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
					params = append(params, strings.TrimSpace(current.String()))
					current.Reset()
					continue
				}
			}
		}

		current.WriteRune(char)
	}

	if current.Len() > 0 {
		params = append(params, strings.TrimSpace(current.String()))
	}

	return params
}

// convertJavaScriptStringLiterals JavaScript文字列リテラル内のダブルクォートをシングルクォートに変換
func convertJavaScriptStringLiterals(jsCode string) string {
	result := jsCode

	// より柔軟なアプローチ: ダブルクォートで囲まれた部分を特定して変換
	// エスケープされたクォートも含めて、正しく処理する

	// パターン1: JSON-like文字列を先に処理 "{\"key\": \"value\"}"
	jsonStringPattern := `"(\{(?:[^"\\]|\\.)*\})"`
	jsonRe := regexCache.GetRegex(jsonStringPattern)
	result = jsonRe.ReplaceAllString(result, "'$1'")

	// パターン2: エスケープを含む文字列リテラル "Say \"Hello\""
	// 完全なエスケープ対応パターン
	escapedStringPattern := `"((?:[^"\\]|\\.)*)"`
	escapedRe := regexCache.GetRegex(escapedStringPattern)
	result = escapedRe.ReplaceAllString(result, "'$1'")

	return result
}

// convertEventHandlerQuotesInHTML HTML出力内のイベントハンドラ属性のダブルクォートを柔軟に変換
func convertEventHandlerQuotesInHTML(html string) string {
	result := html

	// onchange属性の処理
	// パターン: onchange="JavaScript code with "quotes""
	onchangePattern := `(onchange=")([^"]*(?:"[^"]*)*)(")`
	onchangeRe := regexCache.GetRegex(onchangePattern)

	result = onchangeRe.ReplaceAllStringFunc(result, func(match string) string {
		matches := onchangeRe.FindStringSubmatch(match)
		if len(matches) >= 4 {
			prefix := matches[1] // onchange="
			jsCode := matches[2] // JavaScript コード部分（内部にダブルクォートを含む可能性）
			suffix := matches[3] // 最後の "

			// JavaScript内の文字列リテラルを変換
			convertedJS := convertJavaScriptStringLiterals(jsCode)
			return prefix + convertedJS + suffix
		}
		return match
	})

	// onclick属性の処理
	onclickPattern := `(onclick=")([^"]*(?:"[^"]*)*)(")`
	onclickRe := regexCache.GetRegex(onclickPattern)

	result = onclickRe.ReplaceAllStringFunc(result, func(match string) string {
		matches := onclickRe.FindStringSubmatch(match)
		if len(matches) >= 4 {
			prefix := matches[1] // onclick="
			jsCode := matches[2] // JavaScript コード部分
			suffix := matches[3] // 最後の "

			// JavaScript内の文字列リテラルを変換
			convertedJS := convertJavaScriptStringLiterals(jsCode)
			return prefix + convertedJS + suffix
		}
		return match
	})

	return result
}

func processFormFile(params []string) string {
	if len(params) < 1 {
		return ""
	}

	name := strings.Trim(params[0], `'"`)

	attrProcessor := &AttributeProcessor{
		Order: []string{"accept", "capture", "class", "id", "onchange", "onclick"},
		Patterns: map[string]string{
			"accept":   `'accept'\s*=>\s*'([^']+)'`,
			"capture":  `'capture'\s*=>\s*'([^']+)'`,
			"id":       `'id'\s*=>\s*'([^']+)'`,
			"class":    `'class'\s*=>\s*'([^']+)'`,
			"onchange": `'onchange'\s*=>\s*'([^']+)'`,
			"onclick":  `'onclick'\s*=>\s*'([^']+)'`,
		},
	}

	extraAttrs := ""
	multipleAttr := ""

	if len(params) > 1 {
		// multiple属性の特別処理（先に処理）
		multiplePattern := `'multiple'\s*=>\s*(true|false|\d+)`
		if re := regexCache.GetRegex(multiplePattern); re.MatchString(params[1]) {
			matches := re.FindStringSubmatch(params[1])
			if len(matches) > 1 {
				val := matches[1]
				if val == "true" {
					multipleAttr = " multiple"
				}
				// false や 0 の場合は何も追加しない
			}
		}

		extraAttrs = attrProcessor.ProcessAttributes(params[1])
	}

	// HTML出力を生成
	result := fmt.Sprintf(`<input type="file" name="%s"%s%s>`, name, extraAttrs, multipleAttr)

	// onchange/onclick属性のダブルクォート変換を最終出力に適用
	result = convertEventHandlerQuotesInHTML(result)

	return result
}

func replaceFormNumber(text string) string {
	patterns := []string{
		`\{\{\s*Form::number\(\s*([^}]+)\s*\)\s*\}\}`,
		`\{\!\!\s*Form::number\(\s*([^}]+)\s*\)\s*\!\!\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			params := extractParams(re.FindStringSubmatch(match)[1])
			return processFormNumber(params)
		})
	}
	return text
}

func processFormNumber(params []string) string {
	if len(params) < 1 {
		return ""
	}

	name := strings.Trim(params[0], `'"`)
	value := ""
	if len(params) > 1 {
		value = params[1]
	}

	attrProcessor := &AttributeProcessor{
		Order: []string{"placeholder", "class", "id", "min", "max", "step"},
		Patterns: map[string]string{
			"placeholder": `'placeholder'\s*=>\s*'([^']+)'`,
			"class":       `'class'\s*=>\s*'([^']+)'`,
			"id":          `'id'\s*=>\s*'([^']+)'`,
			"min":         `'min'\s*=>\s*(\d+)`,
			"max":         `'max'\s*=>\s*(\d+)`,
			"step":        `'step'\s*=>\s*(\d+)`,
		},
	}

	extraAttrs := ""
	if len(params) > 2 {
		extraAttrs = attrProcessor.ProcessAttributes(params[2])
	}

	// HTMLとして無効な値（null、空文字列）の場合、value属性を出力しない
	valueAttr := ""
	if value != "" {
		rawValue := strings.TrimSpace(value)
		if rawValue != "null" && rawValue != "''" && rawValue != `""` {
			formattedValue := FormatValueAttribute(value)
			valueAttr = fmt.Sprintf(` value="%s"`, formattedValue)
		}
	}

	return fmt.Sprintf(`<input type="number" name="%s"%s%s>`, name, valueAttr, extraAttrs)
}

func processFormInput(inputType string, params []string) string {
	if len(params) < 1 {
		return ""
	}

	name := strings.Trim(params[0], `'"`)
	value := ""
	if len(params) > 1 {
		value = params[1]
	}

	attrProcessor := &AttributeProcessor{
		Order: []string{"placeholder", "class", "id"},
		Patterns: map[string]string{
			"placeholder": `'placeholder'\s*=>\s*'([^']+)'`,
			"class":       `'class'\s*=>\s*'([^']+)'`,
			"id":          `'id'\s*=>\s*'([^']+)'`,
		},
	}

	extraAttrs := ""
	if len(params) > 2 {
		extraAttrs = attrProcessor.ProcessAttributes(params[2])
	}

	formattedValue := FormatValueAttribute(value)
	return fmt.Sprintf(`<input type="%s" name="%s" value="%s"%s>`, inputType, name, formattedValue, extraAttrs)
}

func replaceFormSelect(text string) string {
	patterns := []string{
		`\{\{\s*Form::select\(\s*([^}]+)\s*\)\s*\}\}`,
		`\{\!\!\s*Form::select\(\s*([^}]+)\s*\)\s*\!\!\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			params := extractParams(re.FindStringSubmatch(match)[1])
			return processFormSelect(params)
		})
	}
	return text
}

func processFormSelect(params []string) string {
	if len(params) < 2 {
		return ""
	}

	name := strings.Trim(params[0], `'"`)
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
	patterns := []string{
		`\{\{\s*Form::checkbox\(\s*([^}]+)\s*\)\s*\}\}`,
		`\{\!\!\s*Form::checkbox\(\s*([^}]+)\s*\)\s*\!\!\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			params := extractParams(re.FindStringSubmatch(match)[1])
			return processFormCheckbox(params)
		})
	}
	return text
}

func processFormCheckbox(params []string) string {
	if len(params) < 1 {
		return ""
	}

	name := strings.Trim(params[0], `'"`)
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
		Order: []string{"class", "id", "style", "disabled"},
		Patterns: map[string]string{
			"class":    `'class'\s*=>\s*'([^']+)'`,
			"id":       `'id'\s*=>\s*'([^']+)'`,
			"style":    `'style'\s*=>\s*'([^']+)'`,
			"disabled": `'disabled'\s*=>\s*'([^']*)'`,
		},
	}

	extraAttrs := ""
	if len(params) > 3 {
		extraAttrs = attrProcessor.ProcessAttributes(params[3])
	}

	// 配列形式の名前かどうかで処理を分岐（Laravel の配列形式サポートのため）
	if strings.HasSuffix(name, "[]") {
		return fmt.Sprintf(`<input type="checkbox" name="%s" value="{{ %s }}" @if(in_array(%s, (array)%s)) checked @endif%s>`,
			name, value, value, checked, extraAttrs)
	} else {
		return fmt.Sprintf(`<input type="checkbox" name="%s" value="{{ %s }}" @if(%s) checked @endif%s>`, name, value, checked, extraAttrs)
	}
}

func replaceFormSubmit(text string) string {
	patterns := []string{
		`(?i)\{\!\!\s*Form::submit\(\s*([^}]+)\s*\)\s*\!\!\}`,
		`(?i)\{\{\s*Form::submit\(\s*([^}]+)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			params := extractParams(re.FindStringSubmatch(match)[1])
			return processFormSubmit(params)
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

	// 属性処理の統一
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

func extractParamsAdvanced(paramsStr string) []string {
	paramsStr = strings.ReplaceAll(paramsStr, "\n", " ")
	paramsStr = strings.ReplaceAll(paramsStr, "\t", " ")
	re := regexCache.GetRegex(`\s+`)
	paramsStr = re.ReplaceAllString(paramsStr, " ")
	paramsStr = strings.TrimSpace(paramsStr)
	return extractParams(paramsStr)
}

func extractParams(paramsStr string) []string {
	var params []string
	var current strings.Builder
	var parenCount, bracketCount int
	var inQuotes bool
	var quoteChar rune

	for _, char := range paramsStr {
		if !inQuotes && (char == '"' || char == '\'') {
			inQuotes = true
			quoteChar = char
		} else if inQuotes && char == quoteChar {
			inQuotes = false
			quoteChar = 0
		} else if !inQuotes {
			switch char {
			case '(':
				parenCount++
			case ')':
				parenCount--
			case '[':
				bracketCount++
			case ']':
				bracketCount--
			case ',':
				if parenCount == 0 && bracketCount == 0 {
					params = append(params, strings.TrimSpace(current.String()))
					current.Reset()
					continue
				}
			}
		}
		current.WriteRune(char)
	}

	if current.Len() > 0 {
		params = append(params, strings.TrimSpace(current.String()))
	}

	return params
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
