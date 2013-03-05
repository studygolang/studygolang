package util

import (
	"strconv"
)

// 必须是int类型，否则panic
func MustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
