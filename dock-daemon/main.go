package main

import (
	"dlib"
	"dlib/dbus"
	liblogger "dlib/logger"
	"os"
)

var logger = liblogger.NewLogger("dde-daemon/dock-daemon")

func main() {
	if !dlib.UniqueOnSession("com.deepin.daemon.daemon.dock-daemon") {
		logger.Warning("Anohter com.deepin.daemon.dock-daemon is running")
		return
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Fatalf("%v", err)
		}
	}()

	dlib.InitI18n()

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

	d := NewDockedAppManager()
	err = dbus.InstallOnSession(d)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		os.Exit(1)
	}

	s := NewSetting()
	err = dbus.InstallOnSession(s)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		os.Exit(1)
	}

	go dlib.StartLoop()

	// swm := NewSpecialWindowManager()
	// err = dbus.InstallOnSession(swm)
	// if err != nil {
	// 	logger.Error("register dbus interface failed:", err)
	// }
	// go swm.listenRootWindow()

	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); err != nil {
		logger.Errorf("lost dbus session: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
