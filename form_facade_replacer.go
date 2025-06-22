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
)

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

	// コマンドライン引数の処理
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

	// ファイルかディレクトリかを判定
	info, err := os.Stat(config.TargetPath)
	if err != nil {
		log.Fatalf("エラー: '%s' が存在しません。", config.TargetPath)
	}

	config.IsFile = !info.IsDir()

	if config.IsFile {
		// 単一ファイルの場合、.blade.phpかチェック
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
	fmt.Println(" go run form_facade_replacer.go resources/views/web/hoge/fuga.blade.php")
}

func processBladeFiles(config *ReplacementConfig) error {
	if config.IsFile {
		// 単一ファイルの処理
		return processSingleFile(config, config.TargetPath)
	}

	// ディレクトリの再帰処理
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

func replaceFormPatterns(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	text := string(content)

	// 1. Form::open() の置換
	text = replaceFormOpen(text)

	// 2. Form::close() の置換
	text = replaceFormClose(text)

	// 3. Form::hidden() の置換
	text = replaceFormHidden(text)

	// 4. Form::button() の置換
	text = replaceFormButton(text)

	// 5. Form::textarea() の置換
	text = replaceFormTextarea(text)

	// 6. Form::label() の置換
	text = replaceFormLabel(text)

	// 7. Form::text() の置換
	text = replaceFormText(text)

	// 8. Form::number() の置換
	text = replaceFormNumber(text)

	// 9. Form::select() の置換
	text = replaceFormSelect(text)

	// 10. Form::checkbox() の置換
	text = replaceFormCheckbox(text)

	// 11. Form::Submit() の置換
	text = replaceFormSubmit(text)

	return os.WriteFile(filePath, []byte(text), 0644)
}

func replaceFormOpen(text string) string {
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::open\(\s*\[\s*(.*?)\s*\]\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::open\(\s*\[\s*(.*?)\s*\]\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			return processFormOpen(re.FindStringSubmatch(match)[1])
		})
	}
	return text
}

func processFormOpen(content string) string {
	action := ""
	method := "GET"
	extraAttrs := ""

	// route処理 - 文字列形式と配列形式の両方に対応
	// まず配列形式のroute（パラメータ付き）をチェック
	paramRouteRe := regexp.MustCompile(`'route'\s*=>\s*\[\s*'([^']+)'\s*,\s*(\[[^\]]*\])`)
	if paramMatches := paramRouteRe.FindStringSubmatch(content); len(paramMatches) > 2 {
		// 配列パラメータを含む場合の処理
		action = fmt.Sprintf("{{ route('%s', %s) }}", paramMatches[1], paramMatches[2])
	} else {
		// 次に配列形式のroute（パラメータなし）をチェック
		arrayRouteRe := regexp.MustCompile(`'route'\s*=>\s*\[\s*'([^']+)'\s*\]`)
		if arrayMatches := arrayRouteRe.FindStringSubmatch(content); len(arrayMatches) > 1 {
			action = fmt.Sprintf("{{ route('%s') }}", arrayMatches[1])
		} else {
			// 最後に文字列形式のrouteをチェック
			simpleRouteRe := regexp.MustCompile(`'route'\s*=>\s*'([^']+)'`)
			if simpleMatches := simpleRouteRe.FindStringSubmatch(content); len(simpleMatches) > 1 {
				action = fmt.Sprintf("{{ route('%s') }}", simpleMatches[1])
			}
		}
	}

	// actionが設定されなかった場合のみurl処理を実行
	if action == "" {
		// url処理 - 先読みアサーションを使わない方法に変更
		urlRe := regexp.MustCompile(`'url'\s*=>\s*([^,\]]+)`)
		if matches := urlRe.FindStringSubmatch(content); len(matches) > 1 {
			urlVal := strings.TrimSpace(matches[1])
			// route()関数の場合の特別処理
			if strings.HasPrefix(urlVal, "route(") {
				// 括弧のバランスを考慮してroute()の完全な関数呼び出しを抽出
				routeFuncRe := regexp.MustCompile(`route\([^)]*(?:\([^)]*\)[^)]*)*\)`)
				if routeMatch := routeFuncRe.FindString(content); routeMatch != "" {
					action = fmt.Sprintf("{{ %s }}", routeMatch)
				} else {
					action = fmt.Sprintf("{{ %s }}", urlVal)
				}
			} else {
				action = urlVal
			}
		}
	}

	// method処理
	if methodRe := regexp.MustCompile(`'method'\s*=>\s*'([^']+)'`); methodRe.MatchString(content) {
		method = methodRe.FindStringSubmatch(content)[1]
	}

	// その他の属性処理
	attrPatterns := map[string]string{
		"target": `'target'\s*=>\s*'([^']+)'`,
		"id":     `'id'\s*=>\s*'([^']+)'`,
		"class":  `'class'\s*=>\s*'([^']+)'`,
	}

	for attr, pattern := range attrPatterns {
		if re := regexp.MustCompile(pattern); re.MatchString(content) {
			value := re.FindStringSubmatch(content)[1]
			extraAttrs += fmt.Sprintf(` %s="%s"`, attr, value)
		}
	}

	// GETメソッドの場合はCSRF fieldを含めない
	if strings.ToUpper(method) == "GET" {
		return fmt.Sprintf(`<form action="%s" method="%s"%s>`, action, method, extraAttrs)
	} else {
		return fmt.Sprintf(`<form action="%s" method="%s"%s>
{{ csrf_field() }}`, action, method, extraAttrs)
	}
}

