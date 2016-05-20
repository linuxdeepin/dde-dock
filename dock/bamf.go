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

#define KDE4_PREFIX "kde4"

char* getAppIdFromXid(guint32 xid) {
	static BamfMatcher* matcher = NULL;
	if (matcher == NULL) {
		matcher = bamf_matcher_get_default();
	}
	BamfApplication* app = bamf_matcher_get_application_for_xid(matcher, xid);
	if (app == NULL) {
		return NULL;
	}
	const char* desktop_file = bamf_application_get_desktop_file(app);
	if (desktop_file == NULL) {
		return NULL;
	}

	char* dirname = g_path_get_dirname(desktop_file);
	char* basedir = g_path_get_basename(dirname);
	char* basename = g_path_get_basename(desktop_file);

	free(dirname);

	if (g_strcmp0(basedir, KDE4_PREFIX) != 0) {
		free(basedir);
		return basename;
	} else {
		char* result = g_strconcat(KDE4_PREFIX, "-", basename, NULL);

		free(basedir);
		free(basename);

		return result;
	}
}
*/
import "C"
import "unsafe"
import (
	"github.com/BurntSushi/xgb/xproto"
	"strings"
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

func getAppIDFromXid(xid xproto.Window) string {
	cAppId := C.getAppIdFromXid(C.guint32(uint32(xid)))
	if cAppId == nil {
		return ""
	}
	appId := C.GoString(cAppId)
	defer C.free(unsafe.Pointer(cAppId))
	return strings.TrimSuffix(appId, ".desktop")
}
