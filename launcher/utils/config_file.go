/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package utils

// #cgo pkg-config: glib-2.0
// #include <glib.h>
import "C"

// GReloadUserSpecialDirsCache reloads user special dirs cache.
func GReloadUserSpecialDirsCache() {
	C.g_reload_user_special_dirs_cache()
}

func uniqueStringList(l []string) []string {
	m := make(map[string]bool, 0)
	for _, v := range l {
		m[v] = true
	}
	var n []string
	for k := range m {
		n = append(n, k)
	}
	return n
}
