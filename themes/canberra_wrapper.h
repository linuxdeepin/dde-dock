/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
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
#include <locale.h>
#include <stdio.h>
#include <canberra.h>

static uint32_t id = 0;

int canberra_play_system_sound(char *theme, char *event_id) {
    int ret;
    int curid = ++id;
    setlocale(LC_ALL, "");
    ca_context *co;
    ca_context_create(&co);
    ret = ca_context_play(co, curid, CA_PROP_CANBERRA_XDG_THEME_NAME, theme,
                          CA_PROP_EVENT_ID, event_id, NULL);

    // wait for end
    int playing;
    do {
        g_usleep(500000); // sleep 0.5s
        ret = ca_context_playing(co, curid, &playing);
    } while (playing > 0);

    ca_context_destroy(co);
    if (ret != CA_SUCCESS) {
        fprintf(stderr, "play: id=%d %s\n", curid, ca_strerror(ret));
    }
    return ret;
}

int canberra_play_sound_file(char *file) {
    int ret;
    int curid = ++id;
    setlocale(LC_ALL, "");
    ca_context *co;
    ca_context_create(&co);
    ret = ca_context_play(co, curid, CA_PROP_MEDIA_FILENAME, file, NULL);

    // wait for end
    int playing;
    do {
        g_usleep(500000); // sleep 0.5s
        ret = ca_context_playing(co, curid, &playing);
    } while (playing > 0);

    ca_context_destroy(co);
    if (ret != CA_SUCCESS) {
        fprintf(stderr, "play (by filename): id=%d %s\n", curid, ca_strerror(ret));
    }
    return ret;
}
