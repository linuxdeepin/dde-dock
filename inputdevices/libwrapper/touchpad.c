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
#include <stdlib.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/extensions/XInput2.h>

#include "utils.h"
#include "devices.h"

/**
 * Detailed description of the touch panel properties,
 * see here: http://www.x.org/archive/X11R7.5/doc/man/man4/synaptics.4.html
 **/

int
is_tpad_device(int deviceid)
{
	int ret = is_device_property_exist(deviceid, "Synaptics Off");
	if (ret == 1) {
		return 1;
	}

	return 0;
}

int
enabled_touchpad (int deviceid, int enabled)
{
	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed: %d\n", deviceid);
		return -1;
	}

	if (enabled_device(disp, deviceid, enabled) != 0) {
		fprintf(stderr, "Enable touchpad failed\n");
		XCloseDisplay(disp);
		return -1;
	}

	if (enabled) {
		/**
		 * Property 'Synaptics Off' 8 bit, valid values (0, 1, 2):
		 *	Value 0: Touchpad is enabled
		 *	Value 1: Touchpad is switched off
		 *	Value 2: Only tapping and scrolling is switched off
		 **/
		Atom prop = XInternAtom(
		                disp,
		                "Synaptics Off",
		                True);
		if (prop != None) {
			unsigned char data = 0;
			XIChangeProperty(disp, deviceid, prop, XA_INTEGER, 8,
			                 XIPropModeReplace, &data, 1);
		}
	}

	XCloseDisplay(disp);
	return 0;
}

/**
 * Property "Synaptics Scrolling Distance" 32 bit, 2 values, vert, horiz.
 *	Option "VertScrollDelta" "integer":
 *		Move distance of the finger for a scroll event.
 *	Option "HorizScrollDelta" "integer" :
 *		Move distance of the finger for a scroll event.
 *
 * if delta = 0, use value from property getting
 **/
int
set_touchpad_natural_scroll(int deviceid, int enabled, int delta)
{
	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed: %d\n", deviceid);
		return -1;
	}

	Atom prop = XInternAtom(disp,
	                        "Synaptics Scrolling Distance",
	                        True);
	if (prop == None) {
		fprintf(stderr, "Get 'Synaptics Scrolling Distance' Atom Failed\n");
		XCloseDisplay(disp);
		return -1;
	}

	Atom act_type;
	int act_format;
	unsigned long nitems, bytes_after;
	unsigned char *data = NULL;
	int rc = XIGetProperty(disp, deviceid, prop, 0, 2, False,
	                       XA_INTEGER, &act_type, &act_format,
	                       &nitems, &bytes_after, &data);
	if (rc != Success) {
		fprintf(stderr, "Get 'Synaptics Scrolling Distance' Property Failed\n");
		XCloseDisplay(disp);
		return -1;
	}

	if (act_type == XA_INTEGER && act_format == 32 && nitems >= 2) {
		int *ptr = (int*)data;
		if (enabled) {
			if (delta > 0) {
				ptr[0] = -abs(delta);
				ptr[1] = -abs(delta);
			} else {
				ptr[0] = -abs(ptr[0]);
				ptr[1] = -abs(ptr[1]);
			}
		} else {
			if (delta > 0) {
				ptr[0] = abs(delta);
				ptr[1] = abs(delta);
			} else {
				ptr[0] = abs(ptr[0]);
				ptr[1] = abs(ptr[1]);
			}
		}

		XIChangeProperty(disp, deviceid, prop, act_type, act_format,
		                 XIPropModeReplace, data, nitems);
	}
	XFree(data);
	XCloseDisplay(disp);

	return 0;
}

/**
 * Property "Synaptics Edge Scrolling" 8 bit (BOOL), 3 values, vertical,
 * horizontal, corner. :
 *	Option "VertEdgeScroll" "boolean":
 *		Enable vertical scrolling when dragging along the right edge.
 *	Option "HorizEdgeScroll" "boolean" :
 *		Enable horizontal scrolling when dragging along
 *		the bottom edge.
 *	Option "CornerCoasting" "boolean":
 *		Enable edge scrolling to continue while the finger stays
 *		in an edge corner.
 **/
int
set_edge_scroll(int deviceid, int enabled)
{
	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed: %d\n", deviceid);
		return -1;
	}

	Atom prop = XInternAtom(disp,
	                        "Synaptics Edge Scrolling",
	                        True);
	if (prop == None) {
		fprintf(stderr, "Get 'Synaptics Edge Scrolling' Atom Failed\n");
		XCloseDisplay(disp);
		return -1;
	}

	Atom act_type;
	int act_format;
	unsigned long nitems, bytes_after;
	unsigned char *data = NULL;
	int rc = XIGetProperty(disp, deviceid, prop, 0, 1, False,
	                       XA_INTEGER, &act_type, &act_format,
	                       &nitems, &bytes_after, &data);
	if (rc != Success) {
		fprintf(stderr, "Get 'Synaptics Edge Scrolling' Property Failed\n");
		XCloseDisplay(disp);
		return -1;
	}

	if (act_type == XA_INTEGER && act_format == 8 && nitems >= 3) {
		data[0] = enabled ? 1 : 0;
		data[1] = enabled ? 1 : 0;
		data[2] = enabled ? 1 : 0;

		XIChangeProperty(disp, deviceid, prop, act_type, act_format,
		                 XIPropModeReplace, data, nitems);
	}

	XFree(data);
	XCloseDisplay(disp);

	return 0;
}

