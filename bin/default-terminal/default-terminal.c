/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

static char*
get_terminal_path(const char* preset)
{
    if (preset != NULL) {
        char *path = g_find_program_in_path(preset);
        if (path) {
            return path;
        }
    }

    char *list[] = {
        "deepin-terminal",
        "gnome-terminal",
        "terminator",
        "xfce4-terminal",
        "rxvt",
        NULL
    };

    int i = 0;
    char *tmp = list[i];
    while (tmp != NULL) {
        if (g_strcmp0(tmp, preset) != 0) {
            char *path = g_find_program_in_path(tmp);
            if (path) {
                return path;
            }
        }
	tmp = list[++i];
    }
    // default
    return g_strdup("xterm");
}

int main(int argc, char* *argv)
{
    GSettings *s = g_settings_new("com.deepin.desktop.default-applications.terminal");
    char* exec = g_settings_get_string(s, "exec");
    char* exec_arg = g_settings_get_string(s, "exec-arg");
    g_object_unref(s);
    if (g_strcmp0(exec, argv[0]) == 0) {
        goto out;
    }

    argv[0] = exec;
    if (argc > 1 && (g_strcmp0(exec, "deepin-terminal") == 0 ||
                     g_strcmp0(exec, "gnome-terminal") == 0 ||
                     g_strcmp0(exec, "xfce4-terminal") == 0 ||
                     g_strcmp0(exec, "terminator") == 0)) {
        argv[1] = exec_arg;
    }
    char* app = get_terminal_path(exec);
    int pid = fork();
    if (pid == 0) {
        execv(app, argv);
    } else {
        /* wait(NULL); */
    }
    g_free(app);

out:
    g_free(exec);
    g_free(exec_arg);
    return 0;
}
