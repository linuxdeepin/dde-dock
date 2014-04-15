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

#include "gsd-init.h"
#include "gsd-mouse-manager.h"
#include "gsd-keyboard-manager.h"
#include <gtk/gtk.h>
#include <glib.h>

GsdMouseManager *mouse_manager = NULL;
GsdKeyboardManager *kbd_manager = NULL;

void
gsd_init ()
{
    gtk_init(NULL, NULL);
    GError *error = NULL;

    mouse_manager = gsd_mouse_manager_new();
    kbd_manager = gsd_keyboard_manager_new ();

    if (mouse_manager == NULL || kbd_manager == NULL) {
        g_warning("New Gsd Mouse/Keyboard Failed");
        gsd_finalize();
        return;
    }

    gsd_mouse_manager_start(mouse_manager, &error);

    gsd_keyboard_manager_start (kbd_manager, &error);
    gtk_main();
}

void
gsd_finalize()
{
    if (mouse_manager != NULL) {
        gsd_mouse_manager_stop(mouse_manager);
    }

    if (kbd_manager != NULL) {
        gsd_keyboard_manager_stop(kbd_manager);
    }
}
