package main

import (
	"dlib"
)

func init() {
	initCategory()
	initItems()
	initDBus()
}

func main() {
	if tree != nil {
		defer tree.DestroyTrie(treeId)
	}
	dlib.StartLoop()
}
