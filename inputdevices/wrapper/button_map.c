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
#include <string.h>
#include <X11/Xlib.h>
#include <X11/extensions/XInput.h>

#include "devices.h"

static int get_device_button_num(Display *disp, const char *name);
static unsigned char* get_button_map (XID xid, const char *name,
                                      int *nbuttons);
static int set_button_map(XID xid, const char *name,
                          unsigned char *map, int nbuttons);
static int do_button_map_set (Display *disp, XDevice *dev, 
		unsigned char *map, int nbuttons);

int
set_left_handed(unsigned long xid, const char *name, int enabled)
{
	if (!name) {
		fprintf(stderr, "Args error in set_left_handed\n");
		return -1;
	}

	int nbuttons = 0;
	unsigned char *map = get_button_map(xid, name, &nbuttons);
	if (!map) {
		/*fprintf(stderr, "Get button map failed: %s\n", name);*/
		return -1;
	}

	if (nbuttons >= 3) {
		if (enabled) {
			map[0] = 3;
			map[2] = 1;
		} else {
			map[0] = 1;
			map[2] = 3;
		}

		set_button_map(xid, name, map, nbuttons);
	}

	free(map);

	return 0;
}

int
set_mouse_natural_scroll(unsigned long xid, const char *name, int enabled)
{
	if (!name) {
		fprintf(stderr, "Args error in set_mouse_natural_scroll\n");
		return -1;
	}

	int nbuttons = 0;
	unsigned char *map = get_button_map(xid, name, &nbuttons);
	if (!map) {
		/*fprintf(stderr, "Get button map failed: %s\n", name);*/
		return -1;
	}

	if (nbuttons >= 5) {
		if (enabled) {
			map[3] = 5;
			map[4] = 4;
		} else {
			map[3] = 4;
			map[4] = 5;
		}

		set_button_map(xid, name, map, nbuttons);
	}

	free(map);

	return 0;
}

static int
get_device_button_num(Display *disp, const char *name)
{
	if (!disp || !name) {
		fprintf(stderr, "Args error in find_device_info\n");
		return -1;
	}

	XDeviceInfo *devices = NULL;
	XDeviceInfo *found = NULL;

	int n_devices;
	devices = XListInputDevices(disp, &n_devices);
	if (!devices) {
		fprintf(stderr, "List Input Devices Failed\n");
		return -1;
	}

	int i;
	for (i = 0; i < n_devices; i++) {
		if (devices[i].use != IsXExtensionPointer) {
			continue;
		}

		if (strcmp(name, devices[i].name) == 0) {
			found = &devices[i];
			break;
		}
	}

	int nbuttons = -1;
	if (found) {
		XAnyClassPtr ip = (XAnyClassPtr)found->inputclassinfo;
		int j;

		for (j = 0; j < found->num_classes; j++) {
			if (ip->class == ButtonClass) {
				nbuttons = ((XButtonInfoPtr)ip)->num_buttons;
				break;
			}
			ip = (XAnyClassPtr)((char*)ip + ip->length);
		}
	}
	XFreeDeviceList(devices);

	return nbuttons;
}

static unsigned char*
get_button_map (XID xid, const char *name, int *nbuttons)
{
	if (!name || !nbuttons) {
		fprintf(stderr, "Args error in get_button_map\n");
		return NULL;
	}

	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed\n");
		return NULL;
	}

	int num = get_device_button_num(disp, name);
	if (num <= 0) {
		fprintf(stderr, "Get button number failed\n");
		XCloseDisplay(disp);
		return NULL;
	}

	unsigned char *map;
	map = (unsigned char*)calloc(num, sizeof(unsigned char));
	if (!map) {
		fprintf(stderr, "Alloc memory for button map failed\n");
		XCloseDisplay(disp);
		return NULL;
	}

	XDevice *dev = XOpenDevice(disp, xid);
	if (!dev) {
		fprintf(stderr, "Open Device '%ld' Failed\n", xid);
		XCloseDisplay(disp);
		return NULL;
	}

	/*XGetDeviceButtonMapping(disp, dev, map, num);*/
	int rc = XGetDeviceButtonMapping(disp, dev, map, num);
	if (rc != num) {
		fprintf(stderr, "Get '%s' Button Map Failed\n", name);
		free(map);
		map = NULL;
		XCloseDevice(disp, dev);
		XCloseDisplay(disp);
		return NULL;
	}

	XCloseDevice(disp, dev);
	XCloseDisplay(disp);
	*nbuttons = num;

	return map;
}

static int
set_button_map(XID xid, const char *name, unsigned char *map, int nbuttons)
{
	if (!name || !map) {
		fprintf(stderr, "Args error in set_button_map\n");
		return -1;
	}

	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed\n");
		return -1;
	}

	XDevice *dev = XOpenDevice(disp, xid);
	if (!dev) {
		fprintf(stderr, "Open Device '%ld' Failed\n", xid);
		XCloseDisplay(disp);
		return -1;
	}

	int rc = do_button_map_set(disp, dev, map, nbuttons);
	if (rc == -1) {
		//failed
		fprintf(stderr, "Set '%s' button map failed\n", name);
		XCloseDevice(disp, dev);
		XCloseDisplay(disp);
		return -1;
	}

	XCloseDevice(disp, dev);
	XCloseDisplay(disp);

	return 0;
}

static int
do_button_map_set (Display *disp, XDevice *dev, 
		unsigned char *map, int nbuttons) 
{
	if (!disp || !dev || !map) {
		return -1;
	}

	int rc = XSetDeviceButtonMapping(disp, dev, map, nbuttons);
	if (rc != 0) {
		do_button_map_set(disp, dev, map, nbuttons);
	}

	return 0;
}
