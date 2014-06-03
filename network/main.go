package network

import (
	//"dlib"
	"dlib/dbus"
	liblogger "dlib/logger"
)

var (
	logger     = liblogger.NewLogger(dbusNetworkDest)
	manager    *Manager
	running    bool
	notifyStop = make(chan int, 100)
)

func Start() {
	/*
		logger.BeginTracing()
	*/

	if running {
		logger.Info(dbusNetworkDest, "already running")
		return
	}
	running = true
	defer func() {
		running = false
	}()

	/*
		if !dlib.UniqueOnSession(dbusNetworkDest) {
			logger.Warning("dbus unique:", dbusNetworkDest)
			return
		}
	*/

	initSlices() // initialize slice code

	manager = NewManager()
	err := dbus.InstallOnSession(manager)
	if err != nil {
		logger.Error("register dbus interface failed: ", err)
		return
	}

	// initialize manager after dbus installed
	manager.initManager()
	dbus.DealWithUnhandledMessage()

	notifyStop = make(chan int, 100) // reset signal to avoid repeat stop action
	notfiyDbusStop := make(chan int)
	/*
		go func() {
			err := dbus.Wait()
			if err != nil {
				logger.Error("lost dbus session:", err)
			} else {
				logger.Info("dbus session stoped")
			}
			notfiyDbusStop <- 1
		}()
	*/

	select {
	case <-notifyStop:
	case <-notfiyDbusStop:
	}
	DestroyManager(manager)
}

func Stop() {
	if !running {
		logger.Info(dbusNetworkDest, "already stopped")
		return
	}
	notifyStop <- 1
}
