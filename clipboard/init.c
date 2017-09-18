/*
 * Copyright (C) 2013 ~ 2017 Deepin Technology Co., Ltd.
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

#include <gtk/gtk.h>
// copy frome xfce4-clipman
#include "gsd-clipboard-manager.h"

static GsdClipboardManager* clip_manager = NULL;

int start_clip_manager()
{
    if (clip_manager) {
        return 0;
    }

    clip_manager = gsd_clipboard_manager_new();
    if (clip_manager == NULL) {
        g_warning("New Clipboard Manager Failed");
        return -1;
    }

    GError* err = NULL;
    if (!gsd_clipboard_manager_start(clip_manager, &err)) {
        g_warning("Start Clipboard Manager Failed: %s", err->message);
        g_object_unref(G_OBJECT(clip_manager));
        clip_manager = NULL;
        return -1;
    }
    /*gtk_main();*/

    return 0;
}

int stop_clip_manager()
{
    if (clip_manager != NULL) {
        gsd_clipboard_manager_stop(clip_manager);
        g_object_unref(G_OBJECT(clip_manager));
        clip_manager = NULL;
    }
    /*gtk_main_quit();*/

    return 0;
}
