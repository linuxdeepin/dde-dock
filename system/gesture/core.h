/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef __GESTURE_CORE_H__
#define __GESTURE_CORE_H__

#define GESTURE_TYPE_SWIPE 100
#define GESTURE_TYPE_PINCH 101

// swipe
#define GESTURE_DIRECTION_UP 10
#define GESTURE_DIRECTION_DOWN 11
#define GESTURE_DIRECTION_LEFT 12
#define GESTURE_DIRECTION_RIGHT 13
// pinch
#define GESTURE_DIRECTION_IN 14
#define GESTURE_DIRECTION_OUT 15

int start_loop();
void quit_loop();

#endif
