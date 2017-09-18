/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

#include <X11/Xlib.h>
#include <X11/extensions/XInput2.h>

#include "listen.h"
#include "_cgo_export.h"

static int has_xi2();
static void *listen_device_thread(void *user_data);

static Display *_disp = NULL;
static pthread_t _thrd;
static int _thrd_exit_flag = 0;

int
start_device_listener()
{
    _thrd_exit_flag = 0;

    int xi_opcode = has_xi2();
    if (xi_opcode == -1) {
        _thrd_exit_flag = 1;
        return -1;
    }

    pthread_attr_t attr;

    // Free thread resource when thread terminates
    pthread_attr_init(&attr);
    pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED);
    int ret = pthread_create(&_thrd, &attr,
                             listen_device_thread, (void*)&xi_opcode);
    pthread_attr_destroy(&attr);

    if (ret != 0 ) {
        fprintf(stderr, "Create device event listen thread failed\n");
        _thrd_exit_flag = 1;
        return -1;
    }

    pthread_join(_thrd, NULL);

    return 0;
}

void
end_device_listener()
{
    if (_disp != NULL) {
        XCloseDisplay(_disp);
    }

    if (_thrd_exit_flag != 1) {
        pthread_cancel(_thrd);
    }
}

static int
has_xi2()
{
    Display *disp = XOpenDisplay(0);
    if (!disp) {
        fprintf(stderr, "Open Display Failed in has_xi2\n");
        return -1;
    }

    int xi_opcode, event, error;
    if (!XQueryExtension(disp, "XInputExtension",
                         &xi_opcode, &event, &error)) {
        fprintf(stderr, "XInput extension not available.\n");
        XCloseDisplay(disp);
        return -1;
    }

    // We support XI 2.0
    int major = 2;
    int minor = 0;

    int rc =XIQueryVersion(disp, &major, &minor);
    if ( rc == BadRequest) {
        fprintf(stderr, "No XI2 Support.\n");
        XCloseDisplay(disp);
        return -1;
    } else if (rc != Success) {
        fprintf(stderr, "Internal Error.\n");
        XCloseDisplay(disp);
        return -1;
    }

    XCloseDisplay(disp);
    return xi_opcode;
}

static void*
listen_device_thread(void *user_data)
{
    /*int xi_opcode = *(int*)user_data;*/
    XIEventMask mask;

    mask.deviceid = XIAllDevices;
    mask.mask_len = XIMaskLen(XI_LASTEVENT);
    mask.mask = calloc(mask.mask_len, sizeof(char));

    XISetMask(mask.mask, XI_HierarchyChanged);

    _disp = XOpenDisplay(0);
    if (!_disp) {
        _thrd_exit_flag = 1;
        pthread_exit(NULL);
        return NULL;
    }

    XISelectEvents(_disp, DefaultRootWindow(_disp), &mask, 1);
    XSync(_disp, False);

    free(mask.mask);

    while (1) {
        XEvent ev;
        XGenericEventCookie *cookie = (XGenericEventCookie*)&ev.xcookie;
        XNextEvent(_disp, (XEvent*)&ev);

        if (cookie->type != GenericEvent ||
            /*cookie->extension != xi_opcode ||*/
            !XGetEventData(_disp, cookie)) {
            continue;
        }

        if (cookie->evtype == XI_HierarchyChanged) {
            XIHierarchyEvent *event = cookie->data;
            if (event->flags & XIMasterAdded ||
                event->flags & XISlaveAdded ) {
                /* int deviceid = event->info->deviceid; */
                /*printf("Device Added: %d\n", deviceid);*/
                handleDeviceChanged();
            } else if (event->flags & XIMasterRemoved ||
                       event->flags &XISlaveRemoved ) {
                /* int deviceid = event->info->deviceid; */
                /*printf("Device Removed: %d\n", deviceid);*/
                handleDeviceChanged();
            }
        }
        XFreeEventData(_disp, cookie);
    }

    _thrd_exit_flag = 1;
    XCloseDisplay(_disp);
    pthread_exit(NULL);

    return NULL;
}
