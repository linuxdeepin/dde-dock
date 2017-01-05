/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef __GESTURE_CORE_H__
#define __GESTURE_CORE_H__

#define GESTURE_TYPE_SWIPE 100
#define GESTURE_TYPE_PINCH 101
#define GESTURE_TYPE_TAP 102

// tap
#define GESTURE_DIRECTION_NONE 0
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
