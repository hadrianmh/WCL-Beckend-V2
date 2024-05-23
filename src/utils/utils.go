package utils

import (
	"strings"
	"unicode"
)

func Ucfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}

func StrReplaceAll(str string, rm string, strAppend string) string {
	return strings.ReplaceAll(str, rm, strAppend)
}
