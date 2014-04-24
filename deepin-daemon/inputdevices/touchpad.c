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

#include "devices.h"
#include "utils.h"
#include <stdlib.h>

#define DEVICE_PROP_ID "device-id"

void
set_tpad_enable(int enable)
{
    int id = xi_device_exist(TPAD_NAME_KEY);

    if (id == -1) {
        g_warning("Get Touchpad Device Id Failed");
        return;
    }

    g_print("Trying to set %s, id: %d\n",
            enable ? "Enable TouchPad" : "Disable TouchPad", id);

    if (set_device_enabled (id, enable) == FALSE) {
        g_warning ("Set %s Failed: id (%d)",
                   enable ? "Enable TouchPad" : "Disable TouchPad", id);
    } else {
        g_print ("Set %s Success: id (%d)\n",
                 enable ? "Enable TouchPad" : "Disable TouchPad", id);
    }
}

void
set_natural_scroll(int enable)
{
    GdkDevice *tpad = device_is_exist(TPAD_NAME_KEY);

    if (tpad == NULL) {
        g_warning("TouchPad not exist\n");
        return;
    }

    XDevice *xdev = NULL;
    xdev = open_gdk_device(tpad);

    if (xdev == NULL) {
        g_warning("Get XDevice From TouchPad");
        return;
    }

    g_print("Trying to set %s for \"%s\"\n",
            enable ? "natural (reverse) scrollroll" : "normal scroll",
            gdk_device_get_name(tpad));

    Atom scrolling_distance = XInternAtom(
                                  GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                                  "Synaptics Scrolling Distance", FALSE);
    g_print("Error Trap Push\n");
    gdk_error_trap_push();
    /*gdk_error_trap_pop_ignored();*/
    g_print("Get Device Property\n");
    Atom act_type;
    int act_format;
    unsigned long nitems, bytes_after;
    unsigned char *data = NULL;
    glong *ptr = NULL;
    int rc = XGetDeviceProperty (
                 GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                 xdev, scrolling_distance, 0, 2, FALSE,
                 XA_INTEGER, &act_type, &act_format, &nitems,
                 &bytes_after, &data);

    if (rc == Success && act_type == XA_INTEGER &&
            act_format == 32 && nitems >= 2) {
        ptr = (glong *)data;

        if (enable) {
            ptr[0] = -abs(ptr[0]);
            ptr[1] = -abs(ptr[1]);
        } else {
            ptr[0] = abs(ptr[0]);
            ptr[1] = abs(ptr[1]);
        }

        XChangeDeviceProperty (GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                               xdev, scrolling_distance, XA_INTEGER, act_format,
                               PropModeReplace, data, nitems);
    }

    if (rc == Success) {
        XFree(data);
    }

    if (gdk_error_trap_pop()) {
        g_warning("Error settings touchpad %s for \"%s\"\n",
                  enable ? "natural (reverse) scrollroll" : "normal scroll",
                  gdk_device_get_name(tpad));
    }

    XCloseDevice(GDK_DISPLAY_XDISPLAY(gdk_display_get_default()), xdev);
}

void
set_edge_scroll(int enable)
{
    Atom prop_edge = XInternAtom(
                         GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                         "Synaptics Edge Scrolling", FALSE);

    if (!prop_edge) {
        g_warning("Get Edge Prop Atom Failed");
        return;
    }

    GdkDevice *tpad = device_is_exist(TPAD_NAME_KEY);

    if (tpad == NULL) {
        g_warning("TouchPad not exist\n");
        return;
    }

    XDevice *xdev = NULL;
    xdev = open_gdk_device(tpad);

    if (xdev == NULL) {
        g_warning("Get XDevice From TouchPad");
        return;
    }

    g_print("Trying to set %s\n",
            enable ? "enable edge scroll" : "disable edge scroll");
    gdk_error_trap_push();
    Atom act_type;
    int act_format;
    unsigned long nitems, bytes_after;
    unsigned char *data = NULL;
    int rc = XGetDeviceProperty(GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                                xdev, prop_edge, 0, 1, False,
                                XA_INTEGER, &act_type, &act_format, &nitems,
                                &bytes_after, &data);

    if (rc == Success && act_type == XA_INTEGER &&
            act_format == 8 && nitems >= 2) {
        data[0] = enable ? 1 : 0;
        XChangeDeviceProperty(GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                              xdev, prop_edge, XA_INTEGER, 8,
                              PropModeReplace, data, nitems);
    }

    if (rc == Success) {
        XFree(data);
    }

    if (gdk_error_trap_pop()) {
        g_warning("Error settings touchpad %s\n",
                  enable ? "enable edge scroll" : "disable edge scroll");
    }

    XCloseDevice (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), xdev);
}

