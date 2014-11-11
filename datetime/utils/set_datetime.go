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

package utils

import (
	"dbus/com/deepin/api/setdatetime"
	"fmt"
)

var (
	errUninitialized = fmt.Errorf("SetDateTime Uninitialized")
)

var _setDate *setdatetime.SetDateTime

func InitSetDateTime() error {
	if _setDate != nil {
		return nil
	}

	var err error
	_setDate, err = setdatetime.NewSetDateTime(
		"com.deepin.api.SetDateTime",
		"/com/deepin/api/SetDateTime",
	)
	if err != nil {
		return err
	}

	return nil
}

func DestroySetDateTime() {
	if _setDate == nil {
		return
	}

	setdatetime.DestroySetDateTime(_setDate)
	_setDate = nil
}

func SetDate(value string) error {
	if _setDate == nil {
		return errUninitialized
	}

	_, err := _setDate.SetCurrentDate(value)
	if err != nil {
		return err
	}

	return nil
}

func SetTime(value string) error {
	if _setDate == nil {
		return errUninitialized
	}

	_, err := _setDate.SetCurrentTime(value)
	if err != nil {
		return err
	}

	return nil
}
