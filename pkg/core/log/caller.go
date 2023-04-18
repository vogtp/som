package log

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	srcPrefix    string
	srcPrefixLen int
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return
	}
	idx := strings.Index(filename, `som/pkg/core/log/caller.go`)
	if idx > len(filename) || idx < 0 {
		return
	}
	srcPrefix = filename[:idx]
	srcPrefixLen = len(srcPrefix)
}

func TrimPackagePath(s string) string {
	if len(s) > len(srcPrefix) && strings.HasPrefix(s, srcPrefix) {
		return s[srcPrefixLen:]
	}
	return s
}

func Caller(skip int) string {
	_, filename, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", TrimPackagePath(filename), line)
}
