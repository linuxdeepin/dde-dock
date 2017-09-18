/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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

#include <gio/gio.h>

#define FILE_MANAGER_MIME_TYPE "inode/directory"

/**
 * Open default file manager via mime type
 * Some file managers custom the directory by user
 **/
int
main(int argc, char *argv[])
{
    GAppInfo *app_info = g_app_info_get_default_for_type(FILE_MANAGER_MIME_TYPE, FALSE);
    if (!app_info) {
        g_error("Failed to get default app for %s", FILE_MANAGER_MIME_TYPE);
        return -1;
    }
    g_debug("Executable: %s\n", g_app_info_get_executable(app_info));
    g_debug("Commandline: %s\n", g_app_info_get_commandline(app_info));

    GError *error = NULL;
    gboolean ret = g_app_info_launch(app_info, NULL, NULL, &error);
    if (error) {
        g_error("Failed to launch %s, error: %s", g_app_info_get_name(app_info), error->message);
        g_error_free(error);
        goto EXIT;
    }

EXIT:
    g_object_unref(app_info);
    return ret?0:-1;
}
