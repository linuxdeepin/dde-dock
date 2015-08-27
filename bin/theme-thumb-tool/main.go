package main

import (
	"fmt"
	"os"
)

const (
	TypeAll        = "all"
	TypeGtk        = "gtk"
	TypeIcon       = "icon"
	TypeCursor     = "cursor"
	TypeBackground = "background"
)

var (
	forceFlag bool = false
)

func main() {
	argNum := len(os.Args)
	if argNum == 1 {
		usage()
	}

	if argNum == 3 && os.Args[1] == "--force" {
		forceFlag = true
	}

	switch os.Args[argNum-1] {
	case TypeAll:
		genAllThumbnails(forceFlag)
	case TypeGtk:
		genGtkThumbnails(forceFlag)
	case TypeIcon:
		genIconThumbnails(forceFlag)
	case TypeCursor:
		genCursorThumbnails(forceFlag)
	case TypeBackground:
		genBgThumbnails(forceFlag)
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Desc:")
	fmt.Println("\ttheme-thumb-tool - gtk/icon/cursor/background thumbnail batch generator")
	fmt.Println("Usage: theme-thumb-tool [Option] [Type]")
	fmt.Println("Option:")
	fmt.Println("\t--force: force to generate thumbnail regardless of file exist")
	fmt.Println("Type:")
	fmt.Println("\tall: generate all of the following types thumbnails")
	fmt.Println("\tgtk: generate all gtk theme thumbnails")
	fmt.Println("\ticon: generate all icon theme thumbnails")
	fmt.Println("\tcursor: generate all cursor theme thumbnails")
	fmt.Println("\tbackground: generate all background thumbnails")

	os.Exit(0)
}
