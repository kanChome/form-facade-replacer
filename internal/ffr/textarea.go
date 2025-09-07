// textarea.go: テキストエリア要素の置換ロジック。
package ffr

import (
	"fmt"
	"strings"
)

// --- Textarea ---
// replaceFormTextarea は Blade 内の Form::textarea(...) を HTML に置換する。
func replaceFormTextarea(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::textarea\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::textarea\(\s*(.*?)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				params := extractParamsBalanced(fullMatch[1])
				return processFormTextarea(params)
			}
			return match
		})
	}
	return text
}

// processFormTextarea は rows/placeholder など属性と値を整形し最終HTMLを生成する。
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
	if strings.TrimSpace(value) == "" || strings.TrimSpace(value) == "null" || value == "''" || value == `""` {
		return fmt.Sprintf(`<textarea name="%s"%s></textarea>`, name, extraAttrs)
	}
	formattedValue := FormatValueAttribute(value)
	if formattedValue == "" {
		return fmt.Sprintf(`<textarea name="%s"%s></textarea>`, name, extraAttrs)
	}
	return fmt.Sprintf(`<textarea name="%s"%s>%s</textarea>`, name, extraAttrs, formattedValue)
}
