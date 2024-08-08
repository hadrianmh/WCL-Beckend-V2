package utils

import (
	"strings"
	"time"
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

func PriceFilter(price string) string {
	filter1 := strings.ReplaceAll(price, ".", "")
	filter2 := strings.ReplaceAll(filter1, ",", ".")
	return filter2
}

func InArray(str []string, value string) bool {
	for _, item := range str {
		if item == value {
			return true
		}
	}
	return false
}

func ArrayUnique(str []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range str {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// getImageFormat extracts the image format from a data URL
// ex: data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...
func GetImageFormat(dataURL string) string {
	if strings.HasPrefix(dataURL, "data:") {
		semicolonIndex := strings.Index(dataURL, ";")
		if semicolonIndex > 5 {
			return dataURL[11:semicolonIndex]
		}
	}
	return ""
}

func ValidateReportFormatDate(dateStr string, format string) bool {
	// Try to parse the date string using the allowed format
	if _, err := time.Parse(format, dateStr); err != nil {
		return false
	}
	return true
}