/**
 * Property 'Synaptics Two-Finger Scrolling' 8 bit (BOOL),
 * 2 values, vertical, horizontal.
 *	Option "VertTwoFingerScroll" "boolean":
 *		Enable vertical scrolling when dragging with
 *		two fingers anywhere on the touchpad.
 *	Option "HorizTwoFingerScroll" "boolean" :
 *		Enable horizontal scrolling when dragging with
 *		two fingers anywhere on the touchpad.
 **/
int
set_two_finger_scroll(int deviceid, int vert_enabled, int horiz_enabled)
{
	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed: %d\n", deviceid);
		return -1;
	}

	Atom prop = XInternAtom(disp,
	                        "Synaptics Two-Finger Scrolling",
	                        True);
	if (prop == None) {
		fprintf(stderr, "Get 'Synaptics Two-Finger Scrolling' Atom Failed\n");
		XCloseDisplay(disp);
		return -1;
	}

	Atom act_type;
	int act_format;
	unsigned long nitems, bytes_after;
	unsigned char *data = NULL;
	int rc = XIGetProperty(disp, deviceid, prop, 0, 1, False,
	                       XA_INTEGER, &act_type, &act_format,
	                       &nitems, &bytes_after, &data);
	if (rc != Success) {
		fprintf(stderr, "Get 'Synaptics Two-Finger Scrolling' Property Failed\n");
		XCloseDisplay(disp);
		return -1;
	}

	if (act_type == XA_INTEGER && act_format == 8 && nitems >= 2) {
		data[0] = vert_enabled ? 1 : 0;  // set vertical(垂直) scroll
		data[1] = horiz_enabled ? 1 : 0; // set horizon(水平) scroll

		XIChangeProperty(disp, deviceid, prop, act_type, act_format,
		                 XIPropModeReplace, data, nitems);
	}

	XFree(data);
	XCloseDisplay(disp);

	return 0;
}

/**
 * Property 'Synaptics Tap Action' 8 bit,
 * up to MAX_TAP values (see synaptics.h), 0 disables an element.
 * order: RT, RB, LT, LB, F1, F2, F3.
 *	Option "RTCornerButton" "integer":
 *		Which mouse button is reported on a right top corner tap.
 *	Option "RBCornerButton" "integer":
 *		Which mouse button is reported on a right bottom corner tap.
 *	Option "LTCornerButton" "integer":
 *		Which mouse button is reported on a left top corner tap.
 *	Option "LBCornerButton" "integer":
 *		Which mouse button is reported on a left bottom corner tap.
 *	Option "TapButton1" "integer":
 *		Which mouse button is reported on a non-corner one-finger tap.
 *	Option "TapButton2" "integer":
 *		Which mouse button is reported on a non-corner two-finger tap.
 *	Option "TapButton3" "integer":
 *		Which mouse button is reported on a non-corner
 *		three-finger tap.
 **/
int
set_tab_to_click(int deviceid, int enabled, int left_handed)
{
	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed: %d\n", deviceid);
		return -1;
	}

	Atom prop = XInternAtom(disp,
	                        "Synaptics Tap Action",
	                        True);
	if (prop == None) {
		fprintf(stderr, "Get 'Synaptics Tap Action' Atom Failed\n");
		XCloseDisplay(disp);
		return -1;
	}

	Atom act_type;
	int act_format;
	unsigned long nitems, bytes_after;
	unsigned char *data = NULL;
	int rc = XIGetProperty(disp, deviceid, prop, 0, 2, False,
	                       XA_INTEGER, &act_type, &act_format,
	                       &nitems, &bytes_after, &data);
	if (rc != Success) {
		fprintf(stderr, "Get 'Synaptics Tap Action' Property Failed\n");
		XCloseDisplay(disp);
		return -1;
	}

	if (act_type == XA_INTEGER && act_format == 8 && nitems >= 7) {
		data[4] = enabled ? (left_handed ? 3 : 1) : 0; // left button
		data[5] = enabled ? (left_handed ? 1 : 3) :0; // right button
		data[6] = enabled ? 2 : 0; // middle button

		XIChangeProperty(disp, deviceid, prop, act_type, act_format,
		                 XIPropModeReplace, data, nitems);
	}

	XFree(data);
	XCloseDisplay(disp);

	return 0;
}
