package ffr

import (
	"fmt"
	"strings"
)

func replaceFormRadio(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::radio\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::radio\(\s*(.*?)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				params := extractParamsBalanced(fullMatch[1])
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
	name := ProcessFieldName(params[0])
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
	checkedAttr := ""
	if checked != "" && checked != "false" && checked != "null" {
		checkedAttr = fmt.Sprintf(" @if(%s) checked @endif", checked)
	}
	return fmt.Sprintf(`<input type="radio" name="%s" value="%s"%s%s>`, name, value, checkedAttr, extraAttrs)
}
