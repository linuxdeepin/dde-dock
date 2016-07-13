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
// char* get_icon_file_path(const char* name);
// void init_deepin();
// char* get_data_uri_by_path(const char* path);
import "C"
import (
	"path/filepath"
	"strings"
	"unsafe"
)

func getIconFilePath(name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	dotIndex := strings.LastIndex(name, ".")
	if dotIndex != -1 {
		ext := name[dotIndex+1:]
		logger.Debugf("getIconFilePath ext: %q", ext)
		switch ext {
		case "jpg", "png", "svg", "xpm":
			// remove ext
			name = name[:dotIndex]
		}
	}
	return _getIconFilePath(name)
}

func _getIconFilePath(name string) string {
	logger.Debugf("_getIconFilePath name: %q", name)
	cName := C.CString(name)
	cPath := C.get_icon_file_path(cName)
	path := C.GoString(cPath)
	C.free(unsafe.Pointer(cPath))
	C.free(unsafe.Pointer(cName))
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
