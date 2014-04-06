package main

import (
	"dlib"
	"dlib/dbus"
	"dlib/gio-2.0"
	"fmt"
	"log"
	"os"
)

func main() {
	dlib.InitI18n()
	// DesktopAppInfo.ShouldShow does not know deepin.
	gio.DesktopAppInfoSetDesktopEnv("Deepin")

	initCategory()
	fmt.Println("init category done")

	initItems()
	fmt.Println("init items done")

	initDBus()
	fmt.Println("init dbus done")

	if tree != nil {
		defer tree.DestroyTrie(treeId)
	}
	dbus.DealWithUnhandledMessage()
	go dlib.StartLoop()
	if err := dbus.Wait(); err != nil {
		log.Panicln("lost dbus session:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
