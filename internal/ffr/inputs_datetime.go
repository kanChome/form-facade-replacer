// inputs_datetime.go: 日付/時間/日時入力の置換ロジック。
package ffr

// --- Date ---
// replaceFormDate は Blade 内の Form::date(...) を HTML に置換する。
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

// --- Time ---
// replaceFormTime は Blade 内の Form::time(...) を HTML に置換する。
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

// --- Datetime ---
// replaceFormDatetime は Blade 内の Form::datetime(...) を HTML に置換する。
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
