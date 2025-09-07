package ffr

import (
	"regexp"
	"sync"
)

// 正規表現キャッシュ
type RegexCache struct {
	mu    sync.RWMutex
	cache map[string]*regexp.Regexp
}

var regexCache = &RegexCache{
	cache: make(map[string]*regexp.Regexp),
}

// GetRegex 正規表現の取得（キャッシュあり）
func (rc *RegexCache) GetRegex(pattern string) *regexp.Regexp {
	rc.mu.RLock()
	if re, exists := rc.cache[pattern]; exists {
		rc.mu.RUnlock()
		return re
	}
	rc.mu.RUnlock()

	rc.mu.Lock()
	defer rc.mu.Unlock()

	// ダブルチェック（他のゴルーチンが既に作成している可能性）
	if re, exists := rc.cache[pattern]; exists {
		return re
	}

	re := regexp.MustCompile(pattern)
	rc.cache[pattern] = re
	return re
}

// Bladeの基本パターン
const (
	BladeExclamationPattern = `\{\!\!\s*%s\s*\!\!\}`
	BladeCurlyPattern       = `\{\{\s*%s\s*\}\}`
)
