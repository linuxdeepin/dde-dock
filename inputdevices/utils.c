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

#include "utils.h"
#include <X11/extensions/XInput2.h>

int
xi_device_exist (const char *name)
{
	if (name == NULL) {
		return -1;
	}

	int n_devices;
	XDeviceInfo *infos = XListInputDevices(
	                         GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()),
	                         &n_devices);
	int i;
	int dev_id = -1;

	for (i = 0; i < n_devices; i++) {
		/*g_debug("XID: %lu\n", infos[i].id);*/
		/*g_debug("Atom: %lu\n", infos[i].type);*/
		/*g_debug("Name: %s\n", infos[i].name);*/
		/*g_debug("Num Class: %d\n", infos[i].num_classes);*/
		/*g_debug("Use: %d\n", infos[i].use);*/
		// Filter master device
		if ((g_strcmp0(infos[i].name, "Virtual core pointer") == 0) ||
		        (g_strcmp0(infos[i].name, "Virtual core keyboard") == 0)) {
			continue;
		}

		if ( str_is_contain (infos[i].name, name) ) {
			dev_id = infos[i].id;
			break;
		}
	}

	XFreeDeviceList(infos);

	return dev_id;
}


GdkDevice *
device_is_exist (const char *deviceName)
{
	g_debug("Check Device Exisr: %s\n", deviceName);
	GList *devList, *l;
	GdkDisplay *display = gdk_display_get_default ();

	if (display == NULL) {
		g_warning("Get Default Display Failed: %s", deviceName);
		return NULL;
	}

	g_debug("Get Device Manager\n");
	GdkDeviceManager *devManager = gdk_display_get_device_manager(display);

	if (devManager == NULL) {
		g_warning("Get Device Manager Failed: %s", deviceName);
		return NULL;
	}

	g_debug("Get Device List\n");
	devList = gdk_device_manager_list_devices(devManager,
	          GDK_DEVICE_TYPE_SLAVE);

	if (devList == NULL) {
		g_warning("Get Device List Failed: %s", deviceName);
		return NULL;
	}

	g_debug("Get Device List End\n");
	GdkDevice *device = NULL;

	gboolean flag = FALSE;

	for ( l = devList; l != NULL; l = l->next ) {
		device = l->data;

		const gchar *name = gdk_device_get_name(device);

		g_debug("Device Name: %s\n", name);

		if ( str_is_contain (name, deviceName) ) {
			flag = TRUE;
			break;
		}
	}

	g_list_free (devList);

	if (flag) {
		flag = FALSE;
		return device;
	}

	return NULL;
}

gboolean
str_is_contain (const gchar *src, const gchar *sub)
{
	if ( src == NULL || sub == NULL ) {
		return FALSE;
	}

	gchar *tmp1 = str_to_letter(src);
	gchar *tmp2 = str_to_letter(sub);

	gchar *ret = g_strrstr (tmp1, tmp2);
	g_free(tmp1);
	g_free(tmp2);

	if ( ret == NULL ) {
		return FALSE;
	}

	return TRUE;
}

gchar *
str_to_upper(const gchar *src)
{
	if (src == NULL) {
		return NULL;
	}

	/*g_debug("To Upper: %s\n", src);*/
	return g_utf8_strup(src, -1);
}

gchar *
str_to_letter(const gchar *src)
{
	if (src == NULL) {
		return NULL;
	}

	/*g_debug("To Letter: %s\n", src);*/
	return g_utf8_strdown(src, -1);
}

gboolean
set_device_enabled (int device_id,
                    gboolean enabled)
{
	Atom prop;
	guchar value;

	prop = XInternAtom (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), "Device Enabled", False);

	if (!prop) {
		return FALSE;
	}

	gdk_error_trap_push ();
	g_debug("Start Set device\n");

	value = enabled ? 1 : 0;
	XIChangeProperty (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()),
	                  device_id, prop, XA_INTEGER, 8, PropModeReplace, &value, 1);

	g_debug("Has Set end\n");

	if (gdk_error_trap_pop ()) {
		return FALSE;
	}

	return TRUE;
}

XDevice *
open_gdk_device (GdkDevice *device)
{
	XDevice *xdevice;
	int id;

	g_object_get (G_OBJECT (device), "device-id", &id, NULL);

	gdk_error_trap_push ();

	xdevice = XOpenDevice (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), id);

	if (gdk_error_trap_pop () != 0) {
		return NULL;
	}

	return xdevice;
}
