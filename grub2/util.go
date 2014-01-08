package main

import (
	"strconv"
	"strings"
)

func quoteString(str string) string {
	return strconv.Quote(str)
}

func unquoteString(str string) string {
	if (strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`)) ||
		(strings.HasPrefix(str, `'`) && strings.HasSuffix(str, `'`)) {
		return str[1 : len(str)-1]
	}
	return str
}
