package utils

import (
	"fmt"
	"strings"
)

type Mode byte

const (
	EQ    = Mode('E') // 完全匹配
	Match = Mode('M') // 正则匹配
)

func UriKey(method, path string, mode Mode) string {
	return fmt.Sprintf("URI|%c|%s|%s", mode, strings.ToTitle(method), path)
}
