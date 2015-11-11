/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

#include <X11/Xlib.h>
#include <X11/extensions/record.h>
#include <glib.h>

#include "record.h"
#include "_cgo_export.h"

typedef struct _XRecordGrabInfo {
	Display *ctrl_disp;
	Display *data_disp;
	XRecordRange *range;
	XRecordContext context;
} XRecordGrabInfo;

static void grab_key_event_cb (XPointer user_data, XRecordInterceptData *hook);
static gpointer enable_ctx_thread (gpointer user_data);

static XRecordGrabInfo *grab_info = NULL;
static int pressed_cnt = 0;
static gboolean button_pressed = FALSE;

void
xrecord_grab_init ()
{
	if (grab_info != NULL) {
		g_debug("XRecord grab has been init...\n");
		return;
	}

	grab_info = g_new0 (XRecordGrabInfo, 1);

	if ( !grab_info ) {
		g_warning ("Alloc XRecordGrabInfo memory failed...");
		return;
	}

	grab_info->ctrl_disp = XOpenDisplay (NULL);
	grab_info->data_disp = XOpenDisplay (NULL);

	if ( !grab_info->ctrl_disp || !grab_info->data_disp ) {
		g_warning ("Unable to connect to X server...");
		xrecord_grab_finalize ();
		return;
	}

	gint dummy;

	if ( !XQueryExtension (grab_info->ctrl_disp, "XTEST",
	                       &dummy, &dummy, &dummy) ) {
		g_warning ("XTest extension missing...");
		xrecord_grab_finalize ();
		return;
	}

	if ( !XRecordQueryVersion (grab_info->ctrl_disp, &dummy, &dummy) ) {
		g_warning ("Failed to obtain xrecord version...");
		xrecord_grab_finalize ();
		return;
	}

	grab_info->range = XRecordAllocRange ();

	if ( !grab_info->range ) {
		g_warning ("Alloc XRecordRange memory failed...");
		xrecord_grab_finalize ();
		return;
	}

	grab_info->range->device_events.first = KeyPress;
	grab_info->range->device_events.last = ButtonRelease;

	XRecordClientSpec spec = XRecordAllClients;
	grab_info->context = XRecordCreateContext (
	                         grab_info->data_disp, 0, &spec, 1, &grab_info->range, 1);

	if ( !grab_info->context ) {
		g_warning ("Unable to create context...");
		xrecord_grab_finalize();
		return;
	}

	XSynchronize (grab_info->ctrl_disp, TRUE);
	XFlush (grab_info->ctrl_disp);

	GThread *thrd = g_thread_new ("enable context",
	                              (GThreadFunc)enable_ctx_thread, NULL);

	if ( !thrd ) {
		g_warning ("Unable to create thread...");
		xrecord_grab_finalize ();
		return;
	}
	g_thread_unref(thrd);
}

void
xrecord_grab_finalize ()
{
	if (!grab_info) {
		return;
	}

	if (grab_info->context) {
		XRecordDisableContext(grab_info->data_disp, grab_info->context);
		XRecordFreeContext(grab_info->data_disp, grab_info->context);
	}

	if (grab_info->range) {
		XFree(grab_info->range);
		grab_info->range = NULL;
	}

	if (grab_info->ctrl_disp) {
		XCloseDisplay (grab_info->ctrl_disp);
		grab_info->ctrl_disp = NULL;
	}

	if (grab_info->data_disp) {
		XCloseDisplay (grab_info->data_disp);
		grab_info->data_disp = NULL;
	}

	if (grab_info) {
		g_free (grab_info);
		grab_info = NULL;
	}
}

/*
 * check keyboard && mouse whether has grabbed
 * if grabed, return 1, otherwise return 0
 */
int
is_grabbed ()
{
	Display *dpy;
	Window root;
	int ret;

	dpy = XOpenDisplay (0);
	root = DefaultRootWindow (dpy);

	ret = XGrabKeyboard (dpy, root, False,
	                     GrabModeSync, GrabModeSync, CurrentTime);

	if ( ret == AlreadyGrabbed ) {
		g_debug ("AlreadyGrabbed!\n");
		/*XUngrabKeyboard (dpy, CurrentTime);*/
		XCloseDisplay(dpy);
		return 1;
	}


	XUngrabKeyboard (dpy, CurrentTime);
	XCloseDisplay(dpy);
	return 0;
}

static gpointer
enable_ctx_thread (gpointer user_data)
{
	if ( !XRecordEnableContext (grab_info->data_disp, grab_info->context,
	                            grab_key_event_cb, NULL) ) {
		g_warning ("Unable to enable context...");
		xrecord_grab_finalize ();
	}

	g_thread_exit (NULL);
}

static void
grab_key_event_cb (XPointer user_data, XRecordInterceptData *hook)
{
	if ( hook->category != XRecordFromServer ) {
		XRecordFreeData(hook);
		g_warning ("Data not from X server...");
		return;
	}

	int event_type = hook->data[0];
	KeyCode keycode = hook->data[1];

	/*g_debug ("event type: %d, code: %d\n", (int)event_type, (int)keycode);*/

	switch (event_type) {
	case KeyPress:
		pressed_cnt++;
		break;

	case KeyRelease:
		/*g_debug ("pressed_cnt: %d\n", pressed_cnt);*/

		if ((pressed_cnt == 1) && !button_pressed)  {
			//handler event
			handleSingleKeyEvent(keycode, 0);
		}

		pressed_cnt = 0;
		break;

	case ButtonPress:
		button_pressed = TRUE;
		break;

	case ButtonRelease:
		button_pressed = FALSE;
		break;

	default:
		pressed_cnt = 0;
		break;
	}

	XRecordFreeData(hook);
}
