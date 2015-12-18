package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"pkg.deepin.io/lib/log"
	"strings"
)

// using go build -ldflags "-X main.__VERSION__ version" to set version.
var __VERSION__ = "unknown"

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

type CMD struct {
	app      *kingpin.Application
	verbose  *bool
	logLevel *string
	memprof  *string
	cpuprof  *string
}

func (cmd *CMD) Parse(args []string) string {
	subcmd := kingpin.MustParse(cmd.app.Parse(args))
	return subcmd
}

func (cmd *CMD) LogLevel() log.Priority {
	if *cmd.verbose {
		return log.LevelDebug
	}
	lv, _ := toLogLevel(*cmd.logLevel)
	return lv
}

func (cmd *CMD) MemProf() string {
	return *cmd.memprof
}

func (cmd *CMD) CpuProf() string {
	return *cmd.cpuprof
}

func InitCMD() *CMD {
	cmd := kingpin.New("dde-preload", "session preload daemon")
	cmd.Version("version " + __VERSION__)

	flags := new(CMD)
	flags.app = cmd
	flags.verbose = cmd.Flag("verbose", "Show much more message, the shorthand for --loglevel debug, if specificed, loglevel is ignored.").Short('v').Bool()
	flags.logLevel = cmd.Flag("loglevel", "Set log level, possible value is error/warn/info/debug/no.").Short('l').String()
	flags.memprof = cmd.Flag("memprof", "Write memory profile to specific file").String()
	flags.cpuprof = cmd.Flag("cpuprof", "Write cpu profile to specific file").String()

	return flags
}
