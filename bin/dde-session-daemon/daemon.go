package main

import _ "dde-daemon/keybinding"
import _ "dde-daemon/screensaver"
import _ "dde-daemon/power"
import _ "dde-daemon/audio"
import _ "dde-daemon/themes"

import _ "dde-daemon/clipboard"
import _ "dde-daemon/datetime"
import _ "dde-daemon/mime"
import _ "dde-daemon/screen_edges"
import _ "dde-daemon/bluetooth"
import _ "dde-daemon/grub2"
import _ "dde-daemon/network"
import _ "dde-daemon/mounts"

import _ "dde-daemon/dock"
import _ "dde-daemon/launcher"

import "dlib/glib-2.0"

//#cgo pkg-config:gtk+-3.0
//#include <gtk/gtk.h>
//void init(){gtk_init(0,0);}
import "C"
import . "dlib/gettext"
import "dlib"
import "dlib/logger"
import "os"
import "dlib/dbus"
import "dde-daemon"

var Logger = logger.NewLogger("com.deepin.daemon")

func init() {
	if !dlib.UniqueOnSession("com.deepin.daemon") {
		Logger.Warning("There already has an dde-daemon running.")
		return
	}
	if len(os.Args) >= 2 {
		for _, disabledModuleName := range os.Args[1:] {
			loader.Enable(disabledModuleName, false)
		}
	}
}

func main() {
	Logger.BeginTracing()
	defer Logger.EndTracing()

	InitI18n()
	Textdomain("dde-daemon")

	C.init()

	loader.Run()

	go func() {
		if err := dbus.Wait(); err != nil {
			Logger.Errorf("Lost dbus: %v", err)
			os.Exit(-1)
		} else {
			os.Exit(0)
		}
	}()

	ddeSessionRegister()
	glib.StartLoop()
}
