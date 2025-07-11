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

// バージョン情報（リリース時にldフラグで設定される）
var (
	version   = "dev"
	buildDate = "unknown"
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

	if arg == "--version" || arg == "-v" {
		printVersion()
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
			// 基本パターンをまず試す
			if re := regexCache.GetRegex(pattern); re.MatchString(attrs) {
				matches := re.FindStringSubmatch(attrs)
				var val string
				if len(matches) > 2 && matches[2] != "" {
					val = matches[2] // 数値の場合
				} else {
					val = matches[1] // 文字列の場合
				}

				// PHP文字列連結を含む値を処理
				val = processAttributeValue(val)

				// disabled属性とrequired属性の特別処理
				if (attr == "disabled" && (val == "" || val == "disabled")) ||
					(attr == "required" && (val == "" || val == "required")) {
					extraAttrs += fmt.Sprintf(` %s`, attr)
				} else {
					extraAttrs += fmt.Sprintf(` %s="%s"`, attr, val)
				}
			}
		}
	}
	return extraAttrs
}

// processAttributeValue 属性値のPHP文字列連結を処理
func processAttributeValue(value string) string {
	// 元の値を保持（引用符を含む）
	originalValue := value

	// 外側の引用符を取り除く
	value = strings.Trim(value, `'"`)

	// PHP文字列連結パターンを検出して変換（ProcessFieldNameと同様のロジック）
	if strings.Contains(value, " . ") {
		// パターン1: 'prefix' . $variable . 'suffix'
		concatPattern1 := `^'([^']*)'\s*\.\s*(.+?)\s*\.\s*'([^']*)'$`
		re1 := regexCache.GetRegex(concatPattern1)

		if matches := re1.FindStringSubmatch(value); len(matches) == 4 {
			prefix := matches[1]
			variable := strings.TrimSpace(matches[2])
			suffix := matches[3]
			// 全体を{{ }}で囲む
			return fmt.Sprintf("{{ '%s' . %s . '%s' }}", prefix, variable, suffix)
		}

		// パターン2: 'prefix' . $variable (suffix無し)
		concatPattern2 := `^'([^']*)'\s*\.\s*(.+)$`
		re2 := regexCache.GetRegex(concatPattern2)

		if matches := re2.FindStringSubmatch(value); len(matches) == 3 {
			prefix := matches[1]
			variable := strings.TrimSpace(matches[2])
			// 全体を{{ }}で囲む
			return fmt.Sprintf("{{ '%s' . %s }}", prefix, variable)
		}

		// パターン3: PHP変数のみ（$var['key']形式）
		// 特に $row['id'] のような複雑な変数への対応
		concatPattern3 := `^'([^']*)'\s*\.\s*(\$[a-zA-Z_]\w*(?:\[[^\]]*\])*)\s*$`
		re3 := regexCache.GetRegex(concatPattern3)

		if matches := re3.FindStringSubmatch(value); len(matches) == 3 {
			prefix := matches[1]
			variable := strings.TrimSpace(matches[2])
			// 全体を{{ }}で囲む
			return fmt.Sprintf("{{ '%s' . %s }}", prefix, variable)
		}

		// 文字列連結が見つかった場合は、元の値を{{ }}で囲む
		return fmt.Sprintf("{{ %s }}", originalValue)
	}

	return value
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

func DetectArrayHelper(value string) bool {
	return regexCache.GetRegex(`(?i)^(old|session|request|input)\s*\(`).MatchString(strings.TrimSpace(value))
}

// IsArrayFieldName 配列形式のフィールド名かどうかを判定
func IsArrayFieldName(fieldName string) bool {
	return regexCache.GetRegex(`\[.*\]`).MatchString(fieldName)
}

// ProcessFieldName PHP文字列連結を含むフィールド名をBlade構文に変換
func ProcessFieldName(name string) string {
	// まず全体から外側のクォートを除去
	nameAttr := strings.Trim(name, `'"`)

	// PHP文字列連結パターンを検出して変換
	if strings.Contains(nameAttr, " . ") {
		// パターン1: 'prefix' . $variable . 'suffix' (標準パターン)
		// 変数部分が単一引用符を含む場合に対応
		concatPattern1 := `^'([^']*)'\s*\.\s*(.+?)\s*\.\s*'([^']*)'$`
		re1 := regexCache.GetRegex(concatPattern1)

		if matches := re1.FindStringSubmatch(nameAttr); len(matches) == 4 {
			prefix := matches[1]
			variable := strings.TrimSpace(matches[2])
			suffix := matches[3]
			return fmt.Sprintf("%s{{ %s }}%s", prefix, variable, suffix)
		}

		// パターン2: prefix[' . $variable . ']suffix (埋め込みパターン)
		// 変数部分が単一引用符を含む場合に対応
		concatPattern2 := `^([^']*\[)'\s*\.\s*(.+?)\s*\.\s*'(\][^']*)$`
		re2 := regexCache.GetRegex(concatPattern2)

		if matches := re2.FindStringSubmatch(nameAttr); len(matches) == 4 {
			prefix := matches[1]
			variable := strings.TrimSpace(matches[2])
			suffix := matches[3]
			return fmt.Sprintf("%s{{ %s }}%s", prefix, variable, suffix)
		}
	}

	return nameAttr
}

func FormatValueAttribute(value string) string {
	// 空値、null、空文字の場合は空文字を返す
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" || trimmedValue == "null" || trimmedValue == "''" || trimmedValue == `""` {
		return ""
	}

	// 引用符で囲まれた値の場合、中身をチェック
	if strings.HasPrefix(trimmedValue, "'") && strings.HasSuffix(trimmedValue, "'") {
		innerValue := strings.Trim(trimmedValue, "'")

		// 純粋な数値の場合は引用符を除去
		if regexCache.GetRegex(`^\d+(\.\d+)?$`).MatchString(innerValue) {
			return fmt.Sprintf("{{ %s }}", innerValue)
		}

		// カラーコード（16進数）の場合は引用符を除去
		if regexCache.GetRegex(`^#[0-9a-fA-F]{3,6}$`).MatchString(innerValue) {
			return fmt.Sprintf("{{ %s }}", innerValue)
		}
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
	paramRouteStart := regexCache.GetRegex(`'route'\s*=>\s*\[\s*'([^']+)'\s*,\s*`)
	if startMatch := paramRouteStart.FindStringSubmatchIndex(content); len(startMatch) > 3 {
		routeName := content[startMatch[2]:startMatch[3]]
		paramStart := startMatch[1] // カンマ以降の開始位置

		// パラメータ部分を抽出（バランスした括弧を考慮）
		params := extractRouteParamsBalanced(content[paramStart:])
		if params != "" {
			return fmt.Sprintf("{{ route('%s', %s) }}", routeName, params)
		}
	}

	arrayRouteRe := regexCache.GetRegex(`'route'\s*=>\s*\[\s*'([^']+)'\s*\]`)
	if arrayMatches := arrayRouteRe.FindStringSubmatch(content); len(arrayMatches) > 1 {
		return fmt.Sprintf("{{ route('%s') }}", arrayMatches[1])
	}

	simpleRouteRe := regexCache.GetRegex(`'route'\s*=>\s*'([^']+)'`)
	if simpleMatches := simpleRouteRe.FindStringSubmatch(content); len(simpleMatches) > 1 {
		return fmt.Sprintf("{{ route('%s') }}", simpleMatches[1])
	}

	// route が見つからなかった場合のみ url 処理を実行（route が優先）
	return extractFormUrl(content)
}

// extractRouteParamsBalanced ルートパラメータをバランスした括弧で抽出
func extractRouteParamsBalanced(content string) string {
	var result strings.Builder
	var bracketCount int
	var inQuotes bool
	var quoteChar rune
	var escapeNext bool

	for _, char := range content {
		if escapeNext {
			result.WriteRune(char)
			escapeNext = false
			continue
		}

		if char == '\\' && inQuotes {
			result.WriteRune(char)
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
			case '[':
				bracketCount++
			case ']':
				bracketCount--
				if bracketCount == 0 {
					// 完全なパラメータ配列が見つかった
					result.WriteRune(char)
					return result.String()
				}
			}
		}

		result.WriteRune(char)
	}

	return ""
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

func replaceFormText(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::text\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::text\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormInput("text", params)
			}
			return match
		})
	}
	return text
}

