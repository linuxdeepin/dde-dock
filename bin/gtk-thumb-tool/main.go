/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

// #cgo pkg-config: gtk+-2.0 libmetacity-private
// #cgo CFLAGS: -Wall -g
// #include <stdlib.h>
// #include "common.h"
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(0)
		}
	}()

	if C.try_init() == 0 {
		return
	}
	if len(os.Args) != 3 {
		fmt.Printf("ERROR\n")
		fmt.Printf("Usage: %s <Theme> <Dest>\n", os.Args[0])
		return
	}

	theme := os.Args[1]
	dest := os.Args[2]

	cTheme := C.CString(theme)
	defer C.free(unsafe.Pointer(cTheme))
	cDest := C.CString(dest)
	defer C.free(unsafe.Pointer(cDest))

	ret := C.gen_gtk_thumbnail(cTheme, cDest)
	if int(ret) == -1 {
		fmt.Printf("Generate Gtk Thumbnail Failed", cTheme, cDest)
	}
}
