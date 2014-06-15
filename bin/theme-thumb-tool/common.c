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
#include <gtk/gtk.h>
#include "common.h"

int init_env ()
{
	return gtk_init_check(NULL, NULL);
}

int
get_base_space (int total, int dest)
{
	if (total < dest) {
		return -1;
	}

	return (total/2) - (dest/2);
}

char *
get_user_pictures_dir ()
{
	return (char*)g_get_user_special_dir (G_USER_DIRECTORY_PICTURES);
}
