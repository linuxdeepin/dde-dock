package main

import (
	"dlib/dbus"
	"dlib/logger"
	"flag"
	"os"
)

var (
	Logger   = logger.NewLogger("com.deepin.daemon.Network")
	manager  *Manager
	argDebug bool
)

func main() {
	defer Logger.EndTracing()

	// configure logger
	flag.BoolVar(&argDebug, "d", false, "debug mode")
	flag.BoolVar(&argDebug, "debug", false, "debug mode")
	flag.Parse()
	if argDebug {
		Logger.SetLogLevel(logger.LEVEL_DEBUG)
	}

	manager = NewManager()
	err := dbus.InstallOnSession(manager)
	if err != nil {
		Logger.Error("register dbus interface failed: ", err)
		os.Exit(1)
	}

	// initialize manager after configuring dbus
	manager.initManager()

	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		Logger.Error("lost dbus session: ", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
