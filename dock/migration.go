/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

// this file should be rewrite

//#cgo pkg-config: glib-2.0 gio-unix-2.0 gtk+-3.0
//#include <stdlib.h>
// char* icon_name_to_path(const char* name, int size);
// void init_deepin();
// char* get_data_uri_by_path(const char* path);
import "C"
import "unsafe"

func get_theme_icon(name string, size int) string {
	iconName := C.CString(name)
	defer func() {
		C.free(unsafe.Pointer(iconName))
	}()
	cPath := C.icon_name_to_path(iconName, C.int(size))
	defer C.free(unsafe.Pointer(cPath))
	path := C.GoString(cPath)
	return path
}

func initDeepin() {
	C.init_deepin()
}

func xpm_to_dataurl(icon string) string {
	iconName := C.CString(icon)
	defer func() {
		C.free(unsafe.Pointer(iconName))
	}()
	cDataUri := C.get_data_uri_by_path(iconName)
	defer C.free(unsafe.Pointer(cDataUri))
	return C.GoString(cDataUri)

}