func replaceFormClose(text string) string {
	patterns := []string{
		`\{\!\!\s*Form::close\(\)\s*\!\!\}`,
		`\{\{\s*Form::close\(\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllString(text, "</form>")
	}
	return text
}

func replaceFormHidden(text string) string {
	// マルチライン形式に対応するため(?s)フラグを追加
	patterns := []string{
		`(?s)\{\!\!\s*Form::hidden\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::hidden\(\s*(.*?)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			// 括弧内のコンテンツを抽出
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

	// 文字列結合が含まれる場合は、Bladeの変数展開形式に変換
	nameAttr := name
	if strings.Contains(name, " . ") {
		// 'reason[' . $contentsData['statusList']['CHANGE'] . ']' のようなパターンを
		// reason[{{ $contentsData['statusList']['CHANGE'] }}] に変換
		// 複数のパターンを試行
		patterns := []string{
			// シングルクォート区切りのパターン
			`^'([^']*)'\\s*\\.\\s*([^'\\s]+(?:\\[[^\\]]*\\]\\[[^\\]]*\\])?)\\s*\\.\\s*'([^']*)'$`,
			// より一般的なパターン
			`^'([^']*)'\\s*\\.\\s*(.+?)\\s*\\.\\s*'([^']*)'$`,
		}

		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			if matches := re.FindStringSubmatch(name); len(matches) == 4 {
				nameAttr = fmt.Sprintf("%s{{ %s }}%s", matches[1], matches[2], matches[3])
				break
			}
		}
	}

	extraAttrs := ""
	if len(params) > 2 {
		attrs := params[2]
		if idRe := regexp.MustCompile(`'id'\s*=>\s*'([^']+)'`); idRe.MatchString(attrs) {
			extraAttrs += fmt.Sprintf(` id="%s"`, idRe.FindStringSubmatch(attrs)[1])
		}
		if classRe := regexp.MustCompile(`'class'\s*=>\s*'([^']+)'`); classRe.MatchString(attrs) {
			extraAttrs += fmt.Sprintf(` class="%s"`, classRe.FindStringSubmatch(attrs)[1])
		}
	}

	return fmt.Sprintf(`<input type="hidden" name="%s" value="{{ %s }}"%s>`, nameAttr, value, extraAttrs)
}

func replaceFormButton(text string) string {
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\!\!\s*Form::button\(\s*(.*?)\s*,\s*\[\s*(.*?)\s*\]\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::button\(\s*(.*?)\s*,\s*\[\s*(.*?)\s*\]\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			matches := re.FindStringSubmatch(match)
			return processFormButton(matches[1], matches[2])
		})
	}
	return text
}

func processFormButton(textParam, attrs string) string {
	extraAttrs := ""
	attrPatterns := map[string]string{
		"type":    `'type'\s*=>\s*'([^']+)'`,
		"onclick": `'onclick'\s*=>\s*'([^']+)'`,
		"class":   `'class'\s*=>\s*'([^']+)'`,
		"id":      `'id'\s*=>\s*'([^']+)'`,
	}

	for attr, pattern := range attrPatterns {
		if re := regexp.MustCompile(pattern); re.MatchString(attrs) {
			value := re.FindStringSubmatch(attrs)[1]
			extraAttrs += fmt.Sprintf(` %s="%s"`, attr, value)
		}
	}

	// data属性の処理
	dataRe := regexp.MustCompile(`'(data-[^']+)'\s*=>\s*'([^']+)'`)
	for _, match := range dataRe.FindAllStringSubmatch(attrs, -1) {
		extraAttrs += fmt.Sprintf(` %s="%s"`, match[1], match[2])
	}

	return fmt.Sprintf(`<button%s>{!! %s !!}</button>`, extraAttrs, textParam)
}

func replaceFormTextarea(text string) string {
	// (?s)フラグで改行を含む文字列のマッチを有効化
	patterns := []string{
		`(?s)\{\{\s*Form::textarea\(\s*'([^']+)'\s*,\s*([^,]*(?:\([^)]*\)[^,]*)*[^,]*)\s*,\s*\[(.*?)\]\s*\)\s*\}\}`,
		`(?s)\{\!\!\s*Form::textarea\(\s*'([^']+)'\s*,\s*([^,]*(?:\([^)]*\)[^,]*)*[^,]*)\s*,\s*\[(.*?)\]\s*\)\s*\!\!\}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			matches := re.FindStringSubmatch(match)
			return processFormTextarea(matches[1], matches[2], matches[3])
		})
	}
	return text
}

func processFormTextarea(name, value, attrs string) string {
	extraAttrs := ""
	attrPatterns := map[string]string{
		"cols":        `'cols'\s*=>\s*(\d+)`,
		"rows":        `'rows'\s*=>\s*(?:'([^']+)'|(\d+))`,
		"placeholder": `'placeholder'\s*=>\s*'([^']+)'`,
		"class":       `'class'\s*=>\s*'([^']+)'`,
	}

	for attr, pattern := range attrPatterns {
		if re := regexp.MustCompile(pattern); re.MatchString(attrs) {
			matches := re.FindStringSubmatch(attrs)
			var val string
			if len(matches) > 2 && matches[2] != "" {
				val = matches[2] // 数値の場合
			} else {
				val = matches[1] // 文字列の場合
			}
			extraAttrs += fmt.Sprintf(` %s="%s"`, attr, val)
		}
	}

	return fmt.Sprintf(`<textarea name="%s"%s>{{ %s }}</textarea>`, name, extraAttrs, value)
}

