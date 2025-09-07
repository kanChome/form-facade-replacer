package ffr

import "fmt"

func replaceFormSelect(text string) string {
    patterns := []string{
        `(?s)\{\{\s*Form::select\(\s*(.*?)\s*\)\s*\}\}`,
        `(?s)\{\!\!\s*Form::select\(\s*(.*?)\s*\)\s*\!\!\}`,
    }
    for _, pattern := range patterns {
        re := regexCache.GetRegex(pattern)
        text = re.ReplaceAllStringFunc(text, func(match string) string {
            fullMatch := re.FindStringSubmatch(match)
            if len(fullMatch) > 1 {
                params := extractParamsBalanced(fullMatch[1])
                return processFormSelect(params)
            }
            return match
        })
    }
    return text
}

func processFormSelect(params []string) string {
    if len(params) < 2 { return "" }
    name := ProcessFieldName(params[0])
    options := params[1]
    selected := ""
    if len(params) > 2 { selected = params[2] }
    attrProcessor := &AttributeProcessor{
        Order: []string{"class", "id", "onchange"},
        Patterns: map[string]string{
            "class":    `'class'\s*=>\s*'([^']+)'`,
            "id":       `'id'\s*=>\s*'([^']+)'`,
            "onchange": `'(?:onChange|onchange)'\s*=>\s*'([^']+)'`,
        },
    }
    extraAttrs := ""
    if len(params) > 3 { extraAttrs = attrProcessor.ProcessAttributes(params[3]) }
    return fmt.Sprintf(`<select name="%s"%s>
@foreach(%s as $key => $value)
<option value="{{ $key }}" @if($key == %s) selected @endif>{{ $value }}</option>
@endforeach
</select>`, name, extraAttrs, options, selected)
}

