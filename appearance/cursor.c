/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <gtk/gtk.h>

#include "cursor.h"

static void update_gtk_cursor();

void
handle_gtk_cursor_changed()
{
    gtk_init(NULL, NULL);
    GtkSettings* s = gtk_settings_get_default();
    int sig_id = g_signal_connect(s, "notify::gtk-cursor-theme-name",
                                  update_gtk_cursor, NULL);
    if (sig_id <= 0) {
        return;
    }

    gtk_main();
}

static void
update_gtk_cursor()
{
    GdkCursor* cursor = gdk_cursor_new_for_display(
        gdk_display_get_default(),
        GDK_LEFT_PTR);
    gdk_window_set_cursor(gdk_get_default_root_window(), cursor);
    g_object_unref(G_OBJECT(cursor));
}