func replaceFormLabel(text string) string {
	patterns := []string{
		`\{\!\!\s*Form::label\(\s*([^}]+)\s*\)\s*\!\!\}`,
		`\{\{\s*Form::label\(\s*([^}]+)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
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

	// パラメータ数に応じた処理
	if len(params) == 1 {
		textParam = fmt.Sprintf("'%s'", name)
	} else {
		textParam = params[1]
	}

	extraAttrs := ""
	if len(params) > 2 {
		attrs := params[2]

		if forRe := regexp.MustCompile(`'for'\s*=>\s*'([^']+)'`); forRe.MatchString(attrs) {
			forAttr = forRe.FindStringSubmatch(attrs)[1]
		}

		attrPatterns := map[string]string{
			"class": `'class'\s*=>\s*'([^']+)'`,
			"id":    `'id'\s*=>\s*'([^']+)'`,
			"style": `'style'\s*=>\s*'([^']+)'`,
		}

		for attr, pattern := range attrPatterns {
			if re := regexp.MustCompile(pattern); re.MatchString(attrs) {
				value := re.FindStringSubmatch(attrs)[1]
				extraAttrs += fmt.Sprintf(` %s="%s"`, attr, value)
			}
		}
	}

	return fmt.Sprintf(`<label for="%s"%s>{!! %s !!}</label>`, forAttr, extraAttrs, textParam)
}

func replaceFormText(text string) string {
	patterns := []string{
		`\{\!\!\s*Form::text\(\s*([^}]+)\s*\)\s*\!\!\}`,
		`\{\{\s*Form::text\(\s*([^}]+)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			params := extractParams(re.FindStringSubmatch(match)[1])
			return processFormInput("text", params)
		})
	}
	return text
}

func replaceFormNumber(text string) string {
	patterns := []string{
		`\{\{\s*Form::number\(\s*([^}]+)\s*\)\s*\}\}`,
		`\{\!\!\s*Form::number\(\s*([^}]+)\s*\)\s*\!\!\}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
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

	extraAttrs := ""
	if len(params) > 2 {
		attrs := params[2]
		attrPatterns := map[string]string{
			"placeholder": `'placeholder'\s*=>\s*'([^']+)'`,
			"class":       `'class'\s*=>\s*'([^']+)'`,
			"id":          `'id'\s*=>\s*'([^']+)'`,
			"min":         `'min'\s*=>\s*(\d+)`,
			"max":         `'max'\s*=>\s*(\d+)`,
			"step":        `'step'\s*=>\s*(\d+)`,
		}

		for attr, pattern := range attrPatterns {
			if re := regexp.MustCompile(pattern); re.MatchString(attrs) {
				val := re.FindStringSubmatch(attrs)[1]
				extraAttrs += fmt.Sprintf(` %s="%s"`, attr, val)
			}
		}
	}

	return fmt.Sprintf(`<input type="number" name="%s" value="{{ %s }}"%s>`, name, value, extraAttrs)
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

	extraAttrs := ""
	if len(params) > 2 {
		attrs := params[2]
		attrPatterns := map[string]string{
			"placeholder": `'placeholder'\s*=>\s*'([^']+)'`,
			"class":       `'class'\s*=>\s*'([^']+)'`,
			"id":          `'id'\s*=>\s*'([^']+)'`,
		}

		for attr, pattern := range attrPatterns {
			if re := regexp.MustCompile(pattern); re.MatchString(attrs) {
				val := re.FindStringSubmatch(attrs)[1]
				extraAttrs += fmt.Sprintf(` %s="%s"`, attr, val)
			}
		}
	}

	return fmt.Sprintf(`<input type="%s" name="%s" value="{{ %s }}"%s>`, inputType, name, value, extraAttrs)
}

