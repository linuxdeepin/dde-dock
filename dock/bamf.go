/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

//#cgo pkg-config: libbamf3
/*
#include <stdlib.h>
#include <libbamf/bamf-matcher.h>
#include <libbamf/bamf-application.h>

char * getDesktopFromWindowByBamf(guint32 win) {
	static BamfMatcher* matcher = NULL;
	if (matcher == NULL ) {
		matcher = bamf_matcher_get_default();
	}
	BamfApplication* app = bamf_matcher_get_application_for_xid(matcher, win);
	if (app == NULL) {
		return NULL;
	}
	return g_strdup(bamf_application_get_desktop_file(app));
}

*/
import "C"
import "unsafe"
import (
	"github.com/BurntSushi/xgb/xproto"
)

func getDesktopFromWindowByBamf(win xproto.Window) string {
	cDesktop := C.getDesktopFromWindowByBamf(C.guint32(uint32(win)))
	if cDesktop == nil {
		return ""
	}
	desktop := C.GoString(cDesktop)
	defer C.free(unsafe.Pointer(cDesktop))
	return desktop
}
