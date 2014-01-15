package main

import (
	"fmt"
	"os"
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

func logWarn(format string, v ...interface{}) {
	log(fmt.Sprintf("[WARN] "+format, v...))
}

func logError(format string, v ...interface{}) {
	log(fmt.Sprintf("[WARN] "+format, v...))
}

func logPanic(format string, v ...interface{}) {
	s := fmt.Sprintf("[ERROR] "+format, v...)
	log(s)
	panic(s)
}

func logFatal(format string, v ...interface{}) {
	log(fmt.Sprintf("[ERROR] "+format, v...))
	os.Exit(1)
}
