/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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

package accounts

import (
	"strconv"

	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.policykit1"
	"pkg.deepin.io/lib/dbus1"
)

func polkitAuthWithPid(actionId, user, uid string, pid uint32) (bool, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return false, err
	}
	polkit := policykit1.NewAuthority(systemBus)

	var subject policykit1.Subject
	subject.Kind = "unix-process"
	subject.Details = make(map[string]dbus.Variant)
	subject.Details["pid"] = dbus.MakeVariant(uint32(pid))
	subject.Details["start-time"] = dbus.MakeVariant(uint64(0))
	if user != "" {
		subject.Details["user"] = dbus.MakeVariant(user)
	}
	if uid != "" {
		v, err := strconv.ParseUint(uid, 10, 32)
		if err == nil {
			subject.Details["uid"] = dbus.MakeVariant(uint32(v))
		}
	}

	var details = make(map[string]string)
	details[""] = ""
	var flg uint32 = 1
	var cancelId string

	ret, err := polkit.CheckAuthorization(0, subject,
		actionId, details, flg, cancelId)
	if err != nil {
		return false, err
	}
	return ret.IsAuthorized, nil
}
