package ffr

import (
    "fmt"
    "strings"
)

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
    dataRe := regexCache.GetRegex(`'(data-[^']+)'\s*=>\s*'([^']+)'`)
    for _, match := range dataRe.FindAllStringSubmatch(attrs, -1) {
        extraAttrs += fmt.Sprintf(` %s="%s"`, match[1], match[2])
    }
    extraAttrs += processDynamicAttributes(attrs)
    return fmt.Sprintf(`<button%s>{!! %s !!}</button>`, extraAttrs, textParam)
}

func replaceFormSubmit(text string) string {
    patterns := []string{
        `(?is)\{\!\!\s*Form::submit\(\s*(.*?)\s*\)\s*\!\!\}`,
        `(?is)\{\{\s*Form::submit\(\s*(.*?)\s*\)\s*\}\}`,
    }
    for _, pattern := range patterns {
        re := regexCache.GetRegex(pattern)
        text = re.ReplaceAllStringFunc(text, func(match string) string {
            fullMatch := re.FindStringSubmatch(match)
            if len(fullMatch) > 1 {
                params := extractParamsBalanced(fullMatch[1])
                return processFormSubmit(params)
            }
            return match
        })
    }
    return text
}

func processFormSubmit(params []string) string {
    if len(params) < 1 { return "" }
    textParam := strings.Trim(params[0], `'"`)
    if textParam == "null" { textParam = "" }
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
    if len(params) > 1 { extraAttrs = attrProcessor.ProcessAttributes(params[1]) }
    return fmt.Sprintf(`<button type="submit"%s>%s</button>`, extraAttrs, textParam)
}
