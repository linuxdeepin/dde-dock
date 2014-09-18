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

#ifndef __DEVICES_H__
#define __DEVICES_H__

typedef struct _DeviceInfo {
	char *name;
	int deviceid;
	int enabled;
} DeviceInfo;

DeviceInfo *get_device_info_list(int *n_devices);
void free_device_info(DeviceInfo *infos, int n_devices);

int is_mouse_device(int deviceid);
int set_motion_acceleration(int deviceid, double acceleration);
int set_motion_threshold(int deviceid, double threshold);
int set_left_handed(unsigned long xid, const char *name, int enabled);

int set_mouse_natural_scroll(unsigned long xid, const char *name, int enabled);

int is_tpad_device(int deviceid);
int set_touchpad_enabled(int deviceid, int enabled);
int set_touchpad_natural_scroll(int deviceid, int enabled, int delta);
int set_edge_scroll(int deviceid, int enabled);
int set_two_finger_scroll(int deviceid, int vert_enabled, int horiz_enabled);
int set_tab_to_click(int deviceid, int enabled, int left_handed);

int set_keyboard_repeat(int repeat,
                        unsigned int delay, unsigned int interval);

int is_wacom_device(int deviceid);

int listen_device_changed();
void end_device_listen_thread();

#endif
