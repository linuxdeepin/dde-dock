package main

import _ "pkg.linuxdeepin.com/dde-daemon/inputdevices"
import _ "pkg.linuxdeepin.com/dde-daemon/screensaver"
import _ "pkg.linuxdeepin.com/dde-daemon/power"
import "pkg.linuxdeepin.com/lib/proxy"

import _ "pkg.linuxdeepin.com/dde-daemon/audio"

import _ "pkg.linuxdeepin.com/dde-daemon/appearance"

import _ "pkg.linuxdeepin.com/dde-daemon/clipboard"
import _ "pkg.linuxdeepin.com/dde-daemon/datetime"
import _ "pkg.linuxdeepin.com/dde-daemon/mime"

import _ "pkg.linuxdeepin.com/dde-daemon/screenedge"

import _ "pkg.linuxdeepin.com/dde-daemon/bluetooth"

import _ "pkg.linuxdeepin.com/dde-daemon/network"
import _ "pkg.linuxdeepin.com/dde-daemon/mounts"

import _ "pkg.linuxdeepin.com/dde-daemon/dock"
import _ "pkg.linuxdeepin.com/dde-daemon/launcher"
import _ "pkg.linuxdeepin.com/dde-daemon/keybinding"

import _ "pkg.linuxdeepin.com/dde-daemon/dsc"
import _ "pkg.linuxdeepin.com/dde-daemon/mpris"
import _ "pkg.linuxdeepin.com/dde-daemon/systeminfo"

import _ "pkg.linuxdeepin.com/dde-daemon/sessionwatcher"

import "pkg.linuxdeepin.com/lib/glib-2.0"

//#cgo pkg-config:gtk+-3.0
//#include <gtk/gtk.h>
//void init(){gtk_init(0,0);}
import "C"
import . "pkg.linuxdeepin.com/lib/gettext"
import "pkg.linuxdeepin.com/lib"
import "pkg.linuxdeepin.com/lib/log"
import "os"
import "pkg.linuxdeepin.com/lib/dbus"
import "pkg.linuxdeepin.com/dde-daemon"

var logger = log.NewLogger("com.deepin.daemon")

func main() {
	if !lib.UniqueOnSession("com.deepin.daemon") {
		logger.Warning("There already has an dde-daemon running.")
		return
	}
	if len(os.Args) >= 2 {
		for _, disabledModuleName := range os.Args[1:] {
			loader.Enable(disabledModuleName, false)
		}
	}
	logger.BeginTracing()
	defer logger.EndTracing()

	InitI18n()
	Textdomain("dde-daemon")

	C.init()
	proxy.SetupProxy()

	initPlugins()
	listenDaemonSettings()

	go func() {
		if err := dbus.Wait(); err != nil {
			logger.Errorf("Lost dbus: %v", err)
			os.Exit(-1)
		} else {
			logger.Info("dbus connection is closed by user")
			os.Exit(0)
		}
	}()

	loader.Start()
	defer loader.Stop()

	ddeSessionRegister()
	dbus.DealWithUnhandledMessage()
	glib.StartLoop()
}
