package conv

import "strings"

func JoinStringSlice(slice []string, sep string) string {
	if len(slice) == 0 {
		return ""
	}
	if len(slice) == 1 {
		return slice[0]
	}
	return strings.Join(slice, sep)
}

func ReplaceNewlineWithSpace(input string) string {
	return strings.ReplaceAll(input, "\n", " ")
}
