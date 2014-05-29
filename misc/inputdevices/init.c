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

#include "utils.h"
#include "devices.h"
#include <gdk/gdk.h>

static void listen_device_changed();
static void device_removed_cb(GdkDeviceManager *manager,
                              GdkDevice *device, gpointer user_data);

void
init_gdk_env ()
{
	gdk_init(NULL, NULL);
        listen_device_changed();
}

static void
listen_device_changed ()
{
	GdkDeviceManager *manager = gdk_display_get_device_manager(
	                                gdk_display_get_default());
	if (manager == NULL) {
		g_warning("Get Devices Manager Failed");
		return;
	}

	g_signal_connect (G_OBJECT(manager), "device-removed",
	                  G_CALLBACK(device_removed_cb), NULL);
}

static void
device_removed_cb(GdkDeviceManager *manager,GdkDevice *device,
                  gpointer user_data)
{
	const gchar *name = gdk_device_get_name(device);
	if ( str_is_contain (name, KEYBOARD_KEY_NAME) ) {
		set_tpad_enable(1);
	}
}
