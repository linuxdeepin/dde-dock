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
#define LONG_PRESS_MAX_DISTANCE 1
#define SCREEN_WIDTH 100;
#define SCREEN_HEIGHT 100;

struct raw_multitouch_event {
    double dx_unaccel, dy_unaccel;
    double scale;
    int fingers;
    uint64_t t_start_tap;
    guint tap_id;
    bool dblclick;
};

static void raw_event_reset(struct raw_multitouch_event *event, bool reset_dblclick);

static int is_touchpad(struct libinput_device *dev);
static const char* get_multitouch_device_node(struct libinput_event *ev);

static void handle_events(struct libinput *li, struct movement *m);
static void handle_gesture_events(struct libinput_event *ev, int type);

static gboolean touch_timer_handler(gpointer data);
static void touch_timer_destroy(gpointer data);
static void start_touch_timer();
static void stop_touch_timer();
static int valid_long_press_touch(double x, double y);

static GHashTable *ev_table = NULL;
struct raw_multitouch_event *raw = NULL;
static int quit = 0;

static struct _long_press_timer {
    guint id;
    guint fingers;
    gint sent; // mousedown event sent
    double x, y;
    uint32_t width;
    uint32_t height;
} touch_timer;

static struct _short_press_timer {
    guint id;
} short_press_timer;

guint long_press_timer2 = 0;

static int long_press_duration = ALARM_TIMEOUT_DEFAULT;
static int long_press_duration2 = 1000;
static double long_press_distance = LONG_PRESS_MAX_DISTANCE;
static int short_press_duration = 200;
static int dblclick_duration = 400;

