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

#ifndef __TOUCHPAD_H__
#define __TOUCHPAD_H__

#include <glib.h>

#define TPAD_NAME_KEY "touchpad"
#define MOUSE_NAME_KEY "mouse"
#define KEYBOARD_KEY_NAME "keyboard"

int listen_device_changed ();

// TouchPad Set Func
void set_tpad_enable(int enable);
void set_natural_scroll(int enable);
void set_edge_scroll(int enable);
void set_two_finger_scroll(int enable_vert, int enable_horiz);
void set_tab_to_click (int state, int left_handed);

// Mouse Set Func
void set_motion (char *dev_name,
        double motion_acceleration, int motion_threshold);
void set_middle_button (int enable);
void set_left_handed (int left_handed);

// Keyboard Set Func
void set_keyboard_repeat(int repeat,
        unsigned int interval, unsigned int delay);

#endif

