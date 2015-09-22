package main

import (
	"fmt"
	"pkg.deepin.io/lib/log"
	"strings"
)

func toLogLevel(name string) (log.Priority, error) {
	name = strings.ToLower(name)
	logLevel := log.LevelInfo
	var err error
	switch name {
	case "":
	case "error":
		logLevel = log.LevelError
	case "warn":
		logLevel = log.LevelWarning
	case "info":
		logLevel = log.LevelInfo
	case "debug":
		logLevel = log.LevelDebug
	case "no":
		logLevel = log.LevelDisable
	default:
		err = fmt.Errorf("%s is not support", name)
	}

	return logLevel, err
}

type Flags struct {
	Verbose              *bool
	LogLevel             *string
	IgnoreMissingModules *bool
	ForceStart           *bool
}
