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

#include <math.h>
#include <gdk/gdk.h>
#include "utils.h"
#include "devices.h"

static void configure_button_layout (guchar   *buttons,
                                     gint      n_buttons,
                                     gboolean  left_handed);

void
set_motion (char *dev_name, double motion_acceleration, int motion_threshold)
{
    GdkDevice *device = device_is_exist(dev_name);

    if (device == NULL) {
        g_warning("%s not exist\n", dev_name);
        return;
    }

    XDevice *xdev = NULL;
    xdev = open_gdk_device(device);

    if (xdev == NULL) {
        g_warning("Get XDevice From %s Failed\n", dev_name);
        return;
    }

    XPtrFeedbackControl feedback;
    XFeedbackState *states, *state;
    int num_feedbacks;
    int numerator, denominator;

    if (motion_acceleration >= 1.0) {
        /* we want to get the acceleration, with a resolution of 0.5
         */
        if ((motion_acceleration - floor (motion_acceleration)) < 0.25) {
            numerator = floor (motion_acceleration);
            denominator = 1;
        } else if ((motion_acceleration - floor (motion_acceleration)) < 0.5) {
            numerator = ceil (2.0 * motion_acceleration);
            denominator = 2;
        } else if ((motion_acceleration - floor (motion_acceleration)) < 0.75) {
            numerator = floor (2.0 * motion_acceleration);
            denominator = 2;
        } else {
            numerator = ceil (motion_acceleration);
            denominator = 1;
        }
    } else if (motion_acceleration < 1.0 && motion_acceleration > 0) {
        /* This we do to 1/10ths */
        numerator = floor (motion_acceleration * 10) + 1;
        denominator = 10;
    } else {
        numerator = -1;
        denominator = -1;
    }

    /* Get the list of feedbacks for the device */
    states = XGetFeedbackControl (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), xdev, &num_feedbacks);

    if (states == NULL) {
        goto out;
    }

    state = (XFeedbackState *) states;

    guint i;

    for (i = 0; i < num_feedbacks; i++) {
        if (state->class == PtrFeedbackClass) {
            /* And tell the device */
            feedback.class      = PtrFeedbackClass;
            feedback.length     = sizeof (XPtrFeedbackControl);
            feedback.id         = state->id;
            feedback.threshold  = motion_threshold;
            feedback.accelNum   = numerator;
            feedback.accelDenom = denominator;

            g_debug ("Setting accel %d/%d, threshold %d for device '%s'",
                     numerator, denominator, motion_threshold, gdk_device_get_name (device));

            XChangeFeedbackControl (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()),
                                    xdev,
                                    DvAccelNum | DvAccelDenom | DvThreshold,
                                    (XFeedbackControl *) &feedback);

            break;
        }

        state = (XFeedbackState *) ((char *) state + state->length);
    }

    XFreeFeedbackList (states);

out:

    XCloseDevice (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), xdev);
}

void
set_middle_button (int enable)
{
    Atom prop = XInternAtom (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()),
                             "Evdev Middle Button Emulation", True);

    if (!prop) { /* no evdev devices */
        return;
    }


    GdkDevice *mouse = device_is_exist(MOUSE_NAME_KEY);

    if (mouse == NULL) {
        g_warning("Mouse device not exist\n");
        return;
    }

    XDevice *xdev = NULL;
    xdev = open_gdk_device(mouse);

    if (xdev == NULL) {
        g_warning("Get XDevice From TouchPad");
        return;
    }

    g_print("Set middle button on %s\n", gdk_device_get_name(mouse));
    gdk_error_trap_push();
    Atom act_type;
    int act_format;
    unsigned long nitems, bytes_after;
    unsigned char *data = NULL;
    int rc = XGetDeviceProperty(GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                                xdev, prop, 0, 1, False, XA_INTEGER, &act_type,
                                &act_format, &nitems, &bytes_after, &data);

    if (rc == Success && act_type == XA_INTEGER &&
            act_format == 8 && nitems == 1) {
        data[0] = enable ? 1 : 0;
        XChangeDeviceProperty(GDK_DISPLAY_XDISPLAY(gdk_display_get_default()),
                              xdev, prop, act_type, act_format, PropModeReplace,
                              data, nitems);
    }

    if (rc == Success) {
        XFree(data);
    }

    if (gdk_error_trap_pop()) {
        g_warning ("Error in setting middle button emulation on \"%s\"",
                   gdk_device_get_name (mouse));
    }

    XCloseDevice (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), xdev);
}

void
set_left_handed (int left_handed)
{
    GdkDevice *mouse = device_is_exist(MOUSE_NAME_KEY);

    if (mouse == NULL) {
        g_print("Mouse not exist\n");
        return;
    }

    XDevice *xdev = NULL;
    xdev = open_gdk_device(mouse);

    if (xdev == NULL) {
        g_warning("Get XDevice From Mouse");
        return;
    }

    g_print ("setting handedness on %s\n", gdk_device_get_name (mouse));

    gsize buttons_capacity = 16;
    guchar *buttons = g_new0 (guchar, buttons_capacity);

    gint n_buttons = XGetDeviceButtonMapping (
                         GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), xdev,
                         buttons, buttons_capacity);

    while (n_buttons > buttons_capacity) {
        buttons_capacity = n_buttons;
        buttons = (guchar *) g_realloc (buttons,
                                        buttons_capacity * sizeof (guchar));

        n_buttons = XGetDeviceButtonMapping (
                        GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), xdev,
                        buttons, buttons_capacity);
    }

    configure_button_layout (buttons, n_buttons, left_handed);

    gdk_error_trap_push ();
    XSetDeviceButtonMapping (
        GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), xdev,
        buttons, n_buttons);
    gdk_error_trap_pop_ignored ();

    XCloseDevice (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), xdev);
    g_free (buttons);
}

static void
configure_button_layout (guchar    *buttons,
                         gint      n_buttons,
                         gboolean  left_handed)
{
    const gint left_button = 1;
    gint right_button;
    gint i;

    /* if the button is higher than 2 (3rd button) then it's
     * probably one direction of a scroll wheel or something else
     * uninteresting
     */
    right_button = MIN (n_buttons, 3);

    /* If we change things we need to make sure we only swap buttons.
     * If we end up with multiple physical buttons assigned to the same
     * logical button the server will complain. This code assumes physical
     * button 0 is the physical left mouse button, and that the physical
     * button other than 0 currently assigned left_button or right_button
     * is the physical right mouse button.
     */

    /* check if the current mapping satisfies the above assumptions */
    if (buttons[left_button - 1] != left_button &&
            buttons[left_button - 1] != right_button)
        /* The current mapping is weird. Swapping buttons is probably not a
         * good idea.
         */
    {
        return;
    }

    /* check if we are left_handed and currently not swapped */
    if (left_handed && buttons[left_button - 1] == left_button) {
        /* find the right button */
        for (i = 0; i < n_buttons; i++) {
            if (buttons[i] == right_button) {
                buttons[i] = left_button;
                break;
            }
        }

        /* swap the buttons */
        buttons[left_button - 1] = right_button;
    }
    /* check if we are not left_handed but are swapped */
    else if (!left_handed && buttons[left_button - 1] == right_button) {
        /* find the right button */
        for (i = 0; i < n_buttons; i++) {
            if (buttons[i] == left_button) {
                buttons[i] = right_button;
                break;
            }
        }

        /* swap the buttons */
        buttons[left_button - 1] = left_button;
    }
}
