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
#include <X11/Xlib.h>
#include <X11/extensions/XInput2.h>

#include "utils.h"
#include "devices.h"

static int set_prop_float (Display *disp, int deviceid, 
		const char *prop_name, double value);

int
is_mouse_device(int deviceid)
{
	int ret = is_device_property_exist(deviceid, "Button Labels");
	if (ret == 1) {
		return 1;
	}

	return 0;
}

int
set_motion_acceleration (int deviceid, double acceleration)
{
	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed: %d\n", deviceid);
		return -1;
	}

	int ret = set_prop_float(disp, deviceid, 
			"Device Accel Constant Deceleration", acceleration);
	XCloseDisplay(disp);

	return ret;
}

int
set_motion_threshold(int deviceid, double threshold)
{
	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed: %d\n", deviceid);
		return -1;
	}

	int ret = set_prop_float(disp, deviceid, 
			"Device Accel Adaptive Deceleration", threshold);
	XCloseDisplay(disp);

	return ret;
}

static int
set_prop_float (Display *disp, int deviceid, 
		const char *prop_name, double value)
{
	if (!disp || !prop_name) {
		fprintf(stderr, "Args error in set_prop_float\n");
		return -1;
	}

	Atom prop = XInternAtom(
			disp,
			prop_name, 
			True);
	if (prop == None) {
		fprintf(stderr, "Get '%s' Atom Failed\n", prop_name);
		return -1;
	}

	Atom act_type;
	int act_format;
	unsigned long nitems, bytes_after;
	unsigned char *data;

	int rc = XIGetProperty(disp, deviceid, prop, 0, 1, False, 
			AnyPropertyType, &act_type, &act_format,
			&nitems, &bytes_after, &data);
	if (rc != Success) {
		fprintf(stderr, "Get '%ld' property failed\n", prop);
		return -1;
	}

	Atom float_atom = XInternAtom(disp, "FLOAT", True);
	if (float_atom != None && act_type == float_atom && nitems == 1) {
		*(float*)data = (float)value;
		XIChangeProperty(disp, deviceid, prop, act_type, act_format,
				XIPropModeReplace, data, nitems);
	}
	XFree(data);

	return 0;
}
