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

#include <glib.h>
#include <cairo.h>
#include <gtk/gtk.h>
#include "common.h"

#define ICON_SIZE 48
#define PLACE_ICON_NAME "folder_home"
#define DEVICE_ICON_NAME "block-device"
#define APP_ICON_DEEPIN "deepin-software-center"
#define APP_ICON_NAME "google-chrome"

static GdkPixbuf *get_pixbuf_from_file (const char *name, gboolean is_bg);

static GdkPixbuf *
get_pixbuf_from_file (const char *name, gboolean is_bg)
{
	if (!name) {
		g_warning("Get Icon Pixbuf Failed: args error");
		return NULL;
	}

	GError *error = NULL;
	GdkPixbuf *pixbuf = NULL;
	if (!is_bg) {
		pixbuf = gdk_pixbuf_new_from_file_at_size(name, 24, 24,&error);
	} else {
		pixbuf = gdk_pixbuf_new_from_file(name, &error);
	}
	if (!pixbuf) {
		g_warning("Create %s Pixbuf Failed: %s", name, error->message);
		g_error_free(error);
		return NULL;
	}

	return pixbuf;
}

int
gen_icon_preview(char *bg, char *dest, char *item1, char *item2, char *item3)
{
	if (!bg || !dest || !item1 || !item2 || !item3) {
		g_warning("Generate Icon Preview Failed: args error");
		return -1;
	}

	GdkPixbuf *bg_pixbuf = get_pixbuf_from_file(bg, TRUE);
	if (!bg_pixbuf) {
		return -1;
	}

	int width = gdk_pixbuf_get_width(bg_pixbuf);
	int height = gdk_pixbuf_get_height(bg_pixbuf);

	cairo_surface_t *bg_surface = gdk_cairo_surface_create_from_pixbuf(
	                                  bg_pixbuf, 0, NULL);
	g_object_unref(bg_pixbuf);
	if (!bg_surface) {
		g_warning("Create Bg Surface Failed");
		return -1;
	}

	cairo_t *bg_cairo = cairo_create(bg_surface);
	if (!bg_cairo) {
		g_warning("Create Bg Cairo Failed");
		cairo_surface_destroy(bg_surface);
		return -1;
	}

	GdkPixbuf *pixbuf1 = get_pixbuf_from_file(item1, FALSE);
	GdkPixbuf *pixbuf2 = get_pixbuf_from_file(item2, FALSE);
	GdkPixbuf *pixbuf3 = get_pixbuf_from_file(item3, FALSE);
	if (!pixbuf1 || !pixbuf2 || !pixbuf3) {
		g_warning("Get Icon Pixbuf Failed");
		goto out;
	}

	int width1 = gdk_pixbuf_get_width(pixbuf1) + ICON_SPCAE;
	int height1 = gdk_pixbuf_get_height(pixbuf1);
	int width2 = gdk_pixbuf_get_width(pixbuf2) + width1 + ICON_SPCAE;
	int height2 = gdk_pixbuf_get_height(pixbuf2);
	int width3 = gdk_pixbuf_get_width(pixbuf3) + width2;
	int height3 = gdk_pixbuf_get_height(pixbuf3);
	int h = (height1 > height2)?((height1>height3)?height1:height3):((height2>height3)?height2:height3);

	int base_w = get_base_space(width, width3);
	if ( base_w == -1) {
		goto out;
	}
	int base_h = get_base_space(height, h);
	if (base_h == -1) {
		goto out;
	}

	gdk_cairo_set_source_pixbuf(bg_cairo, pixbuf1,
	                            (gdouble)(base_w), (gdouble)(base_h));
	g_object_unref(pixbuf1);
	cairo_paint(bg_cairo);
	gdk_cairo_set_source_pixbuf(bg_cairo, pixbuf2,
	                            (gdouble)(width1 + base_w), (gdouble)(base_h));
	g_object_unref(pixbuf2);
	cairo_paint(bg_cairo);
	gdk_cairo_set_source_pixbuf(bg_cairo, pixbuf3,
	                            (gdouble)(width2 +base_w), (gdouble)(base_h));
	g_object_unref(pixbuf3);
	cairo_paint(bg_cairo);

	cairo_status_t ret = cairo_surface_write_to_png(bg_surface, dest);
	g_debug("ret: %d", ret);
	gint flag = 1;

out:
	cairo_destroy(bg_cairo);
	cairo_surface_destroy(bg_surface);

	if (flag) {
		return 0;
	}

	return -1;
}
