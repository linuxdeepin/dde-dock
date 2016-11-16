/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef __GESTURE_UTILS_H__
#define __GESTURE_UTILS_H__

#include <libinput.h>

struct libinput* open_from_udev(char *seat, void *user_data, int verbose);
struct libinput* open_from_path(char **path, void *user_data, int verbose);

#endif
