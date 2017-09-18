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

#include <stdio.h>
#include <math.h>

#include <glib.h>
#include <poll.h>

#include "utils.h"
#include "core.h"
#include "_cgo_export.h"

struct raw_multitouch_event {
    double dx_unaccel, dy_unaccel;
    double scale;
    int fingers;
};

static void raw_event_reset(struct raw_multitouch_event *event);

static int is_touchpad(struct libinput_device *dev);
static const char* get_multitouch_device_node(struct libinput_event *ev);

static void handle_events(struct libinput *li);
static void handle_gesture_events(struct libinput_event *ev, int type);

static GHashTable *ev_table = NULL;
static int quit = 0;

int
start_loop()
{
    struct libinput *li = open_from_udev("seat0", NULL, 0);
    if (!li) {
        return -1;
    }

    ev_table = g_hash_table_new_full(g_str_hash,
                                     g_str_equal,
                                     (GDestroyNotify)g_free,
                                     (GDestroyNotify)g_free);
    if (!ev_table) {
        fprintf(stderr, "Failed to initialize event table\n");
        libinput_unref(li);
        return -1;
    }

    // firstly handle all devices
    handle_events(li);

    struct pollfd fds;
    fds.fd = libinput_get_fd(li);
    fds.events = POLLIN;
    fds.revents = 0;

    quit = 0;
    while(!quit && poll(&fds, 1, -1) > -1) {
        handle_events(li);
    }

    return 0;
}

void
quit_loop()
{
    quit = 1;
}

static void
raw_event_reset(struct raw_multitouch_event *event)
{
    if (!event) {
        return ;
    }

    event->dx_unaccel = 0.0;
    event->dy_unaccel = 0.0;
    event->scale = 0.0;
    event->fingers = 0;
}

static const char*
get_multitouch_device_node(struct libinput_event *ev)
{
    struct libinput_device *dev = libinput_event_get_device(ev);
    if (libinput_device_has_capability(dev, LIBINPUT_DEVICE_CAP_TOUCH)) {
        goto out;
    } else if (libinput_device_has_capability(dev, LIBINPUT_DEVICE_CAP_POINTER) &&
               is_touchpad(dev)) {
        goto out;
    } else {
        return NULL;
    }

out:
    return udev_device_get_devnode(libinput_device_get_udev_device(dev));
}

static int
is_touchpad(struct libinput_device *dev)
{
    // TODO: check touchpad whether support multitouch. fingers > 3?
    int cnt = libinput_device_config_tap_get_finger_count(dev);
    return (cnt > 0);
}


/**
 * calculation direction
 * Swipe: (begin -> end)
 *     _dx_unaccel += dx_unaccel, _dy_unaccel += dy_unaccel;
 *     filter small movement threshold abs(_dx_unaccel - _dy_unaccel) < 70
 *     if abs(_dx_unaccel) > abs(_dy_unaccel): _dx_unaccel < 0 ? 'left':'right'
 *     else: _dy_unaccel < 0 ? 'up':'down'
 *
 * Pinch: (begin -> end)
 *     _scale += 1.0 - scale;
 *     if _scale != 0: _scale >= 0 ? 'in':'out'
 **/
