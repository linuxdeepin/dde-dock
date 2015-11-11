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
#include <locale.h>
#include <stdio.h>
#include <string.h>

#include "player.h"

static ca_context* connect_canberra_context(char* device);

static uint32_t id = 0;

int
canberra_play_system_sound(char *theme, char *event_id, char *device)
{
	int ret;
	int curid = ++id;
	setlocale(LC_ALL, "");
	ca_context *ca = connect_canberra_context(device);
	if (ca == NULL) {
		return -1;
	}

	ret = ca_context_play(ca, curid,
	                      CA_PROP_CANBERRA_XDG_THEME_NAME, theme,
	                      CA_PROP_EVENT_ID, event_id, NULL);

	// wait for end
	int playing;
	do {
		g_usleep(500000); // sleep 0.5s
		ret = ca_context_playing(ca, curid, &playing);
	} while (playing > 0);

	ca_context_destroy(ca);
	if (ret != CA_SUCCESS) {
		g_warning("play: id=%d %s\n", curid, ca_strerror(ret));
	}
	return ret;
}

int
canberra_play_sound_file(char *file, char *device)
{
	int ret;
	int curid = ++id;
	setlocale(LC_ALL, "");
	ca_context *ca = connect_canberra_context(device);
	if (ca == NULL) {
		return -1;
	}

	ret = ca_context_play(ca, curid,
	                      CA_PROP_MEDIA_FILENAME, file, NULL);

	// wait for end
	int playing;
	do {
		g_usleep(500000); // sleep 0.5s
		ret = ca_context_playing(ca, curid, &playing);
	} while (playing > 0);

	ca_context_destroy(ca);
	if (ret != CA_SUCCESS) {
		g_warning("play filename: id=%d %s\n", curid, ca_strerror(ret));
	}
	return ret;
}

static ca_context*
connect_canberra_context(char* device)
{
	ca_context* ca = NULL;
	if (ca_context_create(&ca) != 0) {
		g_warning("Create canberra context failed");
		return NULL;
	}

	// set backend driver to 'pulse'
	if (ca_context_set_driver(ca, "pulse") != 0) {
		g_warning("Set 'pulse' as backend driver failed");
		ca_context_destroy(ca);
		return NULL;
	}

	if ((device != NULL) && (strlen(device) > 0)) {
		if (ca_context_change_device(ca, device) != 0) {
			g_warning("Set '%s' as backend device failed");
			ca_context_destroy(ca);
			return NULL;
		}
	}

	if (ca_context_open(ca) != 0) {
		g_warning("Connect the context to sound system failed");
		ca_context_destroy(ca);
		return NULL;
	}

	return ca;
}
