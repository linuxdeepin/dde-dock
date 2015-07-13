package main

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// void gtkInit(){ gtk_init(NULL, NULL); }
import "C"
import (
	"fmt"
	"pkg.deepin.io/dde/daemon/desktop"
	"pkg.deepin.io/lib/glib-2.0"
)

func main() {
	C.gtkInit()
	err := desktop.NewDaemon().Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	glib.StartLoop()
}
