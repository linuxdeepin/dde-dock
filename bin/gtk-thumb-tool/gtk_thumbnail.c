//gcc theme_preview.c `pkg-config --libs --cflags gtk+-2.0  libmetacity-private `
//
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


#include <gdk/gdk.h>
#include <gtk/gtk.h>
#include <metacity-private/common.h>
#include <metacity-private/util.h>
#include <metacity-private/boxes.h>
#include <metacity-private/gradient.h>
#include <metacity-private/theme-parser.h>
#include <metacity-private/theme.h>
#include <metacity-private/gradient.h>
#include <metacity-private/preview-widget.h>

GtkWidget *get_meta_theme (const char* theme_name)
{
	MetaTheme *meta = meta_theme_load(theme_name, NULL);
	if (!meta) {
		g_error("Get Current Meta Theme Failed");
		return NULL;
	}

	GtkWidget *preview =  meta_preview_new();
	if (!preview) {
		g_error("Get Meta Preview Failed");
		return NULL;
	}

	meta_preview_set_theme((MetaPreview*)preview, meta);
	meta_preview_set_title((MetaPreview*)preview, "Test Meta Title");

	return preview;
}

void capture(GtkOffscreenWindow* w, GdkEvent* event, gpointer user_data)
{
	gchar *dest = (gchar*)user_data;
	GdkPixbuf* pbuf = gtk_offscreen_window_get_pixbuf(w);
	gdk_pixbuf_save(pbuf, dest, "png", NULL, NULL);
	gtk_main_quit();
}

int gen_gtk_thumbnail(char *theme, char *dest)
{
	if (theme == NULL || dest == NULL) {
		g_error("gen_gtk_thumbnail args error");
		return -1;
	}
	gtk_init(NULL, NULL);

	GtkWidget* w = gtk_offscreen_window_new();
	gtk_widget_set_size_request(w, 150, 73);
	GtkWidget *t = get_meta_theme(theme);
	if (t == NULL) {
		g_error("Get Meta Theme Failed");
		return -1;
	}
	/*gtk_widget_set_size_request(t, 130, 70);*/
	gtk_container_add((GtkContainer*)w,t);

	GtkWidget *btn = gtk_button_new_with_label("Button");
	if (btn == NULL) {
		g_error("New Button Failed");
		return -1;
	}
	/*gtk_widget_set_size_request(btn, 45, 30);*/
	gtk_container_add((GtkContainer*)t, btn);

	g_signal_connect(G_OBJECT(w), "damage-event", G_CALLBACK(capture), dest);
	gtk_widget_show_all(w);

	gtk_main();

	return 0;
}
