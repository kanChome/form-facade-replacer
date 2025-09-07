package ffr

import (
	"fmt"
	"strings"
)

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

func processFormOpen(content string) string {
	action := extractFormAction(content)
	method := extractFormMethod(content)
	extraAttrs := extractFormAttributes(content)
	return buildFormTag(action, method, extraAttrs)
}

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

func extractFormMethod(content string) string {
	methodRe := regexCache.GetRegex(`'method'\s*=>\s*'([^']+)'`)
	if methodRe.MatchString(content) {
		return methodRe.FindStringSubmatch(content)[1]
	}
	return "GET"
}

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

func buildFormTag(action, method, extraAttrs string) string {
	if strings.ToUpper(method) == "GET" {
		return fmt.Sprintf(`<form action="%s" method="%s"%s>`, action, method, extraAttrs)
	}
	return fmt.Sprintf(`<form action="%s" method="%s"%s>
{{ csrf_field() }}`, action, method, extraAttrs)
}

func replaceFormClose(text string) string {
	return ProcessBladePatterns(text, `Form::close\(\)`, func(content string) string {
		return "</form>"
	})
}