int
start_loop(int verbose, double distance)
{
    struct libinput *li = open_from_udev("seat0", NULL, verbose);
    if (!li) {
        return -1;
    }

    if (verbose) {
        g_setenv("G_MESSAGES_DEBUG", "all", TRUE);
    }
    if (distance > 0) {
        long_press_distance = distance;
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

    // movements have pointer structs inside
    struct movement movements[MOV_SLOTS] = {{{0}}};

    // firstly handle all devices
    handle_events(li, movements);

    struct pollfd fds;
    fds.fd = libinput_get_fd(li);
    fds.events = POLLIN;
    fds.revents = 0;

    quit = 0;
    while(!quit && poll(&fds, 1, -1) > -1) {
        handle_events(li, movements);
        handle_movements(movements);        //handle touch screen
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
    g_debug("[Duration] set: %d --> %d", long_press_duration, duration);
    if (duration == long_press_duration) {
        return;
    }
    long_press_duration = duration;
}


void
set_timer_short_duration(int duration)
{
    g_debug("[Duration short ] set: %d --> %d", short_press_duration, duration);
    if (duration == short_press_duration) {
        return ;
    }
    short_press_duration = duration;
}

void set_dblclick_duration(int duration) 
{
    if (duration == dblclick_duration) {
        return ;
    }
    dblclick_duration = duration;
}

static void
raw_event_reset(struct raw_multitouch_event *event, bool reset_dblclick)
{
    if (!event) {
        return ;
    }

    event->dx_unaccel = 0.0;
    event->dy_unaccel = 0.0;
    event->scale = 0.0;
    event->fingers = 0;
    event->t_start_tap = 0;
    event->tap_id = 0;
    if (reset_dblclick)
        event->dblclick = false;
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

static gboolean
short_press_timer_handler(gpointer data)
{
    point scale = get_last_point_scale();
    handleTouchShortPress(short_press_duration, scale.x, scale.y);
    return FALSE;
}

static gboolean
long_press_timer_handler2(gpointer data)
{
    point scale = get_last_point_scale();
    handleTouchPressTimeout(1, long_press_duration2, scale.x, scale.y);
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
short_press_timer_destory(gpointer data)
{
    short_press_timer.id = 0;
}

static void
long_press_timer_destroy2(gpointer data)
{
    long_press_timer2 = 0;
}

static void
start_touch_timer()
{
    g_debug("start touch timer: %u, fingers: %d, long_press_duration: %d, short_press_duration: %d", touch_timer.id, touch_timer.fingers, long_press_duration, short_press_duration);
    if (touch_timer.id != 0) {
        g_debug("There has an touch_timer: %u", touch_timer.id);
        return;
    }
    touch_timer.id = g_timeout_add_full(G_PRIORITY_DEFAULT, long_press_duration, touch_timer_handler, NULL, touch_timer_destroy);
    short_press_timer.id = g_timeout_add_full(G_PRIORITY_DEFAULT - 1, short_press_duration, short_press_timer_handler, NULL, short_press_timer_destory);
    long_press_timer2 = g_timeout_add_full(G_PRIORITY_DEFAULT - 2, long_press_duration2, long_press_timer_handler2, NULL, long_press_timer_destroy2);
}

static void
stop_touch_timer()
{
    g_debug("stop touch timer: %u, fingers: %d", touch_timer.id, touch_timer.fingers);

    if (touch_timer.id != 0) {
         g_source_remove(touch_timer.id);
         touch_timer.id = 0;
    }
    if (short_press_timer.id != 0) {
        g_source_remove(short_press_timer.id);
        short_press_timer.id = 0;
    }
    if (long_press_timer2 != 0) {
        g_source_remove(long_press_timer2);
        long_press_timer2 = 0;
    }
}

static int
valid_long_press_touch(double x, double y)
{
    if (touch_timer.id == 0) {
        return 0;
    }

    double dx = x - touch_timer.x;
    double dy = y - touch_timer.y;
    return fabs(dx) < long_press_distance && fabs(dy) < long_press_distance;
}

static gboolean
 handle_tap()
{
    g_debug("[Tap] fingers: %d", raw->fingers);
    handleGestureEvent(GESTURE_TYPE_TAP, GESTURE_DIRECTION_NONE, raw->fingers);
    return FALSE;
}

static void
handle_tap_destroy()
{
    if (raw && raw->tap_id) {
        raw->tap_id = 0;
    }
}

static void
handle_tap_stop()
{
    if (raw && raw->tap_id) {
         g_source_remove(raw->tap_id);
         raw->tap_id = 0;
    }
}

static int
handle_tap_delay()
{
    raw->tap_id = g_timeout_add_full(G_PRIORITY_DEFAULT, dblclick_duration,
                                                       handle_tap, NULL, handle_tap_destroy);
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
    raw = g_hash_table_lookup(ev_table, node);
    if (!raw) {
        fprintf(stderr, "Not found '%s' in table\n", node);
        return ;
    }
    struct libinput_event_gesture *gesture = libinput_event_get_gesture_event(ev);
    if (raw->dblclick
    && type != LIBINPUT_EVENT_GESTURE_SWIPE_BEGIN
    && type != LIBINPUT_EVENT_GESTURE_SWIPE_UPDATE
    && type != LIBINPUT_EVENT_GESTURE_SWIPE_END
    && type != LIBINPUT_EVENT_GESTURE_TAP_UPDATE
    && type != LIBINPUT_EVENT_GESTURE_TAP_END) {
        raw->fingers = libinput_event_gesture_get_finger_count(gesture);
        handleSwipeStop(raw->fingers);
        raw->dblclick = false;
    }

    switch (type) {
    case LIBINPUT_EVENT_GESTURE_PINCH_BEGIN:
    case LIBINPUT_EVENT_GESTURE_SWIPE_BEGIN:
        // reset
        raw_event_reset(raw, false);
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

        int fingers = libinput_event_gesture_get_finger_count(gesture);
        if (raw->dblclick) {
            handleSwipeMoving(fingers, dx_unaccel, dy_unaccel);
        }

        break;
    }
    case LIBINPUT_EVENT_GESTURE_PINCH_END:{
        // filter small scale threshold
        if (fabs(raw->scale) < 1) {
            raw_event_reset(raw, true);
            break;
        }

        raw->fingers = libinput_event_gesture_get_finger_count(gesture);
        g_debug("[Pinch] direction: %s, fingers: %d",
                raw->scale>= 0?"in":"out", raw->fingers);
        handleGestureEvent(GESTURE_TYPE_PINCH,
                           (raw->scale >= 0?GESTURE_DIRECTION_IN:GESTURE_DIRECTION_OUT),
                           raw->fingers);
        raw_event_reset(raw, true);
        break;
    }
    case LIBINPUT_EVENT_GESTURE_SWIPE_END:
        raw->fingers = libinput_event_gesture_get_finger_count(gesture);
        if (raw->dblclick) {
            handleSwipeStop(raw->fingers);
            raw_event_reset(raw, true);
            break;
        }
        // filter small movement threshold
        if (fabs(raw->dx_unaccel - raw->dy_unaccel) < 70) {
            raw_event_reset(raw, true);
            break;
        }

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

        raw_event_reset(raw, true);
        break;
    case LIBINPUT_EVENT_GESTURE_TAP_BEGIN:
        g_debug("[Tap begin] time: %u duration: %d fingers: %d \n", raw->t_start_tap, (libinput_event_gesture_get_time_usec(gesture) - raw->t_start_tap) / 1000, raw->fingers);
        if (raw->t_start_tap > 0
        &&  (libinput_event_gesture_get_time_usec(gesture) - raw->t_start_tap) / 1000 <= dblclick_duration
        && raw->fingers == libinput_event_gesture_get_finger_count(gesture)) {
            handleDbclickDown(raw->fingers);
            handle_tap_stop();
            raw_event_reset(raw, true);
            raw->dblclick = true;
        }
        break;
    case LIBINPUT_EVENT_GESTURE_TAP_END:
        if (libinput_event_gesture_get_cancelled(gesture)) {
            raw_event_reset(raw, true);
            break;
        }

        if (!raw->dblclick) {
            raw->fingers = libinput_event_gesture_get_finger_count(gesture);
            raw->t_start_tap = libinput_event_gesture_get_time_usec(gesture);
            handle_tap_delay();
        } else {
            raw_event_reset(raw, true);
        }
        break;
    }
}

static void
handle_touch_events(struct libinput_event *ev, int ty,struct movement *m)
{
    point scale;
    struct libinput_device *dev = libinput_event_get_device(ev);
    if (!dev) {
        fprintf(stderr, "Get device from event failure\n");
        return ;
    }

    switch (ty) {
    case LIBINPUT_EVENT_TOUCH_MOTION:{
        handle_touch_event_motion(ev, m);
        g_debug("Touch motion, id: %u, fingers: %d, sent: %d ",
                touch_timer.id, touch_timer.fingers, touch_timer.sent);
        if (touch_timer.id == 0) {
            break;
        }
        struct libinput_event_touch *touch = libinput_event_get_touch_event(ev);
        // only suupurted events: down and motion
        double x = libinput_event_touch_get_x_transformed(touch, touch_timer.width);
        double y = libinput_event_touch_get_y_transformed(touch, touch_timer.height);
        g_debug("\t[Transformed] X: %lf, Y: %lf", x, y);
        if (valid_long_press_touch(x, y) == 1) {
            break;
        }
        // cancel touch_timer
        stop_touch_timer();
        break;
    }
    case LIBINPUT_EVENT_TOUCH_UP:
        handle_touch_event_up(ev, m);
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

        scale = get_last_point_scale();
        handleTouchUpOrCancel(scale.x, scale.y);
        return;
    case LIBINPUT_EVENT_TOUCH_DOWN: {
        handle_touch_event_down(ev, m);
        g_debug("Touch down, id: %u, fingers: %d, sent: %d ",
                touch_timer.id, touch_timer.fingers, touch_timer.sent);
        if (touch_timer.id != 0 || touch_timer.fingers > 0) {
            stop_touch_timer();
            touch_timer.fingers++;
            break;
        }
        struct libinput_event_touch *touch = libinput_event_get_touch_event(ev);
        start_touch_timer();
        touch_timer.fingers = 1;
        touch_timer.sent = 0;
        touch_timer.width = SCREEN_WIDTH;
        touch_timer.height = SCREEN_HEIGHT;
        double w, h;
        libinput_device_get_size(dev, &w, &h);
        // TODO(jouyouyun): save device size to cache
        if (w > 1) {
            touch_timer.width = w;
        }
        if (h >1) {
            touch_timer.height = h;
        }
        touch_timer.x = libinput_event_touch_get_x_transformed(touch, touch_timer.width);
        touch_timer.y = libinput_event_touch_get_y_transformed(touch, touch_timer.height);
        g_debug("\t[Transformed] X: %lf, Y: %lf", touch_timer.x, touch_timer.y);
        break;
    }
    case LIBINPUT_EVENT_TOUCH_FRAME:
        /* g_debug("Touch frame:"); */
        return;
    case LIBINPUT_EVENT_TOUCH_CANCEL:
        handle_touch_event_cancel(ev, m);
        g_debug("Touch cancel, id: %u, fingers: %d, sent: %d ",
                touch_timer.id, touch_timer.fingers, touch_timer.sent);
        stop_touch_timer();
        touch_timer.fingers = 0;
        if (touch_timer.sent) {
            touch_timer.sent = 0;
            handleTouchEvent(TOUCH_TYPE_RIGHT_BUTTON, BUTTON_TYPE_UP);
        }

        scale = get_last_point_scale();
        handleTouchUpOrCancel(scale.x, scale.y);
        return;
    }
}

static void
handle_events(struct libinput *li, struct movement *m)
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
        case LIBINPUT_EVENT_TOUCH_CANCEL:
        case LIBINPUT_EVENT_TOUCH_FRAME:{
            handle_touch_events(ev, type, m);  
            break;
        }
        default:
            break;
        }
        libinput_event_destroy(ev);
        libinput_dispatch(li);
    }
}
