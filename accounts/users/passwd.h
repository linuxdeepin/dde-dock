/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef __PASSWORD_H__
#define __PASSWORD_H__

char *mkpasswd(const char *words);

int lock_shadow_file();
int unlock_shadow_file();

#endif
