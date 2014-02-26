package main

import (
	"dlib"
	"dlib/dbus"
	"fmt"
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
	dlib.StartLoop()
}
