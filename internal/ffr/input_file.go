// input_file.go: ファイル入力要素の置換ロジック。
package ffr

import (
    "fmt"
    "strings"
)

// --- File ---
// replaceFormFile は Blade 内の Form::file(...) を HTML に置換する。
func replaceFormFile(text string) string {
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

// processFormFile は File 要素の属性（accept/multiple/イベント等）を解決して最終HTMLを生成する。
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
		multiplePattern := `'multiple'\s*=>\s*(true|false|\d+)`
		if re := regexCache.GetRegex(multiplePattern); re.MatchString(params[1]) {
			matches := re.FindStringSubmatch(params[1])
			if len(matches) > 1 {
				val := matches[1]
				if val == "true" {
					multipleAttr = " multiple"
				}
			}
		}
		extraAttrs = attrProcessor.ProcessAttributes(params[1])
	}
	result := fmt.Sprintf(`<input type="file" name="%s"%s%s>`, name, extraAttrs, multipleAttr)
	result = convertEventHandlerQuotesInHTML(result)
	return result
}
