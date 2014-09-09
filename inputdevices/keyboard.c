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

#include <stdio.h>
#include <X11/Xlib.h>
#include <X11/XKBlib.h>

#include "devices.h"

/**
 * repear: set repeat if true
 **/
int
set_keyboard_repeat(int repeat,
                    unsigned int delay, unsigned int interval)
{
	Display *disp = XOpenDisplay(0);
	if (!disp) {
		fprintf(stderr, "Open Display Failed\n");
		return -1;
	}

	if (repeat) {
		XAutoRepeatOn(disp);

		// Use XKB in preference
		int rate_set = XkbSetAutoRepeatRate(disp, XkbUseCoreKbd,
		                                    delay, interval);
		if (!rate_set) {
			fprintf(stderr, "Neither XKeyboard not Xfree86's\
				       	keyboard extensions are available,\
					\n no way to support keyboard\
				       	autorepeat rate settings\n");
		}
	} else {
		XAutoRepeatOff(disp);
	}

	XSync(disp, False);
	XCloseDisplay(disp);

	return 0;
}
