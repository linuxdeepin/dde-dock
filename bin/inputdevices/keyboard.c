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
#include <gdk/gdkx.h>
#include <X11/XKBlib.h>
/*#include <X11/extensions/XKBrules.h>*/

static gboolean xkb_set_keyboard_autorepeat_rate (guint delay,
        guint interval);

/*
 * repeat: set repeat if true
 */
void
set_keyboard_repeat(int repeat, unsigned int interval, unsigned int delay)
{
    gdk_error_trap_push ();

    if (repeat) {
        gboolean rate_set = FALSE;
        XAutoRepeatOn(GDK_DISPLAY_XDISPLAY(gdk_display_get_default()));
        /* Use XKB in preference */
        rate_set = xkb_set_keyboard_autorepeat_rate(interval, delay);

        if (!rate_set) {
            g_warning("Neither XKeyboard not Xfree86's keyboard extensions \
                    are available,\nno way to support keyboard \
                    autorepeat rate settings");
        }
    } else {
        XAutoRepeatOff(GDK_DISPLAY_XDISPLAY(gdk_display_get_default()));
    }

    XSync (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()), FALSE);
    gdk_error_trap_pop_ignored ();
}

static gboolean
xkb_set_keyboard_autorepeat_rate (guint delay, guint interval)
{
    return XkbSetAutoRepeatRate (GDK_DISPLAY_XDISPLAY (gdk_display_get_default ()),
                                 XkbUseCoreKbd,
                                 delay,
                                 interval);
}
