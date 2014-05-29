package network

import (
	"dlib"
	"dlib/dbus"
	liblogger "dlib/logger"
)

var (
	logger  = liblogger.NewLogger(dbusNetworkDest)
	manager *Manager
)

func start() {
	defer logger.EndTracing()

	if !dlib.UniqueOnSession(dbusNetworkDest) {
		logger.Warning("There already has an daemon running:", dbusNetworkDest)
		return
	}

	manager = NewManager()
	err := dbus.InstallOnSession(manager)
	if err != nil {
		logger.Error("register dbus interface failed: ", err)
		return
	}

	// initialize manager after dbus installed
	manager.initManager()

	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		logger.Error("lost dbus session:", err)
		return
	}
}

func stop() {
}
