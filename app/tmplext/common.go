package tmplext

import "strings"

func indent(spaces int, s string) string {
	indent := strings.Repeat(" ", spaces)
	return indent + strings.ReplaceAll(s, "\n", "\n"+indent)
}

func substr(start, length int, s string) string {
	return s[start : start+length]
}
