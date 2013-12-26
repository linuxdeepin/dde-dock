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

#include "list-devices-info.h"

static gint xinput_version(Display *disp);
static void list_devices(Display *disp);
static void list_xi2 (Display *disp);
static void list_xi1 (Display *disp);
#ifdef HAVE_XI2
static void list_class_xi2 (XIAnyClassInfo **classes, int num_classes);
#endif

static GArray *name_array = NULL;

static gint
xinput_version (Display *disp)
{
    XExtensionVersion	*version;
    static int vers = -1;

    if (vers != -1) {
        return vers;
    }

    version = XGetExtensionVersion(disp, INAME);

    if (version && (version != (XExtensionVersion *) NoSuchExtension)) {
        vers = version->major_version;
        XFree(version);
    }

#if HAVE_XI2

    /* Announce our supported version so the server treats us correctly. */
    if (vers >= XI_2_Major) {
        int maj = 2,
            min = 0;

#if HAVE_XI21
        min = 1;
#elif HAVE_XI22
        min = 2;
#endif

        XIQueryVersion(display, &maj, &min);
    }

#endif

    return vers;
}

static void
list_devices(Display *disp)
{
#ifdef HAVE_XI2

    if ( xinput_version(disp) == XI_2_Major ) {
        list_xi2(disp);
        return;
    }

#endif

    list_xi1(disp);
}

#ifdef HAVE_XI2
static void
list_xi2 (Display *disp)
{
    XIDeviceInfo *info, *dev;
    int i;
    int num_devices;

    info = XIQueryDevice (disp, XIAllDevices, &num_devices);

    for ( i = 0; i < num_devices; i++ ) {
        dev = &info[i];
        /*g_print ("%s\t%d\n", dev->name, dev->deviceid);*/
        list_class_xi2 (dev->classes, dev->num_classes);
    }

    XIFreeDeviceInfo (info);
}

static void
list_class_xi2 (XIAnyClassInfo **classes, int num_classes)
{
    int i = 0;

    for ( i = 0; i < num_classes; i++ ) {
        Atom type = classes[i].type;

        if (type > 0) {
            gchar *name = g_strdup(XGetAtomName(disp, type));

            if ( !name ) {
                continue;
            }

            g_array_append_val (name_array, name);
        }
    }
}
#endif

static void
list_xi1 (Display *disp)
{
    XDeviceInfo *info;
    int loop;
    int num_devices;

    info = XListInputDevices (disp, &num_devices);

    for ( loop = 0; loop < num_devices; loop++ ) {
        XDeviceInfo *dev = info + loop;
        Atom type = dev->type;

        if (type > 0) {
            gchar *name = g_strdup(XGetAtomName(disp, type));

            if ( !name ) {
                continue;
            }

            g_array_append_val (name_array, name);
        }
    }
}

int
find_device_by_name (char *name)
{
    Display *disp = XOpenDisplay (NULL);

    if ( !disp ) {
        g_warning ("Unable to connect to X server!");
        goto out;
    }

    int xi_opcode, event, error;

    if ( !XQueryExtension(disp, "XInputExtension",
                          &xi_opcode, &event, &error) ) {
        g_warning ("X Input extension not available.");
        goto out;
    }

    if ( !xinput_version (disp) ) {
        g_warning ("%s extension not avaliable.", INAME);
        goto out;
    }

    name_array = g_array_new (TRUE, TRUE, sizeof(char *));
    list_devices (disp);
    guint len = g_array_get_element_size (name_array);
    int i;

    for ( i = 0; i < len; i++ ) {
        gchar *tmp = g_array_index(name_array, char *, i);
        g_print("name: %s\n", tmp);

        if (g_strcmp0(tmp, name) == 0) {
            g_print ("equal...\n");
            g_array_free (name_array, TRUE);
            return 1;
        }
    }

    g_array_free (name_array, TRUE);

out:

    if ( disp ) {
        XCloseDisplay (disp);
    }

    return 0;
}
