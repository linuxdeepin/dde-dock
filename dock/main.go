package dock

import (
	"dlib/dbus"
	liblogger "dlib/logger"
	"os"
	"os/exec"
)

var (
	logger = liblogger.NewLogger("dde-daemon/dock")
)

func Start() {
	logger.BeginTracing()
	defer logger.EndTracing()

	initDeepin()

	// configure logger
	if stringInSlice("-d", os.Args) || stringInSlice("--debug", os.Args) {
		logger.SetLogLevel(liblogger.LEVEL_DEBUG)
	}

	m := NewEntryProxyerManager()
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

	cm := NewClientManager()
	err = dbus.InstallOnSession(cm)
	if err != nil {
		logger.Error("register dbus interface failed:", err)
	}
	go cm.listenRootWindow()

	region := NewRegion()
	dbus.InstallOnSession(region)

	dbus.DealWithUnhandledMessage()

	initialize()

	go exec.Command("/usr/bin/dde-dock").Run()
}
