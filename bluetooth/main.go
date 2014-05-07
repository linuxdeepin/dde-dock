package main

import (
	idbus "dbus/org/freedesktop/dbus/system"
	"dlib"
	"dlib/dbus"
	liblogger "dlib/logger"
	"flag"
	"os"
)

const (
	dbusBluezDest     = "org.bluez"
	dbusBluezPath     = "/org/bluez"
	dbusBluetoothDest = "com.deepin.daemon.Bluetooth"
	dbusBluetoothPath = "/com/deepin/daemon/Bluetooth"
	dbusBluetoothIfs  = "com.deepin.daemon.Bluetooth"
)

var (
	logger   = liblogger.NewLogger(dbusBluetoothDest)
	argDebug bool
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

	bluetooth := NewBluetooth()
	err := dbus.InstallOnSession(bluetooth)
	if err != nil {
		// don't panic or fatal here
		logger.Error("register dbus interface failed: ", err)
		os.Exit(1)
	}

	bluetooth.initBluetooth()

	// TODO test
	bluezObjectManager, err := idbus.NewObjectManager(dbusBluezDest, "/")
	if err != nil {
		panic(err)
	}
	bluezObjects, err := bluezObjectManager.GetManagedObjects()
	if err != nil {
		panic(err)
	}
	logger.Debug(bluezObjects)

	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		// don't panic or fatal here
		logger.Error("lost dbus session:", err)
		os.Exit(1)
	}
}
