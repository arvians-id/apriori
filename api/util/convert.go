package util

import (
	"regexp"
	"strconv"
	"strings"
)

func StrToInt(str string) int {
	integer, _ := strconv.Atoi(str)
	return integer
}

func IntToStr(number int) string {
	str := strconv.Itoa(number)
	return str
}

func StrToBool(str string) bool {
	boolean, _ := strconv.ParseBool(str)
	return boolean
}

func UpperWords(str string) string {
	str = strings.TrimSpace(str)
	replace := func(word string) string {
		switch word {
		case "with", "in", "a":
			return word
		}
		return strings.Title(word)
	}

	r := regexp.MustCompile(`\w+`)
	str = r.ReplaceAllStringFunc(strings.ToLower(str), replace)

	return str
}
