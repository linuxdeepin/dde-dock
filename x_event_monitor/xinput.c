/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

#include <X11/X.h>
#include <X11/extensions/XInput.h>

#include "xinput.h"
#include "_cgo_export.h"

#include <stdio.h>
#include <stdlib.h>


static int xi2_opcode;

int xinput_version(Display *display)
{
    XExtensionVersion *version;
    static int vers = -1;

    if (vers != -1)
	return vers;

    version = XGetExtensionVersion(display, INAME);

    if (version && (version != (XExtensionVersion*) NoSuchExtension)) {
	vers = version->major_version;
	XFree(version);
    }

    /* Announce our supported version so the server treats us correctly. */
    if (vers >= XI_2_Major)
    {
	const char *forced_version;
	int maj = 2, min = 2;

	forced_version = getenv("XINPUT_XI2_VERSION");
	if (forced_version) {
	    if (sscanf(forced_version, "%d.%d", &maj, &min) != 2) {
		fprintf(stderr, "Invalid format of XINPUT_XI2_VERSION "
			"environment variable. Need major.minor\n");
		exit(1);
	    }
	    printf("Overriding XI2 version to: %d.%d\n", maj, min);
	}

	XIQueryVersion(display, &maj, &min);
    }

    return vers;
}
 void select_events(Display* display)
 {
    XIEventMask m;
    m.deviceid = XIAllMasterDevices;
    m.mask_len = XIMaskLen(XI_LASTEVENT);
    m.mask = calloc(m.mask_len, sizeof(char));
    XISetMask(m.mask, XI_RawKeyPress);
    XISetMask(m.mask, XI_RawKeyRelease);
    XISetMask(m.mask, XI_RawButtonPress);
    XISetMask(m.mask, XI_RawButtonRelease);
    XISetMask(m.mask, XI_RawMotion);
    XISetMask(m.mask, XI_RawTouchBegin);
    XISetMask(m.mask, XI_RawTouchUpdate);
    XISetMask(m.mask, XI_RawTouchEnd);

    XISelectEvents(display, DefaultRootWindow(display), &m, 1);

    free(m.mask);
     XSync(display, False);
 }

int listen(Display *display)
{
    Window root = DefaultRootWindow(display);
    int root_x, root_y, nouse;
    Window noused_window;
    unsigned int mask;
    while(1)
    {
	XEvent ev;
	XGenericEventCookie *cookie = (XGenericEventCookie*)&ev.xcookie;
	XNextEvent(display, (XEvent*)&ev);

	if (XGetEventData(display, cookie) &&
		cookie->type == GenericEvent &&
		cookie->extension == xi2_opcode)
	{
	    switch (cookie->evtype)
	    {
		case XI_RawButtonPress:
		case XI_RawButtonRelease:
		case XI_RawMotion:
		case XI_RawKeyPress:
		case XI_RawKeyRelease:
		case XI_RawTouchBegin:
		case XI_RawTouchUpdate:
		case XI_RawTouchEnd:
		    {
		    XQueryPointer(display, root, &noused_window, &noused_window, &root_x, &root_y, &nouse, &nouse, &mask);
		    XIRawEvent* event = cookie->data;
		    go_handle_raw_event(cookie->evtype, event->detail, root_x, root_y, mask);
		    }

		    break;
	    }
	}

	XFreeEventData(display, cookie);
    }
    return EXIT_SUCCESS;
}

int start_listen()
{
    Display* display = XOpenDisplay(NULL);
    int event, error;

    if (!XQueryExtension(display, "XInputExtension", &xi2_opcode, &event, &error)) {
	fprintf(stderr, "XInput2 not available.\n");
	return -1;
    }

    if (!xinput_version(display)) {
	fprintf(stderr, "XInput2 extension not available\n");
	return -1;
    }

    select_events(display);
    listen(display);
    return 0;
}

int xi_mask_is_set(unsigned char*ptr, char mask)
{
    return XIMaskIsSet(ptr, mask);
}
