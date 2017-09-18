/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

/**
 * This code was taken from this post (I did not write it):
 *
 * http://www.linuxquestions.org/questions/linux-software-2/how-to-show-desktop-in-xfce4-601161/
 *
 **/

#include <X11/Xatom.h>
#include <X11/Xlib.h>

#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main(int argc, char *argv[])
{
    Display *disp;
    Window root;
    Atom showing_desktop_atom, actual_type;
    int actual_format, error, current = 0;
    unsigned long nitems, after;
    unsigned char *data = NULL;

    /* Open the default display */
    if(!(disp = XOpenDisplay(NULL))) {
        fprintf(stderr, "Cannot open display \"%s\".\n", XDisplayName(NULL));
        return -1;
    }

    /* This is the default root window */
    root = DefaultRootWindow(disp);

    /* find the Atom for _NET_SHOWING_DESKTOP */
    showing_desktop_atom = XInternAtom(disp, "_NET_SHOWING_DESKTOP", False);

    /* Obtain the current state of _NET_SHOWING_DESKTOP on the default root window */
    error = XGetWindowProperty(disp, root, showing_desktop_atom, 0, 1, False, XA_CARDINAL,
                               &actual_type, &actual_format, &nitems, &after, &data);
    if(error != Success) {
        fprintf(stderr, "Get '_NET_SHOWING_DESKTOP' property error %d!\n", error);
        XCloseDisplay(disp);
        return -1;
    }

    /* The current state should be in data[0] */
    if(data) {
        current = data[0];
        XFree(data);
        data = NULL;
    }
    printf("Current state: %d\n", current);

    /* If nitems is 0, forget about data[0] and assume that current should be False */
    if(!nitems) {
        fprintf(stderr, "Unexpected result.\n");
        fprintf(stderr, "Assuming unshown desktop!\n");
        current = False;
    }

    /* Initialize Xevent struct */
    XEvent xev = {
        .xclient = {
            .type = ClientMessage,
            .send_event = True,
            .display = disp,
            .window = root,
            .message_type = showing_desktop_atom,
            .format = 32,
            .data.l[0] = !current /* That’s what we want the new state to be */
        }
    };

    /* Send the event to the window manager */
    XSendEvent(disp, root, False, SubstructureRedirectMask | SubstructureNotifyMask, &xev);

    /* Output the new state ("visible" or "hidden") so the calling program
     * can react accordingly.
     */
    /* printf("%s\n", current ? "hidden" : "visible"); */

    XCloseDisplay(disp);
    return 0;
}
