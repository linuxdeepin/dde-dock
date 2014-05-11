package main

import (
	"dlib"
	"dlib/dbus"
	liblogger "dlib/logger"
	"flag"
	"os"
)

var (
	logger   = liblogger.NewLogger(dbusBluetoothDest)
	argDebug bool

	bluetooth *Bluetooth
)

func main() {
	defer logger.EndTracing()

	if !dlib.UniqueOnSession(dbusBluetoothDest) {
		logger.Warning("There already has an daemon running:", dbusBluetoothDest)
		return
	}

	// configure logger
	flag.BoolVar(&argDebug, "d", false, "debug mode")
	flag.BoolVar(&argDebug, "debug", false, "debug mode")
	flag.Parse()
	if argDebug {
		logger.SetLogLevel(liblogger.LEVEL_DEBUG)
	}

	bluetooth = NewBluetooth()
	err := dbus.InstallOnSession(bluetooth)
	if err != nil {
		// don't panic or fatal here
		logger.Error("register dbus interface failed: ", err)
		os.Exit(1)
	}

	// initialize bluetooth after dbus interface installed
	bluetooth.initBluetooth()

	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		// don't panic or fatal here
		logger.Error("lost dbus session:", err)
		os.Exit(1)
	}
}
