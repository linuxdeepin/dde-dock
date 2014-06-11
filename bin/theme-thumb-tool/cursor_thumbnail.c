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

#include <cairo.h>
#include <glib.h>
#include <gdk/gdk.h>
#include "common.h"

#define LEFT_PTR "left_ptr"
#define LEFT_PTR_WATCH "left_ptr_watch"
#define QUESTION_ARROW "question_arrow"

static GdkPixbuf *gen_pixbuf_from_cursor(const gchar *name);

static GdkPixbuf *
gen_pixbuf_from_cursor(const gchar *name)
{
	if (!name) {
		g_error("Cursor Name Is NULL");
		return NULL;
	}

	GdkDisplay *dsp = gdk_display_get_default();
	if (!dsp) {
		g_error("Get Default Display Failed");
		return NULL;
	}

	GdkCursor *cursor = gdk_cursor_new_from_name(dsp, name);
	if (!cursor) {
		g_error ("New Cursor From Name Failed");
		return NULL;
	}

	GdkPixbuf *pixbuf = gdk_cursor_get_image(cursor);
	g_object_unref(cursor);

	return pixbuf;
}

int
gen_cursor_preview(char *bg, char *dest)
{
	if (!bg || !dest) {
		g_error("Cursor Preview Args NULL");
		return -1;
	}

	GError *error = NULL;
	GdkPixbuf *bg_pixbuf = gdk_pixbuf_new_from_file(bg, &error);
	if (!bg_pixbuf) {
		g_error("Create Bg Pixbuf Failed: %s", error->message);
		g_error_free(error);
		return -1;
	}
	int width = gdk_pixbuf_get_width(bg_pixbuf);
	int height = gdk_pixbuf_get_height(bg_pixbuf);
	g_debug("Bg width: %d, height: %d", width, height);

	cairo_surface_t *bg_surface = gdk_cairo_surface_create_from_pixbuf(
			bg_pixbuf, 0, NULL);
	g_object_unref(bg_pixbuf);
	if (!bg_surface) {
		g_error("Create Bg Cairo Surface Failed");
		return -1;
	}

	cairo_t *bg_cairo = cairo_create(bg_surface);
	if (!bg_cairo) {
		g_error("Create Bg Cairo Failed");
		g_error_free(error);
		cairo_surface_destroy(bg_surface);
		return -1;
	}

	GdkPixbuf *pixbuf1 = gen_pixbuf_from_cursor(LEFT_PTR);
	GdkPixbuf *pixbuf2 = gen_pixbuf_from_cursor(LEFT_PTR_WATCH);
	GdkPixbuf *pixbuf3 = gen_pixbuf_from_cursor(QUESTION_ARROW);
	if (!pixbuf1 || !pixbuf2 || !pixbuf3) {
		g_error("Create Pixbuf From Cursor Failed");
		goto out;
	}
	int width1 = gdk_pixbuf_get_width(pixbuf1) + ICON_SPCAE;
	int height1 = gdk_pixbuf_get_height(pixbuf1);
	int width2 = gdk_pixbuf_get_width(pixbuf2) + width1 + ICON_SPCAE;
	int height2 = gdk_pixbuf_get_height(pixbuf2);
	int width3 = gdk_pixbuf_get_width(pixbuf3) + width2;
	int height3 = gdk_pixbuf_get_height(pixbuf3);
	int h = (height1 > height2)?((height1>height3)?height1:height3):((height2>height3)?height2:height3);
	g_debug("Tmp height: %d", h);

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
