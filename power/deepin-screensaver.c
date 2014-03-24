/* -*- Mode: C; tab-width: 8; indent-tabs-mode: t; c-basic-offset: 8 -*-
 *
 * Copyright (C) 2013-2014 Linux Deepin
 * Copyright (C) 2013-2014 onerhao <onerhao@gmail.com>
 *
 * Licensed under the GNU General Public License Version 2
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

#include "config.h"

#include <string.h>
#include <unistd.h>
#include <stdio.h>
#include <glib.h>
#include <gio/gio.h>

#include "gnome-settings-session.h"

#define DEEPIN_SCREENSAVER_NAME         "com.deepin.daemon.Power"
#define DEEPIN_SCREENSAVER_OBJECT       "/org/freedesktop/ScreenSaver"
#define DEEPIN_SCREENSAVER_INTERFACE    "org.freedesktop.ScreenSaver"

GDBusProxy *
deepin_screensaver_get_proxy (void)
{
        static GDBusProxy *ss_proxy;
        GError *error =  NULL;

        if (ss_proxy != NULL) {
                g_object_ref (ss_proxy);
        } else {
                ss_proxy = g_dbus_proxy_new_for_bus_sync (G_BUS_TYPE_SESSION,
                                                               G_DBUS_PROXY_FLAGS_NONE,
                                                               NULL,
                                                               DEEPIN_SCREENSAVER_NAME,
                                                               DEEPIN_SCREENSAVER_OBJECT,
                                                               DEEPIN_SCREENSAVER_INTERFACE,
                                                               NULL,
                                                               &error);
                if (error) {
                        g_warning ("Failed to connect to the session manager: %s", error->message);
                        g_error_free (error);
                } else {
                        g_object_add_weak_pointer (G_OBJECT (ss_proxy), (gpointer*)&ss_proxy);
                }
        }

        return ss_proxy;
}
