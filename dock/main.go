package dock

import (
	"dlib"
	"dlib/dbus"
	"dlib/gettext"
	"dlib/glib-2.0"
	liblogger "dlib/logger"
	"os"
	"os/exec"
)

var (
	logger = liblogger.NewLogger("dde-daemon/dock")
)

func Start() {
	defer logger.EndTracing()

	if !dlib.UniqueOnSession("com.deepin.daemon.Dock") {
		logger.Warning("Anohter com.deepin.daemon.Dock is running")
		return
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Fatalf("%v", err)
		}
	}()

	gettext.Bindtextdomain("dde-daemon", "/usr/share/locale")
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

	go glib.StartLoop()

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

	if err := dbus.Wait(); err != nil {
		logger.Errorf("lost dbus session: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
