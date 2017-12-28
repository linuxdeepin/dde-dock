/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package x_event_monitor

import (
	"strings"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
)

const (
	MotionFlag = int32(1 << 0)
	ButtonFlag = int32(1 << 1)
	KeyFlag    = int32(1 << 2)
	AllFlag    = MotionFlag | ButtonFlag | KeyFlag
)

var hasMotionFlag = func() func(int32) bool {
	return func(flag int32) bool {
		if flag < 0 || flag > AllFlag {
			return false
		}

		if flag&MotionFlag == MotionFlag {
			return true
		}

		return false
	}
}()

var hasButtonFlag = func() func(int32) bool {
	return func(flag int32) bool {
		if flag < 0 || flag > AllFlag {
			return false
		}

		if flag&ButtonFlag == ButtonFlag {
			return true
		}

		return false
	}
}()

var hasKeyFlag = func() func(int32) bool {
	return func(flag int32) bool {
		if flag < 0 || flag > AllFlag {
			return false
		}

		if flag&KeyFlag == KeyFlag {
			return true
		}

		return false
	}
}()

func isInArea(x, y int32, area coordinateRange) bool {
	if (x >= area.X1 && x <= area.X2) &&
		(y >= area.Y1 && y <= area.Y2) {
		return true
	}

	return false
}

func isInIdList(md5Str string, list []string) bool {
	for _, v := range list {
		if md5Str == v {
			return true
		}
	}

	return false
}

var keyCode2Str = func() func(int32) string {
	XU, err := xgbutil.NewConn()
	if err != nil {
		logger.Error("Can't connect to Xserver")
		return func(int32) string { return "" }
	}
	keybind.Initialize(XU)
	return func(code int32) string {
		keyStr := keybind.LookupString(XU, 0, xproto.Keycode(code))
		if keyStr == " " {
			keyStr = "space"
		}
		keyStr = strings.ToLower(keyStr)
		return keyStr
	}
}()
