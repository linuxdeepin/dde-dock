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

#include <string.h>
#include <glib.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/extensions/XTest.h>
#include "send_key_event.h"

static void send_key_event (Display *dsp, Window win,
                            int keycode, gboolean is_press);

static void
send_key_event (Display *dsp, Window win,
                int keysym, gboolean is_press)
{
    XEvent event;

    memset(&event, 0, sizeof(XEvent));
    event.xkey.display = dsp;
    event.xkey.window = win;
    event.xkey.time = CurrentTime;
    event.xkey.x = 1;
    event.xkey.y = 1;
    event.xkey.x_root = 1;
    event.xkey.y_root = 1;
    event.xkey.same_screen = True;
    event.xkey.keycode = XKeysymToKeycode(dsp, keysym);
    event.xkey.state = 0;

    if (!event.xkey.keycode) {
        return;
    }

    gint status = 0;

    if (is_press) {
        XTestFakeKeyEvent(dsp, event.xkey.keycode, True, CurrentTime);
        event.xkey.type = KeyPress;
        status = XSendEvent(dsp, win, True, 0xfff, &event);
    } else {
        XTestFakeKeyEvent(dsp, event.xkey.keycode, False, CurrentTime);
        event.xkey.type = KeyRelease;
        status = XSendEvent(dsp, win, True, 0xfff, &event);
    }

    if (status == 0 ) {
        g_print("Send Key Event Failed!\n");
        return ;
    }

    XFlush(dsp);
}

// shortcut: <Super>+W
void
initate_windows()
{
    Display *dsp = NULL;

    if (!(dsp = XOpenDisplay(NULL))) {
        g_print("Cannot open display!\n");
        return;
    }

    Window win = DefaultRootWindow(dsp);

    send_key_event(dsp, win, XK_Super_L, True);
    send_key_event(dsp, win, XK_W, True);
    send_key_event(dsp, win, XK_Super_L, False);
    send_key_event(dsp, win, XK_W, False);
    /*send_key_event(dsp, win, XK_BackSpace, True);*/
    /*send_key_event(dsp, win, XK_BackSpace, False);*/

    XCloseDisplay(dsp);
}

/*
int
main ()
{
    initate_windows();

    return 0;
}
*/
