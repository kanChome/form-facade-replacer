package ffr

import (
	"fmt"
	"strings"
)

// 属性処理
type AttributeProcessor struct {
	Order    []string
	Patterns map[string]string
}

func (ap *AttributeProcessor) ProcessAttributes(attrs string) string {
	var extraAttrs string
	for _, attr := range ap.Order {
		if pattern, exists := ap.Patterns[attr]; exists {
			if re := regexCache.GetRegex(pattern); re.MatchString(attrs) {
				matches := re.FindStringSubmatch(attrs)
				var val string
				if len(matches) > 2 && matches[2] != "" {
					val = matches[2]
				} else {
					val = matches[1]
				}
				val = processAttributeValue(val)
				if (attr == "disabled" && (val == "" || val == "disabled")) ||
					(attr == "required" && (val == "" || val == "required")) {
					extraAttrs += fmt.Sprintf(` %s`, attr)
				} else {
					extraAttrs += fmt.Sprintf(` %s="%s"`, attr, val)
				}
			}
		}
	}
	return extraAttrs
}

// Bladeパターン適用
func ProcessBladePatterns(text, formMethod string, processor func(string) string) string {
	patterns := []string{
		fmt.Sprintf(BladeExclamationPattern, formMethod),
		fmt.Sprintf(BladeCurlyPattern, formMethod),
	}
	for _, pattern := range patterns {
		re := regexCache.GetRegex(pattern)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			matches := re.FindStringSubmatch(match)
			if len(matches) > 1 {
				return processor(matches[1])
			}
			return processor("")
		})
	}
	return text
}

// 属性値処理
func processAttributeValue(value string) string {
	originalValue := value
	value = strings.Trim(value, `'"`)
	if strings.Contains(value, " . ") {
		concatPattern1 := `^'([^']*)'\s*\.\s*(.+?)\s*\.\s*'([^']*)'$`
		re1 := regexCache.GetRegex(concatPattern1)
		if matches := re1.FindStringSubmatch(value); len(matches) == 4 {
			prefix := matches[1]
			variable := strings.TrimSpace(matches[2])
			suffix := matches[3]
			return fmt.Sprintf("{{ '%s' . %s . '%s' }}", prefix, variable, suffix)
		}
		concatPattern2 := `^'([^']*)'\s*\.\s*(.+)$`
		re2 := regexCache.GetRegex(concatPattern2)
		if matches := re2.FindStringSubmatch(value); len(matches) == 3 {
			prefix := matches[1]
			variable := strings.TrimSpace(matches[2])
			return fmt.Sprintf("{{ '%s' . %s }}", prefix, variable)
		}
		concatPattern3 := `^'([^']*)'\s*\.\s*(\$[a-zA-Z_]\w*(?:\[[^\]]*\])*)\s*$`
		re3 := regexCache.GetRegex(concatPattern3)
		if matches := re3.FindStringSubmatch(value); len(matches) == 3 {
			prefix := matches[1]
			variable := strings.TrimSpace(matches[2])
			return fmt.Sprintf("{{ '%s' . %s }}", prefix, variable)
		}
		return fmt.Sprintf("{{ %s }}", originalValue)
	}
	return value
}

// 値の整形
func DetectArrayHelper(value string) bool {
	return regexCache.GetRegex(`(?i)^(old|session|request|input)\s*\(`).MatchString(strings.TrimSpace(value))
}

func IsArrayFieldName(fieldName string) bool {
	return regexCache.GetRegex(`\[.*\]`).MatchString(fieldName)
}

func ProcessFieldName(name string) string {
	nameAttr := strings.Trim(name, `'"`)
	if strings.Contains(nameAttr, " . ") {
		concatPattern1 := `^'([^']*)'\s*\.\s*(.+?)\s*\.\s*'([^']*)'$`
		re1 := regexCache.GetRegex(concatPattern1)
		if matches := re1.FindStringSubmatch(nameAttr); len(matches) == 4 {
			prefix := matches[1]
			variable := strings.TrimSpace(matches[2])
			suffix := matches[3]
			return fmt.Sprintf("%s{{ %s }}%s", prefix, variable, suffix)
		}
		concatPattern2 := `^([^']*\[)'\s*\.\s*(.+?)\s*\.\s*'(\][^']*)$`
		re2 := regexCache.GetRegex(concatPattern2)
		if matches := re2.FindStringSubmatch(nameAttr); len(matches) == 4 {
			prefix := matches[1]
			variable := strings.TrimSpace(matches[2])
			suffix := matches[3]
			return fmt.Sprintf("%s{{ %s }}%s", prefix, variable, suffix)
		}
	}
	return nameAttr
}

func FormatValueAttribute(value string) string {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" || trimmedValue == "null" || trimmedValue == "''" || trimmedValue == `""` {
		return ""
	}
	if strings.HasPrefix(trimmedValue, "'") && strings.HasSuffix(trimmedValue, "'") {
		innerValue := strings.Trim(trimmedValue, "'")
		if regexCache.GetRegex(`^\d+(\.\d+)?$`).MatchString(innerValue) {
			return fmt.Sprintf("{{ %s }}", innerValue)
		}
		if regexCache.GetRegex(`^#[0-9a-fA-F]{3,6}$`).MatchString(innerValue) {
			return fmt.Sprintf("{{ %s }}", innerValue)
		}
	}
	return fmt.Sprintf("{{ %s }}", value)
}

func FormatHiddenValueAttribute(value string, fieldName string) string {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" || trimmedValue == "null" || trimmedValue == "''" || trimmedValue == `""` {
		return ""
	}
	if IsArrayFieldName(fieldName) {
		return fmt.Sprintf("{{ is_array(%s) ? implode(',', %s) : %s }}", value, value, value)
	}
	return fmt.Sprintf("{{ %s }}", value)
}
