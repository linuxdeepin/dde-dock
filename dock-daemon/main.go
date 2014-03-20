package main

import (
	"dlib/dbus"
	liblogger "dlib/logger"
	"os"
)

var logger = liblogger.NewLogger("dde-daemon/dock-daemon")

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Fatalf("%v", err)
		}
	}()

	// configure logger
	logger.SetRestartCommand("/usr/lib/deepin-daemon/grub2", "--debug")
	if stringInSlice("-d", os.Args) || stringInSlice("--debug", os.Args) {
		logger.SetLogLevel(liblogger.LEVEL_DEBUG)
	}

	m := NewManager()
	err := dbus.InstallOnSession(m)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		os.Exit(1)
	}

	m.watchEntries()
	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); err != nil {
		logger.Errorf("lost dbus session: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
