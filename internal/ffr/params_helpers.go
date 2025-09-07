package ffr

import (
	"strings"
)

// パラメータ抽出（バランス括弧対応）
func extractParamsBalanced(paramsStr string) []string {
	var params []string
	var current strings.Builder
	depth := 0
	inQuotes := false
	var quoteChar rune
	escape := false

	for _, ch := range paramsStr {
		if escape {
			current.WriteRune(ch)
			escape = false
			continue
		}
		if ch == '\\' && inQuotes {
			current.WriteRune(ch)
			escape = true
			continue
		}
		if (ch == '\'' || ch == '"') && !inQuotes {
			inQuotes = true
			quoteChar = ch
		} else if inQuotes && ch == quoteChar {
			inQuotes = false
		}
		if !inQuotes {
			switch ch {
			case '(', '[', '{':
				depth++
			case ')', ']', '}':
				depth--
			case ',':
				if depth == 0 {
					params = append(params, strings.TrimSpace(current.String()))
					current.Reset()
					continue
				}
			}
		}
		current.WriteRune(ch)
	}
	if strings.TrimSpace(current.String()) != "" {
		params = append(params, strings.TrimSpace(current.String()))
	}
	return params
}

// 高度な抽出（フォールバック用）
func extractParamsAdvanced(paramsStr string) []string {
	return extractParamsBalanced(paramsStr)
}

func extractParams(paramsStr string) []string {
	return extractParamsBalanced(paramsStr)
}

// JS文字列リテラル/イベント属性の一部変換
func convertJavaScriptStringLiterals(jsCode string) string {
	result := jsCode
	re := regexCache.GetRegex(`\"([^\"\\]+)\"`)
	for re.MatchString(result) {
		result = re.ReplaceAllString(result, "'$1'")
	}
	return result
}

func convertEventHandlerQuotesInHTML(html string) string {
	result := html
	onclickPattern := `(?i)(onclick=\")([^\"]*(?:\"[^\"]*\")*?[^\"]*?)(\"(?:\s|>))`
	onclickRe := regexCache.GetRegex(onclickPattern)
	result = onclickRe.ReplaceAllStringFunc(result, func(match string) string {
		matches := onclickRe.FindStringSubmatch(match)
		if len(matches) >= 4 {
			prefix := matches[1]
			jsCode := matches[2]
			suffix := matches[3]
			convertedJS := convertJavaScriptStringLiterals(jsCode)
			return prefix + convertedJS + suffix
		}
		return match
	})

	onchangePattern := `(?i)(onchange=\")([^\"]*(?:\"[^\"]*\")*?[^\"]*?)(\"(?:\s|>))`
	onchangeRe := regexCache.GetRegex(onchangePattern)
	result = onchangeRe.ReplaceAllStringFunc(result, func(match string) string {
		matches := onchangeRe.FindStringSubmatch(match)
		if len(matches) >= 4 {
			prefix := matches[1]
			jsCode := matches[2]
			suffix := matches[3]
			convertedJS := convertJavaScriptStringLiterals(jsCode)
			return prefix + convertedJS + suffix
		}
		return match
	})
	return result
}
