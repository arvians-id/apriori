package utils

import "strconv"

func StrToInt(str string) int {
	integer, _ := strconv.Atoi(str)
	return integer
}
