package core

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// LeftTrim trim characters from the left-side of the input.
// If second argument is empty, it's will be remove leading spaces.
func LeftTrim(str, chars string) string {
	pattern := "^\\s+"
	if chars != "" {
		pattern = "^[" + chars + "]+"
	}
	r, _ := regexp.Compile(pattern)
	return string(r.ReplaceAll([]byte(str), []byte("")))
}

// RightTrim trim characters from the right-side of the input.
// If second argument is empty, it's will be remove spaces.
func RightTrim(str, chars string) string {
	pattern := "\\s+$"
	if chars != "" {
		pattern = "[" + chars + "]+$"
	}
	r, _ := regexp.Compile(pattern)
	return string(r.ReplaceAll([]byte(str), []byte("")))
}

// Trim trim characters from both sides of the input.
// If second argument is empty, it's will be remove spaces.
func Trim(str, chars string) string {
	return LeftTrim(RightTrim(str, chars), chars)
}

// ToCamelCase converts from underscore separated form to camel case form.
// Ex.: my_func => MyFunc
func ToCamelCase(value interface{}) string {
	s := ToString(value)
	return strings.Replace(strings.Title(strings.Replace(strings.ToLower(s), "_", " ", -1)), " ", "", -1)
}

// ToLower convert the value string into lowercase format.
func ToLower(value interface{}) string {
	return strings.ToLower(ToString(value))
}

// ToLowerCamelCase converts from underscore separated form to lower camel case form.
// Ex.: my_func => myFunc
func ToLowerCamelCase(value interface{}) string {
	a := []rune(ToCamelCase(value))
	if len(a) > 0 {
		a[0] = unicode.ToLower(a[0])
	}
	return string(a)
}

// ToString convert the input to a string.
func ToString(value interface{}) string {
	res := fmt.Sprintf("%v", value)
	return string(res)
}
