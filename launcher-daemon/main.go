package main

import (
	// "code.google.com/p/gettext-go/gettext"
	"dlib"
	"dlib/dbus"
	"dlib/gio-2.0"
	l "dlib/logger"
	"log"
	"os"
)

var logger *l.Logger

func main() {
	dlib.InitI18n()
	// DesktopAppInfo.ShouldShow does not know deepin.
	gio.DesktopAppInfoSetDesktopEnv("Deepin")

	initCategory()
	logger = l.NewLogger("dde-daemon/launcher-daemon")
	logger.Info("init category done")

	initItems()
	logger.Info("init items done")

	initDBus()
	logger.Info("init dbus done")

	if tree != nil {
		defer tree.DestroyTrie(treeId)
	}
	dbus.DealWithUnhandledMessage()
	go dlib.StartLoop()
	if err := dbus.Wait(); err != nil {
		log.Panicln("lost dbus session:", err)
		os.Exit(1)
	}
}
