package ffr

import (
	"fmt"
)

func replaceFormText(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::text\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::text\(\s*(.*?)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				params := extractParamsBalanced(fullMatch[1])
				return processFormInput("text", params)
			}
			return match
		})
	}
	return text
}

func replaceFormEmail(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::email\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::email\(\s*(.*?)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				params := extractParamsBalanced(fullMatch[1])
				return processFormInput("email", params)
			}
			return match
		})
	}
	return text
}

func replaceFormPassword(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::password\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::password\(\s*(.*?)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				params := extractParamsBalanced(fullMatch[1])
				return processFormPassword(params)
			}
			return match
		})
	}
	return text
}

func replaceFormUrl(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::url\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::url\(\s*(.*?)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				params := extractParamsBalanced(fullMatch[1])
				return processFormInput("url", params)
			}
			return match
		})
	}
	return text
}

func replaceFormTel(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::tel\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::tel\(\s*(.*?)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				params := extractParamsBalanced(fullMatch[1])
				return processFormInput("tel", params)
			}
			return match
		})
	}
	return text
}

func replaceFormSearch(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::search\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::search\(\s*(.*?)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				params := extractParamsBalanced(fullMatch[1])
				return processFormInput("search", params)
			}
			return match
		})
	}
	return text
}

// 動的な Form::input(type, name, value, attrs)
func replaceFormInput(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::input\(\s*(.*?)\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::input\(\s*(.*?)\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			fullMatch := re.FindStringSubmatch(match)
			if len(fullMatch) > 1 {
				params := extractParamsBalanced(fullMatch[1])
				return processFormInputDynamic(params)
			}
			return match
		})
	}
	return text
}

func processFormInput(inputType string, params []string) string {
	if len(params) < 1 {
		return ""
	}
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
	name := ProcessFieldName(params[0])
	value := ""
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
