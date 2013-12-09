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

#include "is-laptop.h"

int
IsLaptop (void)
{
    GError *err = NULL;
    GdkScreen *screen = gdk_screen_get_default ();
    GnomeRRScreen *rr_screen = gnome_rr_screen_new (screen, &err);

    if ( err ) {
        g_warning ("get gnome rr screen failed: %s\n", err->message);
        g_error_free (err);
        return -1;
    }

    GnomeRRConfig *result = gnome_rr_config_new_current (rr_screen, NULL);
    GnomeRROutputInfo **outputs = gnome_rr_config_get_outputs (result);
    gnome_rr_config_set_clone (result, FALSE);

    int i = 0;

    for (; outputs[i] != NULL; ++i) {
        GnomeRROutputInfo *info = outputs[i];
        GnomeRROutput *rr_output = gnome_rr_screen_get_output_by_name (
                                       rr_screen,
                                       gnome_rr_output_info_get_name (info));

        if (gnome_rr_output_is_laptop (rr_output)) {
            return 1;
        }
    }

    return -1;
}
