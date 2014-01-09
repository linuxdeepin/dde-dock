package main

import (
	"fmt"
)

const (
	LOG_INFO = iota
	LOG_ERROR
)

// TODO just for test
func log(format string, v ...interface{}) {
	fmt.Printf("==> "+format+"\n", v...)
}

func logInfo(format string, v ...interface{}) {
	log(fmt.Sprintf("[INFO] "+format, v...))
}

func logError(format string, v ...interface{}) {
	log(fmt.Sprintf("[ERROR] "+format, v...))
}
