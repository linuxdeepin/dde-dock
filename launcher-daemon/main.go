package main

import (
	"dlib"
	"dlib/dbus"
	"fmt"
	"log"
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
		log.Fatal("lost dbus session:", err)
	}
}
