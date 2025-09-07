package ffr

import "fmt"

func replaceFormLabel(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::label\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::label\(\s*(.*?)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				params := extractParamsBalanced(fullMatch[1])
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
	name := ProcessFieldName(params[0])
	forAttr := name
	textParam := ""
	if len(params) == 1 {
		textParam = fmt.Sprintf("'%s'", name)
	} else {
		textParam = params[1]
	}
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
		forRe := regexCache.GetRegex(`'for'\s*=>\s*'([^']+)'`)
		if forRe.MatchString(attrs) {
			forAttr = forRe.FindStringSubmatch(attrs)[1]
		}
		extraAttrs = attrProcessor.ProcessAttributes(attrs)
	}
	return fmt.Sprintf(`<label for="%s"%s>{!! %s !!}</label>`, forAttr, extraAttrs, textParam)
}
