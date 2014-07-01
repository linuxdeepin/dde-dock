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

#ifndef __UTILS_H__
#define __UTILS_H__

#include <glib.h>
#include <gdk/gdk.h>
#include <gdk/gdkx.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/Xatom.h>

#include <X11/extensions/XInput.h>
#include <X11/extensions/XIproto.h>

typedef struct _DeviceInfo {
	char *name;
	char *atom_name;
	int xid;
	int atom;
} DeviceInfo;

gchar *str_to_upper(const gchar *src);
gchar *str_to_letter(const gchar *src);
gboolean str_is_contain (const gchar *src, const gchar *sub);
int xi_device_exist (const char *name);
GdkDevice* device_is_exist (const char *deviceName);
gboolean set_device_enabled (int device_id,
                    gboolean enabled);
XDevice *open_gdk_device (GdkDevice *device);

#endif
