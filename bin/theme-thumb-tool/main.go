/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	TypeAll        = "all"
	TypeGtk        = "gtk"
	TypeIcon       = "icon"
	TypeCursor     = "cursor"
	TypeBackground = "background"
)

var (
	forceFlag = kingpin.Flag("force", "Force generate thumbnails").Short('f').Bool()
	destDir   = kingpin.Flag("output", "Thumbnails output directory").Default("").Short('o').String()
	thumbType = kingpin.Arg("type", "Thumbnail type, such as: gtk, icon, cursor...").Default("all").String()
)

func main() {
	kingpin.Parse()
	argNum := len(os.Args)
	if argNum == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage()
	}

	var thumbFiles []string
	switch *thumbType {
	case TypeAll:
		thumbFiles = genAllThumbnails(*forceFlag)
	case TypeGtk:
		thumbFiles = genGtkThumbnails(*forceFlag)
	case TypeIcon:
		thumbFiles = genIconThumbnails(*forceFlag)
	case TypeCursor:
		thumbFiles = genCursorThumbnails(*forceFlag)
	case TypeBackground:
		thumbFiles = genBgThumbnails(*forceFlag)
	default:
		usage()
	}
	moveThumbFiles(thumbFiles)
}

func usage() {
	fmt.Println("Desc:")
	fmt.Println("\ttheme-thumb-tool - gtk/icon/cursor/background thumbnail batch generator")
	fmt.Println("Usage:")
	fmt.Println("\ttheme-thumb-tool [Option] [Type]")
	fmt.Println("Option:")
	fmt.Println("\t-f --force: force to generate thumbnail regardless of file exist")
	fmt.Println("\t-o --output: thumbnails output directory")
	fmt.Println("Type:")
	fmt.Println("\tall: generate all of the following types thumbnails")
	fmt.Println("\tgtk: generate all gtk theme thumbnails")
	fmt.Println("\ticon: generate all icon theme thumbnails")
	fmt.Println("\tcursor: generate all cursor theme thumbnails")
	fmt.Println("\tbackground: generate all background thumbnails")

	os.Exit(0)
}

func moveThumbFiles(files []string) {
	if len(*destDir) == 0 {
		return
	}

	err := os.MkdirAll(*destDir, 0755)
	if err != nil {
		fmt.Printf("Create '%s' failed: %v\n", *destDir, err)
		return
	}
	for _, file := range files {
		dest := path.Join(*destDir, path.Base(file))
		if !*forceFlag && dutils.IsFileExist(dest) {
			continue
		}
		err := dutils.CopyFile(file, dest)
		os.Remove(file)
		if err != nil {
			fmt.Printf("Move '%s' to '%s' failed: %v\n", file, dest, err)
			continue
		}
	}
}
