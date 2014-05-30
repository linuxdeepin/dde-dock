package main

import "dde-daemon/clipboard"
import "dde-daemon/audio"
import "dde-daemon/power"
import "dde-daemon/display"
import "dde-daemon/keybinding"
import "dde-daemon/datetime"
import "dde-daemon/mime"
import "dde-daemon/mounts"
import "dde-daemon/screen_edges"
import "dde-daemon/themes"

//import "dde-daemon/dock"
//import "dde-daemon/launcher"

//import "dde-daemon/inputdevices"
import "dlib/glib-2.0"

//#cgo pkg-config:gtk+-3.0
//#include <gtk/gtk.h>
//void init(){gtk_init(0,0);}
import "C"
import "time"
import . "dlib/gettext"
import "dlib"
import "dlib/logger"

var Logger = logger.NewLogger("com.deepin.daemon")

func main() {
	if !dlib.UniqueOnSession("com.deepin.daemon") {
		Logger.Warning("There already has an dde-daemon running.")
		return
	}
	InitI18n()
	Textdomain("dde-daemon")

	C.init()
	go clipboard.Start()
	go audio.Start()
	go power.Start()
	go display.Start()
	<-time.After(time.Second * 3)

	//go dock.Start()
	//go launcher.Start()

	go keybinding.Start()
	//go inputdevices.Start()
	go datetime.Start()
	go mime.Start()
	go mounts.Start()
	go themes.Start()
	go screen_edges.Start()

	startMprisDaemon()

	<-time.After(time.Second)
	glib.StartLoop()
}
