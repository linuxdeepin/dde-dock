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

#include <stdio.h>
#include <string.h>
#include <X11/Xlib.h>
#include <X11/Xcursor/Xcursor.h>
#include <X11/extensions/Xfixes.h>

static char* findAlternative(const char *name);
static XcursorImages* xcLoadImages(const char *theme,
		const char *image, int size);
static unsigned long loadCursorHandle(Display *disp,
		const char *theme, const char *name, int size);

static char*
findAlternative(const char *name)
{
	if (!name) {
		return NULL;
	}

	// Qt uses non-standard names for some core cursors.
	// If Xcursor fails to load the cursor, Qt creates it with the correct name
	// using the core protcol instead (which in turn calls Xcursor).
	// We emulate that process here.
	// Note that there's a core cursor called cross, but it's not the one Qt expects.
	// Precomputed MD5 hashes for the hardcoded bitmap cursors in Qt and KDE.
	// Note that the MD5 hash for left_ptr_watch is for the KDE version of that cursor.
	static char* xcursor_alter[] = {
		"cross", "crosshair",
		"up_arrow", "center_ptr",
		"wait", "watch",
		"ibeam", "xterm",
		"size_all", "fleur",
		"pointing_hand", "hand2",
		"size_ver", "00008160000006810000408080010102",
		"size_hor", "028006030e0e7ebffc7f7070c0600140",
		"size_bdiag", "c7088f0f3e6c8088236ef8e1e3e70000",
		"size_fdiag", "fcf1c3c7cd4491d801f1e1c78f100000",
		"whats_this", "d9ce0ab605698f320427677b458ad60b",
		"split_h", "14fef782d02440884392942c11205230",
		"split_v", "2870a09082c103050810ffdffffe0204",
		"forbidden", "03b6e0fcb3499374a867c041f52298f0",
		"left_ptr_watch", "3ecb610c1bf2410f44200f48c40d3599",
		"hand2", "e29285e634086352946a0e7090d73106",
		"openhand", "9141b49c8149039304290b508d208c40",
		"closedhand", "05e88622050804100c20044008402080",
		NULL};

	int i;
	for (i = 0; xcursor_alter[i] != NULL; i+=2) {
		if (strcmp(xcursor_alter[i], name) == 0) {
			return xcursor_alter[i+1];
		}
	}

	return NULL;
}

static XcursorImages*
xcLoadImages(const char *theme, const char *image, int size)
{
	if (!theme || !image) {
		return NULL;
	}

	return XcursorLibraryLoadImages(image, theme, size);
}

static unsigned long
loadCursorHandle(Display *disp, const char *theme, const char *name, int size)
{
	if (size == -1) {
		size = XcursorGetDefaultSize(disp);
	}

	// Load the cursor images
	XcursorImages *images = NULL;
	images = xcLoadImages(theme, name, size);
	if (!images) {
		images = xcLoadImages(theme,
		                      findAlternative(name), size);
		if (!images) {
			return 0;
		}
	}

	unsigned long handle = (unsigned long)XcursorImagesLoadCursor(disp,
	                       images);
	XcursorImagesDestroy(images);

	return handle;
}

int
apply_qt_cursor(const char *theme)
{
	if (!theme) {
		fprintf(stderr, "Cursor theme is NULL\n");
		return -1;
	}

	/**
	 *  Only running once
	 *  Setting it by startdde is not effective, need to set again, why?
	 **/
	static int running = 0;
	if (running == 2) {
		return 0;
	}

	/**
	 * Fixed Qt cursor not work when cursor theme changed.
	 * For details see: lxqt-config/lxqt-config-cursor
	 *
	 * XFixes multiple qt cursor name, a X Error will be occured.
	 * Now only XFixes qt cursor name 'left_ptr'
	 * Why?
	 **/
	static char* list[] = {
		// Qt cursors
		"left_ptr",
		// "up_arrow",
		// "cross",
		// "wait",
		// "left_ptr_watch",
		// "ibeam",
		// "size_ver",
		// "size_hor",
		// "size_bdiag",
		// "size_fdiag",
		// "size_all",
		// "split_v",
		// "split_h",
		// "pointing_hand",
		// "openhand",
		// "closedhand",
		// "forbidden",
		// "whats_this",
		// X core cursors
		// "X_cursor",
		"right_ptr",
		"hand1",
		"hand2",
		"watch",
		"xterm",
		"crosshair",
		"left_ptr_watch",
		// "center_ptr",  // invalid Cursor parameter, why?
		"sb_h_double_arrow",
		"sb_v_double_arrow",
		"fleur",
		"top_left_corner",
		"top_side",
		"top_right_corner",
		"right_side",
		"bottom_right_corner",
		"bottom_side",
		"bottom_left_corner",
		"left_side",
		"question_arrow",
		"pirate",
		NULL};

	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open display failed\n");
		return -1;
	}

	int i;
	for (i = 0; list[i] != NULL; i++) {
		Cursor cursor = (Cursor)loadCursorHandle(disp, theme,
		                list[i], -1);
		XFixesChangeCursorByName(disp, cursor, list[i]);
		// FIXME: do we need to free the cursor?
	}
	XCloseDisplay(disp);
	running++;

	return 0;
}
