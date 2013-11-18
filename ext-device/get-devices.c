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

#include "get-devices.h"
#include <gdk/gdk.h>

static gboolean SubStrContainer (const gchar *src, const gchar *sub);

int
DeviceIsExist (const char *deviceName)
{
	int flag = 0;
	GList *devList, *l;
	GdkDisplay *display = gdk_display_get_default ();
	GdkDeviceManager *devManager = gdk_display_get_device_manager(display);

	devList = gdk_device_manager_list_devices(devManager, 
			GDK_DEVICE_TYPE_SLAVE);
	for ( l = devList; l != NULL; l = l->next ) {
		GdkDevice *device = l->data;

		const gchar *name = gdk_device_get_name(device);
		if ( SubStrContainer (name, deviceName) ) {
			flag = 1;
			break;
		}
	}
	g_list_free (devList);

	return flag;
}

static gboolean
SubStrContainer (const gchar *src, const gchar *sub)
{
	if ( src == NULL || sub == NULL ) {
		return FALSE;
	}

	if ( g_strrstr (src, sub) == NULL ) {
		return FALSE;
	}

	return TRUE;
}
