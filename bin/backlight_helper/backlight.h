/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef __BACKLIGHT_H__
#define __BACKLIGHT_H__

int init_udev();
void finalize_udev();

char** get_syspath_list(int* num);
void free_syspath_list(char** list, int num);

char* get_syspath_by_type(char* type);

int get_brightness(char* syspath);
int get_max_brightness(char* syspath);
int set_brightness(char* syspath, int value);

#endif
