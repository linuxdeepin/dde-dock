/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef __GET_BACKLIGHT_PATH__
#define __GET_BACKLIGHT_PATH__

double get_backlight(char* driver_type);
void set_backlight(double v, char* driver_type);

#endif
