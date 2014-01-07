package main

import (
	"fmt"
)

const (
	LOG_INFO = iota
	LOG_ERROR
)

// TODO just for test
func log(msg string) {
	fmt.Printf("==> %s\n", msg)
}

func logError(msg string) {
	log(fmt.Sprintf("[ERROR] %s", msg))
}
