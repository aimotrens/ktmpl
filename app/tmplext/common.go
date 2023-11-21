package tmplext

import (
	"fmt"
	"strings"
)

func indent(spaces int, s string) string {
	indent := strings.Repeat(" ", spaces)
	return indent + strings.ReplaceAll(s, "\n", "\n"+indent)
}

func substr(start, length int, s string) string {
	return s[start : start+length]
}

func iterate(from, to int) []int {
	var result []int
	for i := from; i <= to; i++ {
		result = append(result, i)
	}
	return result
}

func format(format string, obj any) string {
	return fmt.Sprintf(format, obj)
}

func endsWith(suffix, s string) bool {
	return strings.HasSuffix(s, suffix)
}

func startsWith(prefix, s string) bool {
	return strings.HasPrefix(s, prefix)
}

func contains(substring, s string) bool {
	return strings.Contains(s, substring)
}