void
set_two_finger_scroll(int enable_vert, int enable_horiz)
{
    Atom prop_twofinger = XInternAtom(
                              GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                              "Synaptics Two-Finger Scrolling", FALSE);

    if (!prop_twofinger) {
        g_warning("Get Edge Prop Atom Failed");
        return;
    }

    GdkDevice *tpad = device_is_exist(TPAD_NAME_KEY);

    if (tpad == NULL) {
        g_warning("TouchPad not exist\n");
        return;
    }

    XDevice *xdev = NULL;
    xdev = open_gdk_device(tpad);

    if (xdev == NULL) {
        g_warning("Get XDevice From TouchPad");
        return;
    }

    g_print("Trying to set %s, %s\n",
            enable_vert ? "enable two finger vert scroll" : "disable two finger vert scroll",
            enable_horiz ? "enable two finger horiz scroll" : "disable two finger horiz scroll");
    gdk_error_trap_push();
    Atom act_type;
    int act_format;
    unsigned long nitems, bytes_after;
    unsigned char *data = NULL;
    int rc = XGetDeviceProperty(GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                                xdev, prop_twofinger, 0, 1, False,
                                XA_INTEGER, &act_type, &act_format, &nitems,
                                &bytes_after, &data);

    if (rc == Success && act_type == XA_INTEGER &&
            act_format == 8 && nitems >= 2) {
        data[0] = enable_vert ? 1 : 0; // set vertical(垂直) scroll
        data[1] = enable_horiz ? 1 : 0; // set horizon(水平) scroll
        XChangeDeviceProperty(GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                              xdev, prop_twofinger, XA_INTEGER, 8,
                              PropModeReplace, data, nitems);
    }

    if (rc == Success) {
        XFree(data);
    }

    if (gdk_error_trap_pop()) {
        g_warning("Error settings touchpad %s, %s\n",
                  enable_vert ? "enable two finger vert scroll" : "disable two finger vert scroll",
                  enable_horiz ? "enable two finger horiz scroll" : "disable two finger horiz scroll");
    }

    XCloseDevice (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), xdev);
}

/*
 * state: touchpad is enable
 * left_handed: true/false
 */
void
set_tab_to_click (int state, int left_handed)
{
    GdkDevice *tpad = device_is_exist(TPAD_NAME_KEY);

    if (tpad == NULL) {
        g_print("TouchPad not exist\n");
        return;
    }

    XDevice *xdev = NULL;
    xdev = open_gdk_device(tpad);

    if (xdev == NULL) {
        g_warning("Get XDevice From TouchPad");
        return;
    }

    Atom prop = XInternAtom (GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                             "Synaptics Tap Action", FALSE);

    if (!prop) {
        g_warning("Get Prop 'Synaptics Tap Action' Failed");
        return;
    }

    g_print("Settings tap to click on %s\n",
            left_handed ? "Use Left Hand" : "Use Right Hand");
    gdk_error_trap_push();
    Atom act_type;
    int act_format;
    unsigned long nitems, bytes_after;
    unsigned char *data = NULL;
    int rc = XGetDeviceProperty(GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                                xdev, prop, 0, 2, False, XA_INTEGER, &act_type,
                                &act_format, &nitems, &bytes_after, &data);

    if (rc == Success && act_type == XA_INTEGER &&
            act_format == 8 && nitems >= 7) {
        /* Set RLM mapping for 1/2/3 fingers*/
        data[4] = (state) ? ((left_handed) ? 3 : 1) : 0;
        data[5] = (state) ? ((left_handed) ? 1 : 3) : 0;
        data[6] = (state) ? 2 : 0;
        XChangeDeviceProperty (
            GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()),
            xdev, prop, XA_INTEGER, 8,
            PropModeReplace, data, nitems);
    }

    if (rc == Success) {
        XFree(data);
    }

    if (gdk_error_trap_pop()) {
        g_warning("Error settings touchpad %s\n",
                  left_handed ? "Use Left Hand" : "Use Right Hand");
    }

    XCloseDevice (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), xdev);
}

/*
 * Disable TouchPad when typing
 */
// Has been implemented in GoLang
/*void*/
/*set_disable_w_typing(int enable_w)*/
/*{*/
/*}*/

/*
 * Re-enable touchpad when any other pointing device isn't present
 * 当没有鼠标设备被设置时，重新启用触摸板
 */
// Has been implemented in dde-daemon/deepin-daemon
/*void*/
/*ensure_tpad_active()*/
/*{*/
/*}*/
