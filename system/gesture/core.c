/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

#define ALARM_TIMEOUT_DEFAULT 700 // 700ms
#define LONG_PRESS_MAX_DISTANCE 5

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

static gboolean touch_timer_handler(gpointer data);
static void touch_timer_destroy(gpointer data);
static void start_touch_timer();
static void stop_touch_timer();
static int valid_long_press_touch(double x, double y);

static GHashTable *ev_table = NULL;
static int quit = 0;

static struct _long_press_timer {
    guint id;
    guint fingers;
    gint sent; // mousedown event sent
    double x, y;
} touch_timer;
static int longpress_duration = ALARM_TIMEOUT_DEFAULT;

int
start_loop(void)
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
quit_loop(void)
{
    quit = 1;
}

void
set_timer_duration(int duration)
{
    g_debug("[Duration] set: %d --> %d", longpress_duration, duration);
    if (duration == longpress_duration) {
        return ;
    }
    longpress_duration = duration;
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

static gboolean
touch_timer_handler(gpointer data)
{
    g_debug("touch timer arrived: %u, data: (%f. %f)", touch_timer.id,
			touch_timer.x, touch_timer.y);
    handleTouchEvent(TOUCH_TYPE_RIGHT_BUTTON, BUTTON_TYPE_DOWN);
    touch_timer.sent = 1;
    return FALSE;
}

static void
touch_timer_destroy(gpointer data)
{
    g_debug("touch timer destroy: %u, data: (%f, %f)", touch_timer.id,
			touch_timer.x, touch_timer.y);
    if (touch_timer.id != 0) {
		touch_timer.id = 0;
    }
    touch_timer.x = 0;
    touch_timer.y = 0;
}

static void
start_touch_timer()
{
    g_debug("stop touch timer: %u, fingers: %d", touch_timer.id, touch_timer.fingers);
    if (touch_timer.id != 0) {
		g_debug("There has an touch_timer: %u", touch_timer.id);
		return;
    }
    touch_timer.id = g_timeout_add_full(G_PRIORITY_DEFAULT, longpress_duration,
										touch_timer_handler, NULL, touch_timer_destroy);
}

static void
stop_touch_timer()
{
    g_debug("stop touch timer: %u, fingers: %d", touch_timer.id, touch_timer.fingers);
    if (touch_timer.id == 0) {
		return ;
    }
    g_source_remove(touch_timer.id);
    touch_timer.id = 0;
}

static int
valid_long_press_touch(double x, double y)
{
    if (touch_timer.id == 0) {
		return 0;
    }

    double dx = x - touch_timer.x;
    double dy = y - touch_timer.y;
    return fabs(dx) < LONG_PRESS_MAX_DISTANCE && fabs(dy) < LONG_PRESS_MAX_DISTANCE;
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
        g_debug("[Pinch] direction: %s, fingers: %d",
                raw->scale>= 0?"in":"out", raw->fingers);
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
            g_debug("[Swipe] direction: %s, fingers: %d",
                    raw->dx_unaccel < 0?"left":"right", raw->fingers);
            handleGestureEvent(GESTURE_TYPE_SWIPE,
                               (raw->dx_unaccel < 0?GESTURE_DIRECTION_LEFT:GESTURE_DIRECTION_RIGHT),
                               raw->fingers);
        } else {
            // up/down movement
            g_debug("[Swipe] direction: %s, fingers: %d",
                    raw->dy_unaccel < 0?"up":"down", raw->fingers);
            handleGestureEvent(GESTURE_TYPE_SWIPE,
                               (raw->dy_unaccel < 0?GESTURE_DIRECTION_UP:GESTURE_DIRECTION_DOWN),
                               raw->fingers);
        }
        raw_event_reset(raw);
        break;
    case LIBINPUT_EVENT_GESTURE_TAP_BEGIN:
        break;
    case LIBINPUT_EVENT_GESTURE_TAP_END:
        if (libinput_event_gesture_get_cancelled(gesture)) {
			break;
        }
        raw->fingers = libinput_event_gesture_get_finger_count(gesture);
        g_debug("[Tap] fingers: %d", raw->fingers);
        handleGestureEvent(GESTURE_TYPE_TAP, GESTURE_DIRECTION_NONE, raw->fingers);
        break;
    }
}

