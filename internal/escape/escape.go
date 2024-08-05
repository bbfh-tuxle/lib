package escape

import "strings"

func EscapeString(str string) string {
	str = strings.ReplaceAll(str, "\n", "\t")
	str = strings.ReplaceAll(str, "\r", "")
	return str
}

func UnescapeString(str string) string {
	str = strings.ReplaceAll(str, "\t", "\n")
	return str
}