func replaceFormFile(text string) string {
	// {!! Form::file(...) !!} と {{ Form::file(...) }} の両方に対応
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::file\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::file\(\s*(.*?)\s*\)\s*\}\}`,
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

	re := regexCache.GetRegex(`"([^"\\]+)"`)

	for re.MatchString(result) {
		result = re.ReplaceAllString(result, "'$1'")
	}

	return result
}

// convertEventHandlerQuotesInHTML HTML出力内のイベントハンドラ属性のダブルクォートを柔軟に変換
func convertEventHandlerQuotesInHTML(html string) string {
	result := html

	onclickPattern := `(?i)(onclick=")([^"]*(?:"[^"]*")*?[^"]*?)("(?:\s|>))`
	onclickRe := regexCache.GetRegex(onclickPattern)

	result = onclickRe.ReplaceAllStringFunc(result, func(match string) string {
		matches := onclickRe.FindStringSubmatch(match)
		if len(matches) >= 4 {
			prefix := matches[1] // 元の大文字小文字を保持
			jsCode := matches[2]
			suffix := matches[3]
			convertedJS := convertJavaScriptStringLiterals(jsCode)
			return prefix + convertedJS + suffix
		}
		return match
	})

	onchangePattern := `(?i)(onchange=")([^"]*(?:"[^"]*")*?[^"]*?)("(?:\s|>))`
	onchangeRe := regexCache.GetRegex(onchangePattern)

	result = onchangeRe.ReplaceAllStringFunc(result, func(match string) string {
		matches := onchangeRe.FindStringSubmatch(match)
		if len(matches) >= 4 {
			prefix := matches[1] // 元の大文字小文字を保持
			jsCode := matches[2]
			suffix := matches[3]
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
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\{\s*Form::number\(\s*(.*?)\s*\)\s*\}\}`,
		`(?s)\{\!\!\s*Form::number\(\s*(.*?)\s*\)\s*\!\!\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormNumber(params)
			}
			return match
		})
	}
	return text
}

