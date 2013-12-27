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

#include "grab-xrecord.h"

typedef struct _XRecordGrabInfo {
    Display *ctrl_disp;
    Display *data_disp;
    XRecordRange *range;
    XRecordContext context;
} XRecordGrabInfo;

static void grab_key_event_cb (XPointer user_data, XRecordInterceptData *hook);
static gpointer enable_ctx_thread (gpointer user_data);
static void exec_action (int code);

static XRecordGrabInfo *grab_info = NULL;
static GHashTable *key_table = NULL;
static int key_press_cnt = 0;

void
grab_xrecord_init ()
{
    key_table = g_hash_table_new_full (g_direct_hash, g_direct_equal,
                                       NULL, (GDestroyNotify)g_free);
    grab_info = g_new0 (XRecordGrabInfo, 1);

    if ( !grab_info ) {
        g_warning ("Alloc XRecordGrabInfo memory failed...");
        grab_xrecord_finalize ();
    }

    grab_info->ctrl_disp = XOpenDisplay (NULL);
    grab_info->data_disp = XOpenDisplay (NULL);

    if ( !grab_info->ctrl_disp || !grab_info->data_disp ) {
        g_warning ("Unable to connect to X server...");
        grab_xrecord_finalize ();
    }

    gint dummy;

    if ( !XQueryExtension (grab_info->ctrl_disp, "XTEST",
                           &dummy, &dummy, &dummy) ) {
        g_warning ("XTest extension missing...");
        grab_xrecord_finalize ();
    }

    if ( !XRecordQueryVersion (grab_info->ctrl_disp, &dummy, &dummy) ) {
        g_warning ("Failed to obtain xrecord version...");
        grab_xrecord_finalize ();
    }

    grab_info->range = XRecordAllocRange ();

    if ( !grab_info->range ) {
        g_warning ("Alloc XRecordRange memory failed...");
        grab_xrecord_finalize ();
    }

    grab_info->range->device_events.first = KeyPress;
    grab_info->range->device_events.last = KeyRelease;

    XRecordClientSpec spec = XRecordAllClients;
    grab_info->context = XRecordCreateContext (
                             grab_info->data_disp, 0, &spec, 1, &grab_info->range, 1);

    if ( !grab_info->context ) {
        g_warning ("Unable to create context...");
        grab_xrecord_finalize();
    }

    XSynchronize (grab_info->ctrl_disp, TRUE);
    XFlush (grab_info->ctrl_disp);

    GThread *thrd = g_thread_new ("enable context",
                                  (GThreadFunc)enable_ctx_thread, NULL);

    if ( !thrd ) {
        g_warning ("Unable to create thread...");
        grab_xrecord_finalize ();
    }
}

void
grab_xrecord_finalize ()
{
    if (key_table) {
        g_hash_table_remove_all (key_table);
        key_table = NULL;
    }

    if (!grab_info) {
        return;
    }

    if (grab_info->context) {
        XRecordDisableContext(grab_info->data_disp, grab_info->context);
        XRecordFreeContext(grab_info->data_disp, grab_info->context);
    }

    if (grab_info->range) {
        XFree(grab_info->range);
        grab_info->range = NULL;
    }

    if (grab_info->ctrl_disp) {
        XCloseDisplay (grab_info->ctrl_disp);
        grab_info->ctrl_disp = NULL;
    }

    if (grab_info->data_disp) {
        XCloseDisplay (grab_info->data_disp);
        grab_info->data_disp = NULL;
    }

    if (grab_info) {
        g_free (grab_info);
        grab_info = NULL;
    }
}

static gpointer
enable_ctx_thread (gpointer user_data)
{
    if ( !XRecordEnableContext (grab_info->data_disp, grab_info->context,
                                grab_key_event_cb, NULL) ) {
        g_warning ("Unable to enable context...");
        grab_xrecord_finalize ();
    }

    g_thread_exit (NULL);
}

static void
grab_key_event_cb (XPointer user_data, XRecordInterceptData *hook)
{
    if ( hook->category != XRecordFromServer ) {
        g_warning ("Data not from X server...");
        return;
    }

    int event_type = hook->data[0];
    KeyCode keycode = hook->data[1];

    g_debug ("event type: %d, code: %d\n", (int)event_type, (int)keycode);

    switch (event_type) {
        case KeyPress:
            key_press_cnt++;
            break;

        case KeyRelease:
            g_debug ("key_press_cnt: %d\n", key_press_cnt);

            if (key_press_cnt == 1) {
                exec_action (keycode);
            }

            key_press_cnt = 0;
            break;

        default:
            break;
    }
}

void
grab_xrecord_key (int keycode, const char *action)
{
    if (keycode < 0 || !action) {
        g_warning ("grab key args error. keycode: %d, action:%s...",
                   keycode, action);
        return;
    }

    g_debug ("insert code: %d, action: %s\n", keycode, action);
    g_hash_table_insert (key_table,
                         GINT_TO_POINTER(keycode), g_strdup(action));
}

void
ungrab_xrecord_key (int keycode)
{
    if (keycode < 0) {
        g_warning ("ungrab key args error. keycode: %d...", keycode);
        return ;
    }

    g_hash_table_remove (key_table, GINT_TO_POINTER(keycode));
}

static void
exec_action (int code)
{
    gchar *action = g_hash_table_lookup (key_table, GINT_TO_POINTER (code));
    g_debug ("exec action: %s\n", action);

    if (action) {
        g_spawn_command_line_sync (action, NULL, NULL, NULL, NULL);
    }
}
