package ffr

import (
    "fmt"
    "strings"
)

func replaceFormNumber(text string) string {
    patterns := []string{
        `(?s)\{\{\s*Form::number\(\s*(.*?)\s*\)\s*\}\}`,
        `(?s)\{\!\!\s*Form::number\(\s*(.*?)\s*\)\s*\!\!\}`,
    }
    for _, pattern := range patterns {
        re := regexCache.GetRegex(pattern)
        text = re.ReplaceAllStringFunc(text, func(match string) string {
            fullMatch := re.FindStringSubmatch(match)
            if len(fullMatch) > 1 {
                params := extractParamsBalanced(fullMatch[1])
                return processFormNumber(params)
            }
            return match
        })
    }
    return text
}

func processFormNumber(params []string) string {
    if len(params) < 1 { return "" }
    name := ProcessFieldName(params[0])
    value := ""
    if len(params) > 1 { value = params[1] }
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
    if len(params) > 2 { extraAttrs = attrProcessor.ProcessAttributes(params[2]) }
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

func replaceFormRange(text string) string {
    patterns := []string{
        `(?s)\{\!\!\s*Form::range\(\s*(.*?)\s*\)\s*\!\!\}`,
        `(?s)\{\{\s*Form::range\(\s*(.*?)\s*\)\s*\}\}`,
    }
    for _, pattern := range patterns {
        re := regexCache.GetRegex(pattern)
        text = re.ReplaceAllStringFunc(text, func(match string) string {
            fullMatch := re.FindStringSubmatch(match)
            if len(fullMatch) > 1 {
                params := extractParamsBalanced(fullMatch[1])
                return processFormInput("range", params)
            }
            return match
        })
    }
    return text
}