func replaceFormSelect(text string) string {
	patterns := []string{
		`\{\{\s*Form::select\(\s*([^}]+)\s*\)\s*\}\}`,
		`\{\!\!\s*Form::select\(\s*([^}]+)\s*\)\s*\!\!\}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
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

	extraAttrs := ""
	if len(params) > 3 {
		attrs := params[3]
		attrPatterns := map[string]string{
			"class":    `'class'\s*=>\s*'([^']+)'`,
			"id":       `'id'\s*=>\s*'([^']+)'`,
			"onChange": `'onChange'\s*=>\s*'([^']+)'`,
			"onchange": `'onchange'\s*=>\s*'([^']+)'`,
		}

		for attr, pattern := range attrPatterns {
			if re := regexp.MustCompile(pattern); re.MatchString(attrs) {
				value := re.FindStringSubmatch(attrs)[1]
				// onChangeとonchangeは同じ属性として扱う
				if attr == "onChange" || attr == "onchange" {
					extraAttrs += fmt.Sprintf(` onchange="%s"`, value)
				} else {
					extraAttrs += fmt.Sprintf(` %s="%s"`, attr, value)
				}
			}
		}
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
		re := regexp.MustCompile(pattern)
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

	// 配列形式の名前（例：status[]）かどうかをチェック
	if strings.HasSuffix(name, "[]") {
		// 配列形式の場合、in_array関数を使ってシンプルにチェック
		// 値が配列に含まれていない場合やfalseの場合はチェックされない
		return fmt.Sprintf(`<input type="checkbox" name="%s" value="{{ %s }}" @if(in_array(%s, (array)%s)) checked @endif>`,
			name, value, value, checked)
	} else {
		// 通常のcheckboxの場合
		return fmt.Sprintf(`<input type="checkbox" name="%s" value="{{ %s }}" @if(%s) checked @endif>`, name, value, checked)
	}
}

func replaceFormSubmit(text string) string {
	patterns := []string{
		`\{\!\!\s*Form::submit\(\s*([^}]+)\s*\)\s*\!\!\}`,
		`\{\{\s*Form::submit\(\s*([^}]+)\s*\)\s*\}\}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
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

	extraAttrs := ""
	if len(params) > 1 {
		attrs := params[1]

		attrPatterns := map[string]string{
			"class":    `'class'\s*=>\s*'([^']+)'`,
			"id":       `'id'\s*=>\s*'([^']+)'`,
			"style":    `'style'\s*=>\s*'([^']+)'`,
			"onclick":  `'onclick'\s*=>\s*'([^']+)'`,
			"disabled": `'disabled'\s*=>\s*'([^']*)'`,
		}

		for attr, pattern := range attrPatterns {
			if re := regexp.MustCompile(pattern); re.MatchString(attrs) {
				value := re.FindStringSubmatch(attrs)[1]
				if attr == "disabled" {
					if value == "" || value == "disabled" {
						extraAttrs += " disabled"
					} else {
						extraAttrs += fmt.Sprintf(` %s="%s"`, attr, value)
					}
				} else {
					extraAttrs += fmt.Sprintf(` %s="%s"`, attr, value)
				}
			}
		}
	}

	return fmt.Sprintf(`<button type="submit"%s>%s</button>`, extraAttrs, textParam)
}

func extractParamsAdvanced(paramsStr string) []string {
	// マルチライン形式のパラメータを解析
	// 改行と空白を正規化
	paramsStr = strings.ReplaceAll(paramsStr, "\n", " ")
	paramsStr = strings.ReplaceAll(paramsStr, "\t", " ")
	// 連続する空白を単一の空白に統一
	re := regexp.MustCompile(`\s+`)
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
