package dock

import (
	"dlib/dbus"
	liblogger "dlib/logger"
	"os"
	"os/exec"
)

var (
	logger = liblogger.NewLogger("dde-daemon/dock")

	region          *Region           = nil
	setting         *Setting          = nil
	hideModemanager *HideStateManager = nil
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

	setting = NewSetting()
	err = dbus.InstallOnSession(setting)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		os.Exit(1)
	}

	hideModemanager = NewHideStateManager(setting.GetHideMode())
	err = dbus.InstallOnSession(hideModemanager)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		os.Exit(1)
	}
	hideModemanager.UpdateState()

	cm := NewClientManager()
	err = dbus.InstallOnSession(cm)
	if err != nil {
		logger.Error("register dbus interface failed:", err)
	}
	go cm.listenRootWindow()

	region = NewRegion()
	dbus.InstallOnSession(region)

	dbus.DealWithUnhandledMessage()

	initialize()
}
