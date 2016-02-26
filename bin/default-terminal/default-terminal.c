/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <gio/gio.h>

int main(int argc, char* *argv)
{
    GSettings *s = g_settings_new("com.deepin.desktop.default-applications.terminal");
    char* exec = g_settings_get_string(s, "exec");
    g_object_unref(s);
    if (g_strcmp0(exec, argv[0]) == 0)
        return 0;
    argv[0] = exec;
    if (argc > 1 && (g_strcmp0(exec, "gnome-terminal") == 0 ||
            g_strcmp0(exec, "terminator") == 0)) {
        argv[1] = g_strdup("-x"); //Need free this?
    }
    char* app = g_find_program_in_path(exec);
    if (app == NULL) {
        app = "x-terminal-emulator";
    }
    int pid = fork();
    if (pid == 0) {
        execv(app, argv);
    } else {
        /* wait(NULL); */
    }
    g_free(app);
    return 0;
}

