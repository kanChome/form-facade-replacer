package ffr

import (
	"fmt"
	"strings"
)

func replaceFormCheckbox(text string) string {
	patterns := []string{
		`(?s)\{\{\s*Form::checkbox\(\s*(.*?)\s*\)\s*\}\}`,
		`(?s)\{\!\!\s*Form::checkbox\(\s*(.*?)\s*\)\s*\!\!\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				params := extractParamsBalanced(fullMatch[1])
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
	name := ProcessFieldName(params[0])
	value := ""
	if len(params) > 1 {
		value = strings.Trim(params[1], `'"`)
	}
	checked := ""
	if len(params) > 2 {
		checked = params[2]
	}
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
		dataRe := regexCache.GetRegex(`'(data-[^']+)'\s*=>\s*'([^']+)'`)
		for _, match := range dataRe.FindAllStringSubmatch(params[3], -1) {
			extraAttrs += fmt.Sprintf(` %s="%s"`, match[1], match[2])
		}
		extraAttrs += attrProcessor.ProcessAttributes(params[3])
	}
	var result string
	if strings.HasSuffix(name, "[]") {
		result = fmt.Sprintf(`<input type="checkbox" name="%s" value="{{ %s }}" @if(in_array(%s, (array)%s)) checked @endif%s>`, name, value, value, checked, extraAttrs)
	} else {
		result = fmt.Sprintf(`<input type="checkbox" name="%s" value="{{ %s }}" @if(%s) checked @endif%s>`, name, value, checked, extraAttrs)
	}
	result = convertEventHandlerQuotesInHTML(result)
	return result
}