static void
handle_touch_events(struct libinput_event *ev, int ty)
{
    struct libinput_device *dev = libinput_event_get_device(ev);
    if (!dev) {
        fprintf(stderr, "Get device from event failure\n");
        return ;
    }

    switch (ty) {
    case LIBINPUT_EVENT_TOUCH_MOTION:{
		g_debug("Touch motion, id: %u, fingers: %d, sent: %d ",
				touch_timer.id, touch_timer.fingers, touch_timer.sent);
		if (touch_timer.id == 0) {
			break;
		}
		struct libinput_event_touch *touch = libinput_event_get_touch_event(ev);
		// only suupurted events: down and motion
		double x = libinput_event_touch_get_x(touch);
		double y = libinput_event_touch_get_y(touch);
		g_debug("\tX: %lf, Y: %lf", x, y);
		if (valid_long_press_touch(x, y) == 1) {
			break;
		}
		// cancel touch_timer
		stop_touch_timer();
		break;
    }
    case LIBINPUT_EVENT_TOUCH_UP:
        g_debug("Touch up, id: %u, fingers: %d, sent: %d ",
				touch_timer.id, touch_timer.fingers, touch_timer.sent);
		stop_touch_timer();
		if (touch_timer.fingers > 0) {
			touch_timer.fingers--;
		}
		if (touch_timer.sent) {
			touch_timer.sent = 0;
			handleTouchEvent(TOUCH_TYPE_RIGHT_BUTTON, BUTTON_TYPE_UP);
		}
		return;
    case LIBINPUT_EVENT_TOUCH_DOWN: {
        g_debug("Touch down, id: %u, fingers: %d, sent: %d ",
				touch_timer.id, touch_timer.fingers, touch_timer.sent);
		if (touch_timer.id != 0 || touch_timer.fingers > 0) {
			stop_touch_timer();
			touch_timer.fingers++;
			break;
		}
		struct libinput_event_touch *touch = libinput_event_get_touch_event(ev);
		touch_timer.x = libinput_event_touch_get_x(touch);
		touch_timer.y = libinput_event_touch_get_y(touch);
		g_debug("\tX: %lf, Y: %lf", touch_timer.x, touch_timer.y);
		start_touch_timer();
		touch_timer.fingers = 1;
		touch_timer.sent = 0;
		break;
    }
    case LIBINPUT_EVENT_TOUCH_FRAME:
        /* g_debug("Touch frame:"); */
        return;
    case LIBINPUT_EVENT_TOUCH_CANCEL:
        g_debug("Touch cancel, id: %u, fingers: %d, sent: %d ",
				touch_timer.id, touch_timer.fingers, touch_timer.sent);
		stop_touch_timer();
		touch_timer.fingers = 0;
		if (touch_timer.sent) {
			touch_timer.sent = 0;
			handleTouchEvent(TOUCH_TYPE_RIGHT_BUTTON, BUTTON_TYPE_UP);
		}
		return;
    }
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
                /* g_debug("Device added: %s", path); */
                g_hash_table_insert(ev_table, g_strdup(path),
                                    g_new0(struct raw_multitouch_event, 1));
            }
            break;
        }
        case LIBINPUT_EVENT_DEVICE_REMOVED: {
            const char *path = get_multitouch_device_node(ev);
            if (path) {
                /* g_debug("Will remove '%s' to table", path); */
                g_hash_table_remove(ev_table, path);
            }
            break;
        }
        case LIBINPUT_EVENT_GESTURE_PINCH_BEGIN:
        case LIBINPUT_EVENT_GESTURE_PINCH_UPDATE:
        case LIBINPUT_EVENT_GESTURE_PINCH_END:
        case LIBINPUT_EVENT_GESTURE_SWIPE_BEGIN:
        case LIBINPUT_EVENT_GESTURE_SWIPE_UPDATE:
        case LIBINPUT_EVENT_GESTURE_SWIPE_END:
        case LIBINPUT_EVENT_GESTURE_TAP_BEGIN:
        case LIBINPUT_EVENT_GESTURE_TAP_UPDATE:
        case LIBINPUT_EVENT_GESTURE_TAP_END:{
            handle_gesture_events(ev, type);
            break;
        }
        case LIBINPUT_EVENT_TOUCH_MOTION:
        case LIBINPUT_EVENT_TOUCH_UP:
        case LIBINPUT_EVENT_TOUCH_DOWN:
        case LIBINPUT_EVENT_TOUCH_FRAME:
        case LIBINPUT_EVENT_TOUCH_CANCEL:{
            handle_touch_events(ev, type);
            break;
        }
        default:
            break;
        }
        libinput_event_destroy(ev);
        libinput_dispatch(li);
    }
}