static void
handle_gesture_events(struct libinput_event *ev, int type)
{
    struct libinput_device *dev = libinput_event_get_device(ev);
    if (!dev) {
        fprintf(stderr, "Get device from event failure\n");
        return ;
    }

    const char *node = udev_device_get_devnode(libinput_device_get_udev_device(dev));
    struct raw_multitouch_event *raw = g_hash_table_lookup(ev_table, node);
    if (!raw) {
        fprintf(stderr, "Not found '%s' in table\n", node);
        return ;
    }
    struct libinput_event_gesture *gesture = libinput_event_get_gesture_event(ev);
    switch (type) {
    case LIBINPUT_EVENT_GESTURE_PINCH_BEGIN:
    case LIBINPUT_EVENT_GESTURE_SWIPE_BEGIN:
        // reset
        raw_event_reset(raw);
        break;
    case LIBINPUT_EVENT_GESTURE_PINCH_UPDATE:{
        double scale = libinput_event_gesture_get_scale(gesture);
        raw->scale += 1.0-scale;
        break;
    }
    case LIBINPUT_EVENT_GESTURE_SWIPE_UPDATE:{
        // update
        double dx_unaccel = libinput_event_gesture_get_dx_unaccelerated(gesture);
        double dy_unaccel = libinput_event_gesture_get_dy_unaccelerated(gesture);
        raw->dx_unaccel += dx_unaccel;
        raw->dy_unaccel += dy_unaccel;
        break;
    }
    case LIBINPUT_EVENT_GESTURE_PINCH_END:{
        // filter small scale threshold
        if (fabs(raw->scale) < 1) {
            raw_event_reset(raw);
            break;
        }

        raw->fingers = libinput_event_gesture_get_finger_count(gesture);
        /* printf("[Pinch] direction: %s, fingers: %d\n", */
        /*        raw->scale>= 0?"in":"out", raw->fingers); */
        handleGestureEvent(GESTURE_TYPE_PINCH,
                           (raw->scale >= 0?GESTURE_DIRECTION_IN:GESTURE_DIRECTION_OUT),
                           raw->fingers);
        raw_event_reset(raw);
        break;
    }
    case LIBINPUT_EVENT_GESTURE_SWIPE_END:
        // filter small movement threshold
        if (fabs(raw->dx_unaccel - raw->dy_unaccel) < 70) {
            raw_event_reset(raw);
            break;
        }

        raw->fingers = libinput_event_gesture_get_finger_count(gesture);
        if (fabs(raw->dx_unaccel) > fabs(raw->dy_unaccel)) {
            // right/left movement
            /* printf("[Swipe] direction: %s, fingers: %d\n", */
            /*        raw->dx_unaccel < 0?"left":"right", raw->fingers); */
            handleGestureEvent(GESTURE_TYPE_SWIPE,
                               (raw->dx_unaccel < 0?GESTURE_DIRECTION_LEFT:GESTURE_DIRECTION_RIGHT),
                               raw->fingers);
        } else {
            // up/down movement
            /* printf("[Swipe] direction: %s, fingers: %d\n", */
            /*        raw->dy_unaccel < 0?"up":"down", raw->fingers); */
            handleGestureEvent(GESTURE_TYPE_SWIPE,
                               (raw->dy_unaccel < 0?GESTURE_DIRECTION_UP:GESTURE_DIRECTION_DOWN),
                               raw->fingers);
        }
        raw_event_reset(raw);
        break;
    }
}

static void
handle_touch_events(struct libinput_event *ev, int ty)
{
    switch (ty) {
    case LIBINPUT_EVENT_TOUCH_MOTION:
        printf("Touch motion:\n");
        break;
    case LIBINPUT_EVENT_TOUCH_UP:
        printf("Touch up:\n");
        return;
    case LIBINPUT_EVENT_TOUCH_DOWN:
        printf("Touch down:\n");
        break;
    case LIBINPUT_EVENT_TOUCH_FRAME:
        printf("Touch frame:\n");
        return;
    case LIBINPUT_EVENT_TOUCH_CANCEL:
        printf("Touch cancel:\n");
        return;
    }

    struct libinput_event_touch *touch = libinput_event_get_touch_event(ev);
    double x = libinput_event_touch_get_x(touch);
    double y = libinput_event_touch_get_y(touch);
    printf("\tX: %lf, Y: %lf\n", x, y);
}

static void
handle_events(struct libinput *li)
{
    struct libinput_event *ev;
    libinput_dispatch(li);
    while ((ev = libinput_get_event(li))) {
        int type =libinput_event_get_type(ev);
        switch (type) {
        case LIBINPUT_EVENT_DEVICE_ADDED:{
            const char *path = get_multitouch_device_node(ev);
            if (path) {
                /* printf("Device added: %s\n", path); */
                g_hash_table_insert(ev_table, g_strdup(path),
                                    g_new0(struct raw_multitouch_event, 1));
            }
            break;
        }
        case LIBINPUT_EVENT_DEVICE_REMOVED: {
            const char *path = get_multitouch_device_node(ev);
            if (path) {
                /* printf("Will remove '%s' to table\n", path); */
                g_hash_table_remove(ev_table, path);
            }
            break;
        }
        case LIBINPUT_EVENT_GESTURE_PINCH_BEGIN:
        case LIBINPUT_EVENT_GESTURE_PINCH_UPDATE:
        case LIBINPUT_EVENT_GESTURE_PINCH_END:
        case LIBINPUT_EVENT_GESTURE_SWIPE_BEGIN:
        case LIBINPUT_EVENT_GESTURE_SWIPE_UPDATE:
        case LIBINPUT_EVENT_GESTURE_SWIPE_END:{
            handle_gesture_events(ev, type);
            break;
        }
        case LIBINPUT_EVENT_TOUCH_MOTION:
        case LIBINPUT_EVENT_TOUCH_UP:
        case LIBINPUT_EVENT_TOUCH_DOWN:
        case LIBINPUT_EVENT_TOUCH_FRAME:
        case LIBINPUT_EVENT_TOUCH_CANCEL:{
            // TODO
            /* handle_touch_events(ev, type); */
            break;
        }
        default:
            break;
        }
        libinput_event_destroy(ev);
        libinput_dispatch(li);
    }
}
