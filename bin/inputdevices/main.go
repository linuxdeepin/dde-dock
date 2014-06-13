package main

import _ "dde-daemon/inputdevices"

import "dlib/glib-2.0"

//#cgo pkg-config:gtk+-3.0
//#include <gtk/gtk.h>
//void init(){gtk_init(0,0);}
import "C"
import . "dlib/gettext"
import "dlib"
import "dlib/logger"
import "dlib/dbus"
import "dde-daemon"
import "os"

var Logger = logger.NewLogger("com.deepin.daemon.InputDevices")

func main() {
	if !dlib.UniqueOnSession("com.deepin.daemon.InputDevices") {
		Logger.Warning("There already has an dde-daemon running.")
		return
	}
	Logger.BeginTracing()
	defer Logger.EndTracing()

	InitI18n()
	Textdomain("dde-daemon")

	C.init()

	loader.Start()
	defer loader.Stop()

	go func() {
		if err := dbus.Wait(); err != nil {
			Logger.Errorf("Lost dbus: %v", err)
			os.Exit(-1)
		} else {
			os.Exit(0)
		}
	}()

	ddeSessionRegister()
	dbus.DealWithUnhandledMessage()
	glib.StartLoop()
}
