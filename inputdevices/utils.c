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
#include <X11/Xatom.h>
#include <X11/extensions/XInput2.h>

#include "utils.h"
#include "devices.h"

DeviceInfo*
get_device_info_list(int *n_devices)
{
	if (!n_devices) {
		return NULL;
	}

	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed\n");
		return NULL;
	}

	int num;
	XIDeviceInfo *infos = XIQueryDevice(disp, XIAllDevices, &num);
	if (!infos) {
		fprintf(stderr, "List Input Device Failed\n");
		XCloseDisplay(disp);
		return NULL;
	}

	int i;
	int j = 0;
	DeviceInfo *list = NULL;
	for (i = 0; i < num; i++) {
		if ((infos[i].use != XISlavePointer &&
		        infos[i].use != XISlaveKeyboard)) {
			continue;
		}

		DeviceInfo *tmp = calloc(j+1, sizeof(DeviceInfo));
		if (!tmp) {
			fprintf(stderr, "Alloc memory failed\n");
			continue;
		}

		if (j != 0) {
			memcpy(tmp, list, j * sizeof(DeviceInfo));
		}

		unsigned long size = strlen(infos[i].name);
		tmp[j].name = calloc(size+1, sizeof(char));
		if (!tmp[j].name) {
			fprintf(stderr, "Alloc memory for name failed\n");
			continue;
		}
		memcpy(tmp[j].name, infos[i].name, size);
		/*tmp[j].name = infos[i].name;*/
		tmp[j].deviceid = infos[i].deviceid;
		tmp[j].enabled = (int)infos[i].enabled;

		if (j != 0) {
			free(list);
			list = NULL;
		}
		list = tmp;
		tmp = NULL;
		j++;
	}
	XIFreeDeviceInfo(infos);
	XCloseDisplay(disp);
	*n_devices = j;

	return list;
}

void
free_device_info(DeviceInfo *infos, int n_devices)
{
	if (!infos) {
		return ;
	}

	int i;
	for (i = 0; i < n_devices; i++) {
		free(infos[i].name);
		infos[i].name = NULL;
	}

	free(infos);
	infos = NULL;
}

/**
 * return:
 *	-1: error or not exist
 *	0: success
 */
int
set_device_enabled (Display *disp, int deviceid, int enabled)
{
	if (!disp) {
		fprintf(stderr, "Display NULL\n");
		return -1;
	}

	Atom prop = XInternAtom(
	                disp,
	                "Device Enabled",
	                True);
	if (prop == None) {
		fprintf(stderr, "Get 'Device Enabled' Atom Failed\n");
		return -1;
	}

	unsigned char data = enabled ? 1 : 0;

	XIChangeProperty(disp, deviceid, prop, XA_INTEGER, 8,
	                 XIPropModeReplace, &data, 1);

	return 0;
}

/**
 * return:
 *	-1: error
 *	0: not exist
 *	1: exist
 */
int
is_device_property_exist(Display *disp, int deviceid, const char *prop_name)
{
	if (!disp || !prop_name) {
		fprintf(stderr, "Args error\n");
		return -1;
	}

	int nprops;
	Atom *props = XIListProperties(disp, deviceid, &nprops);
	if (!props) {
		fprintf(stderr, "List Device '%s' Properties Failed\n",
		        prop_name);
		return -1;
	}

	int flags = 0;
	while (nprops--) {
		char *name = XGetAtomName(disp, props[nprops]);
		if (strcmp(name, prop_name) == 0 ) {
			flags = 1;
			XFree(name);
			break;
		}
		XFree(name);
	}
	XFree(props);

	if (flags) {
		return 1;
	}

	return 0;
}

int
is_wacom_device(int deviceid)
{
	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed: %d\n", deviceid);
		return -1;
	}

	int ret = is_device_property_exist(disp, deviceid, "Wacom Tool Type");

	XCloseDisplay(disp);

	if (ret == 1) {
		return 1;
	}

	return 0;
}
