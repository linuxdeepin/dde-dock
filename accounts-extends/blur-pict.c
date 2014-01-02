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

#include <time.h>
#include <stdlib.h>
#include <string.h>
#include <glib.h>
#include <glib/gstdio.h>
#include <gdk-pixbuf/gdk-pixbuf.h>
#include "gaussianiir2d.h"
#include "blur-pict.h"

#define	BG_BLUR_PICT_CACHE_DIR	"gaussian-background"
#define BG_EXT_URI	 "tEXt::Blur::URI"
#define BG_EXT_MTIME	 "tEXt::Blur::MTime"

static time_t get_file_mtime (const char *file);

int
generate_blur_pict (const char *src_path, const char *dest_path)
{
    GError *error = NULL;
    GdkPixbuf *pixbuf = gdk_pixbuf_new_from_file (src_path, &error);

    if ( error ) {
        g_debug ("new pixbuf failed: %s", error->message);
        g_error_free (error);
        return FALSE;
    }

    int width = gdk_pixbuf_get_width (pixbuf);
    int height = gdk_pixbuf_get_height (pixbuf);
    int rowstride = gdk_pixbuf_get_rowstride (pixbuf);
    int n_channels = gdk_pixbuf_get_n_channels (pixbuf);
    guchar *image_data = gdk_pixbuf_get_pixels (pixbuf);

    clock_t start = clock ();
    gaussianiir2d_pixbuf_c(image_data, width, height,
                           rowstride, n_channels, 50, 50);
    clock_t end = clock ();
    g_debug ("time : %f", (end - start) / (float)CLOCKS_PER_SEC);

    char *tmp_path;
    int tmp_fd;
    tmp_path = g_strconcat (dest_path, ".XXXXXX", NULL);
    tmp_fd = g_mkstemp (tmp_path);

    if (tmp_fd == -1) {
        g_free (tmp_path);
        return FALSE;
    }

    close (tmp_fd);

    char *src_uri_bs64;
    char mtime_str[21];
    gboolean saved_ok;
    error = NULL;

    time_t src_mtime = get_file_mtime (src_path);
    g_snprintf (mtime_str, 21, "%ld",  src_mtime);
    src_uri_bs64 = g_base64_encode ((const guchar *) src_path,
                                    strlen (src_path) + 1);
    saved_ok = gdk_pixbuf_save (pixbuf,
                                tmp_path,
                                "png", &error,
                                BG_EXT_URI, src_uri_bs64,
                                BG_EXT_MTIME, mtime_str,
                                NULL);

    g_free (src_uri_bs64);

    if (saved_ok) {
        g_chmod (tmp_path, 0664);
        g_rename (tmp_path, dest_path);
    } else {
        g_debug ("save pixbuf failed: %s", error->message);
        g_error_free (error);
        g_free (tmp_path);
        return FALSE;
    }

    g_free (tmp_path);
    return TRUE;
}

int
blur_pict_is_valid (const char *src_path, const char *dest_path)
{
    g_return_val_if_fail ((dest_path && src_path), FALSE);

    GError *error = NULL;
    GdkPixbuf *pixbuf = gdk_pixbuf_new_from_file (dest_path, &error);

    if ( error ) {
        g_debug ("new pixbuf failed: %s", error->message);
        g_error_free (error);
        return FALSE;
    }

    //1. check if the original uri matches the provided @uri
    const char *blur_uri_bs64 = gdk_pixbuf_get_option (pixbuf, BG_EXT_URI);

    if ( !blur_uri_bs64 ) {
        return FALSE;
    }

    char *src_uri_bs64 = g_base64_encode ((const guchar *)src_path,
                                          strlen (src_path) + 1);
    gboolean is_equal = (strcmp (src_uri_bs64, blur_uri_bs64) == 0);
    g_free (src_uri_bs64);

    if ( !is_equal ) {
        return FALSE;
    }

    //2. check if the modification time matches
    time_t src_mtime = get_file_mtime (src_path);
    const char *blur_mtime_str = gdk_pixbuf_get_option (pixbuf, BG_EXT_MTIME);

    if ( !blur_mtime_str ) {
        return FALSE;
    }

    time_t blur_mtime = atol (blur_mtime_str);

    if ( src_mtime != blur_mtime ) {
        return FALSE;
    }

    return TRUE;
}

static time_t
get_file_mtime (const char *file)
{
    g_return_val_if_fail (file, 0);

    time_t mtime;
    struct stat _stat_buffer;
    memset (&_stat_buffer, 0, sizeof(struct stat));

    if ( stat (file, &_stat_buffer) == 0) {
        mtime = _stat_buffer.st_mtime;
    }

    return mtime;
}
