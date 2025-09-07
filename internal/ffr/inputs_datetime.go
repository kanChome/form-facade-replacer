package ffr

func replaceFormDate(text string) string {
    patterns := []string{
        `(?s)\{\!\!\s*Form::date\(\s*(.*?)\s*\)\s*\!\!\}`,
        `(?s)\{\{\s*Form::date\(\s*(.*?)\s*\)\s*\}\}`,
    }
    for _, pattern := range patterns {
        re := regexCache.GetRegex(pattern)
        text = re.ReplaceAllStringFunc(text, func(match string) string {
            fullMatch := re.FindStringSubmatch(match)
            if len(fullMatch) > 1 {
                params := extractParamsBalanced(fullMatch[1])
                return processFormInput("date", params)
            }
            return match
        })
    }
    return text
}

func replaceFormTime(text string) string {
    patterns := []string{
        `(?s)\{\!\!\s*Form::time\(\s*(.*?)\s*\)\s*\!\!\}`,
        `(?s)\{\{\s*Form::time\(\s*(.*?)\s*\)\s*\}\}`,
    }
    for _, pattern := range patterns {
        re := regexCache.GetRegex(pattern)
        text = re.ReplaceAllStringFunc(text, func(match string) string {
            fullMatch := re.FindStringSubmatch(match)
            if len(fullMatch) > 1 {
                params := extractParamsBalanced(fullMatch[1])
                return processFormInput("time", params)
            }
            return match
        })
    }
    return text
}

func replaceFormDatetime(text string) string {
    patterns := []string{
        `(?s)\{\!\!\s*Form::datetime\(\s*(.*?)\s*\)\s*\!\!\}`,
        `(?s)\{\{\s*Form::datetime\(\s*(.*?)\s*\)\s*\}\}`,
    }
    for _, pattern := range patterns {
        re := regexCache.GetRegex(pattern)
        text = re.ReplaceAllStringFunc(text, func(match string) string {
            fullMatch := re.FindStringSubmatch(match)
            if len(fullMatch) > 1 {
                params := extractParamsBalanced(fullMatch[1])
                return processFormInput("datetime-local", params)
            }
            return match
        })
    }
    return text
}
