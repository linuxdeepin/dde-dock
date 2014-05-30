package launcher

import (
	// "code.google.com/p/gettext-go/gettext"
	"dlib"
	"dlib/dbus"
	. "dlib/gettext"
	"dlib/gio-2.0"
	"dlib/glib-2.0"
	l "dlib/logger"
	"log"
	"os"
)

var logger *l.Logger = l.NewLogger("dde-daemon/launcher-daemon")

func Start() {
	defer logger.EndTracing()

	if !dlib.UniqueOnSession("com.deepin.dde.daemon.Launcher") {
		logger.Warning("Another com.deepin.daemon.Launcher is running.")
		return
	}
	InitI18n()
	// DesktopAppInfo.ShouldShow does not know deepin.
	gio.DesktopAppInfoSetDesktopEnv("Deepin")

	initCategory()
	logger.Info("init category done")

	initItems()
	logger.Info("init items done")

	initDBus()
	logger.Info("init dbus done")

	if tree != nil {
		defer tree.DestroyTrie(treeId)
	}
	dbus.DealWithUnhandledMessage()
	go glib.StartLoop()
	if err := dbus.Wait(); err != nil {
		log.Panicln("lost dbus session:", err)
		os.Exit(1)
	}
}
