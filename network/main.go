package network

import (
	"dlib"
	"dlib/dbus"
	liblogger "dlib/logger"
	"flag"
	"os"
)

var (
	logger   = liblogger.NewLogger(dbusNetworkDest)
	manager  *Manager
	argDebug bool
)

func main() {
	defer logger.EndTracing()

	if !dlib.UniqueOnSession(dbusNetworkDest) {
		logger.Warning("There already has an daemon running:", dbusNetworkDest)
		return
	}

	// configure logger
	flag.BoolVar(&argDebug, "d", false, "debug mode")
	flag.BoolVar(&argDebug, "debug", false, "debug mode")
	flag.Parse()
	if argDebug {
		logger.SetLogLevel(liblogger.LEVEL_DEBUG)
	}

	manager = NewManager()
	err := dbus.InstallOnSession(manager)
	if err != nil {
		logger.Error("register dbus interface failed: ", err)
		os.Exit(1)
	}

	// initialize manager after configuring dbus
	manager.initManager()

	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		logger.Error("lost dbus session:", err)
		os.Exit(1)
	}
}
