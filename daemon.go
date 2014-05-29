package main

import "dde-daemon/clipboard"
import "dde-daemon/audio"
import "dde-daemon/power"
import "dde-daemon/display"
import "dde-daemon/keybinding"

//import "dde-daemon/dock"
//import "dde-daemon/launcher"

//import "dde-daemon/inputdevices"
import "dlib/glib-2.0"

//#cgo pkg-config:gtk+-3.0
//#include <gtk/gtk.h>
//void init(){gtk_init(0,0);}
import "C"

func main() {
	C.init()
	go clipboard.Start()
	go audio.Start()
	go power.Start()
	go display.Start()

	//go dock.Start()
	//go launcher.Start()

	go keybinding.Start()
	//go inputdevices.Start()

	glib.StartLoop()
}
