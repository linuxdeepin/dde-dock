package main

import (
	"dlib"
	"dlib/dbus"
	"fmt"
	"log"
	"os"
)

func main() {
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
