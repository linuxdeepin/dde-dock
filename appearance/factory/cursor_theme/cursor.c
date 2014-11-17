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

#include <gtk/gtk.h>
#include "cursor.h"

static void
update_cursor()
{
	GdkCursor *cursor = gdk_cursor_new(GDK_LEFT_PTR);
	gdk_window_set_cursor(gdk_get_default_root_window(), cursor);
	g_object_unref(G_OBJECT(cursor));
}

void
handle_cursor_changed()
{
	static gboolean init = FALSE;

	if (init) {
		return;
	}

	init = TRUE;
	GtkSettings* s = gtk_settings_get_default();
	g_signal_connect(s, "notify::gtk-cursor-theme-name",
	                 update_cursor, NULL);
}

