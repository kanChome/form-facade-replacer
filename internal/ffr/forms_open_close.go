// forms_open_close.go: フォーム開始/終了と関連抽出ヘルパの置換ロジック。
package ffr

import (
	"fmt"
	"strings"
)

// --- Open/Close ---
// replaceFormOpen は Blade 内の Form::open([...]) を <form> に置換する。
func replaceFormOpen(text string) string {
	patterns := []string{
		`(?s)\{\!\!\s*Form::open\(\s*\[\s*(.*?)\s*\]\s*\)\s*\!\!\}`,
		`(?s)\{\{\s*Form::open\(\s*\[\s*(.*?)\s*\]\s*\)\s*\}\}`,
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			return processFormOpen(re.FindStringSubmatch(match)[1])
		})
	}
	return text
}

// processFormOpen は open のオプション（action/method/attrs）を解析して form タグを生成する。
func processFormOpen(content string) string {
	action := extractFormAction(content)
	method := extractFormMethod(content)
	extraAttrs := extractFormAttributes(content)
	return buildFormTag(action, method, extraAttrs)
}

// extractFormAction は route/url 指定から action を抽出する（route を優先）。
func extractFormAction(content string) string {
	paramRouteStart := regexCache.GetRegex(`'route'\s*=>\s*\[\s*'([^']+)'\s*,\s*`)
	if startMatch := paramRouteStart.FindStringSubmatchIndex(content); len(startMatch) > 3 {
		routeName := content[startMatch[2]:startMatch[3]]
		paramStart := startMatch[1]
		params := extractRouteParamsBalanced(content[paramStart:])
		if params != "" {
			return fmt.Sprintf("{{ route('%s', %s) }}", routeName, params)
		}
	}
	arrayRouteRe := regexCache.GetRegex(`'route'\s*=>\s*\[\s*'([^']+)'\s*\]`)
	if arrayMatches := arrayRouteRe.FindStringSubmatch(content); len(arrayMatches) > 1 {
		return fmt.Sprintf("{{ route('%s') }}", arrayMatches[1])
	}
	simpleRouteRe := regexCache.GetRegex(`'route'\s*=>\s*'([^']+)'`)
	if simpleMatches := simpleRouteRe.FindStringSubmatch(content); len(simpleMatches) > 1 {
		return fmt.Sprintf("{{ route('%s') }}", simpleMatches[1])
	}
	return extractFormUrl(content)
}

// extractRouteParamsBalanced は route 配列引数をバランスした括弧で抽出する。
func extractRouteParamsBalanced(content string) string {
	var result strings.Builder
	var bracketCount int
	var inQuotes bool
	var quoteChar rune
	var escapeNext bool
	for _, char := range content {
		if escapeNext {
			result.WriteRune(char)
			escapeNext = false
			continue
		}
		if char == '\\' && inQuotes {
			result.WriteRune(char)
			escapeNext = true
			continue
		}
		if !inQuotes && (char == '"' || char == '\'') {
			inQuotes = true
			quoteChar = char
		} else if inQuotes && char == quoteChar {
			inQuotes = false
			quoteChar = 0
		}
		if !inQuotes {
			switch char {
			case '[':
				bracketCount++
			case ']':
				bracketCount--
				if bracketCount == 0 {
					result.WriteRune(char)
					return result.String()
				}
			}
		}
		result.WriteRune(char)
	}
	return ""
}

// extractFormUrl は open の url オプションを抽出する。
func extractFormUrl(content string) string {
	urlRe := regexCache.GetRegex(`'url'\s*=>\s*([^,\]]+)`)
	if matches := urlRe.FindStringSubmatch(content); len(matches) > 1 {
		urlVal := strings.TrimSpace(matches[1])
		if strings.HasPrefix(urlVal, "route(") {
			routeFuncRe := regexCache.GetRegex(`route\([^)]*(?:\([^)]*\)[^)]*)*\)`)
			if routeMatch := routeFuncRe.FindString(content); routeMatch != "" {
				return fmt.Sprintf("{{ %s }}", routeMatch)
			}
			return fmt.Sprintf("{{ %s }}", urlVal)
		}
		return urlVal
	}
	return ""
}

// extractFormMethod は HTTP メソッドを抽出する（既定は GET）。
func extractFormMethod(content string) string {
	methodRe := regexCache.GetRegex(`'method'\s*=>\s*'([^']+)'`)
	if methodRe.MatchString(content) {
		return methodRe.FindStringSubmatch(content)[1]
	}
	return "GET"
}

// extractFormAttributes は id/class/target などの追加属性を整形する。
func extractFormAttributes(content string) string {
	attrProcessor := &AttributeProcessor{
		Order: []string{"class", "id", "target"},
		Patterns: map[string]string{
			"target": `'target'\s*=>\s*'([^']+)'`,
			"id":     `'id'\s*=>\s*'([^']+)'`,
			"class":  `'class'\s*=>\s*'([^']+)'`,
		},
	}
	return attrProcessor.ProcessAttributes(content)
}

// buildFormTag は method に応じて CSRF を付与しつつ form タグを構築する。
func buildFormTag(action, method, extraAttrs string) string {
	if strings.ToUpper(method) == "GET" {
		return fmt.Sprintf(`<form action="%s" method="%s"%s>`, action, method, extraAttrs)
	}
	return fmt.Sprintf(`<form action="%s" method="%s"%s>
{{ csrf_field() }}`, action, method, extraAttrs)
}

// replaceFormClose は Form::close() を </form> に置換する。
func replaceFormClose(text string) string {
	return ProcessBladePatterns(text, `Form::close\(\)`, func(content string) string {
		return "</form>"
	})
}