func processFormNumber(params []string) string {
	if len(params) < 1 {
		return ""
	}

	// PHP文字列連結を含むフィールド名を適切に処理
	name := ProcessFieldName(params[0])
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
			"step":        `'step'\s*=>\s*(\d+(?:\.\d+)?)`,
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

	// PHP文字列連結を含むフィールド名を適切に処理
	name := ProcessFieldName(params[0])
	value := ""
	if len(params) > 1 {
		value = params[1]
	}

	attrProcessor := &AttributeProcessor{
		Order: []string{"placeholder", "class", "id", "required"},
		Patterns: map[string]string{
			"placeholder": `'placeholder'\s*=>\s*'([^']+)'`,
			"class":       `'class'\s*=>\s*'([^']+)'`,
			"id":          `'id'\s*=>\s*'([^']+)'`,
			"required":    `'required'\s*=>\s*'([^']*)'`,
		},
	}

	extraAttrs := ""
	if len(params) > 2 {
		extraAttrs = attrProcessor.ProcessAttributes(params[2])
	}

	formattedValue := FormatValueAttribute(value)
	return fmt.Sprintf(`<input type="%s" name="%s" value="%s"%s>`, inputType, name, formattedValue, extraAttrs)
}

func processFormPassword(params []string) string {
	if len(params) < 1 {
		return ""
	}

	// PHP文字列連結を含むフィールド名を適切に処理
	name := ProcessFieldName(params[0])

	value := ""

	// パスワード用の属性処理（required属性サポートを含む）
	attrProcessor := &AttributeProcessor{
		Order: []string{"placeholder", "class", "id", "required"},
		Patterns: map[string]string{
			"placeholder": `'placeholder'\s*=>\s*'([^']+)'`,
			"class":       `'class'\s*=>\s*'([^']+)'`,
			"id":          `'id'\s*=>\s*'([^']+)'`,
			"required":    `'required'\s*=>\s*'([^']*)'`,
		},
	}

	extraAttrs := ""
	if len(params) > 1 {
		extraAttrs = attrProcessor.ProcessAttributes(params[1])
	}

	return fmt.Sprintf(`<input type="password" name="%s" value="%s"%s>`, name, value, extraAttrs)
}

func replaceFormInput(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::input\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::input\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormInputDynamic(params)
			}
			return match
		})
	}
	return text
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

func replaceFormEmail(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::email\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::email\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormInput("email", params)
			}
			return match
		})
	}
	return text
}

func replaceFormPassword(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::password\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::password\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormPassword(params)
			}
			return match
		})
	}
	return text
}

func replaceFormUrl(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::url\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::url\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				params := extractParamsBalanced(paramStr)
				return processFormInput("url", params)
			}
			return match
		})
	}
	return text
}

func replaceFormTel(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::tel\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::tel\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				params := extractParamsBalanced(paramStr)
				return processFormInput("tel", params)
			}
			return match
		})
	}
	return text
}

func replaceFormSearch(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::search\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::search\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				params := extractParamsBalanced(paramStr)
				return processFormInput("search", params)
			}
			return match
		})
	}
	return text
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

func replaceFormRange(text string) string {
	// 複雑なネストに対応したパターンに変更
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::range\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::range\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				paramStr := fullMatch[1]
				// バランスを考慮したパラメータ抽出に変更
				params := extractParamsBalanced(paramStr)
				return processFormInput("range", params)
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
