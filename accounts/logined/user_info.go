/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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

package logined

import (
	"dbus/org/freedesktop/login1"
	"fmt"
	"pkg.deepin.io/lib/dbus"
)

// UserInfo Show logined user info, if type is tty or ssh, no desktop and display
type UserInfo struct {
	// Active  bool
	UID     uint32
	Name    string
	Desktop string
	Display string

	userPath dbus.ObjectPath
}

// UserInfos Logined user list
type UserInfos []*UserInfo

func newUserInfo(userPath dbus.ObjectPath) (*UserInfo, error) {
	core, err := login1.NewUser(dbusLogin1Dest, userPath)
	if err != nil {
		return nil, err
	}
	defer login1.DestroyUser(core)

	var info = UserInfo{
		UID:      core.UID.Get(),
		Name:     core.Name.Get(),
		userPath: userPath,
	}

	if info.Name == "" {
		return nil, fmt.Errorf("Invalid user path: %s", userPath)
	}

	list := core.Display.Get()
	if len(list) != 2 {
		return &info, nil
	}

	session, err := login1.NewSession(dbusLogin1Dest, list[1].(dbus.ObjectPath))
	if err != nil {
		return &info, nil
	}
	defer login1.DestroySession(session)

	//info.Active = session.Active.Get()
	info.Display = session.Display.Get()
	info.Desktop = session.Desktop.Get()

	return &info, nil
}

// Add Add user to list, if exist and equal, return false
// else replace it, return true
func (infos UserInfos) Add(info *UserInfo) (UserInfos, bool) {
	idx := infos.Index(info.UID)
	if idx != -1 {
		if infos[idx].Equal(info) {
			return infos, false
		}
		infos[idx] = info
	} else {
		infos = append(infos, info)
	}
	return infos, true
}

// Index Find the user position in list, if not found, return -1
func (infos UserInfos) Index(uid uint32) int32 {
	for i, v := range infos {
		if v.UID != uid {
			continue
		}

		return int32(i)
	}
	return -1
}

// Delete Delete user from list, if deleted, return true
func (infos UserInfos) Delete(uid uint32) (UserInfos, bool) {
	var (
		tmp     UserInfos
		deleted = false
	)
	for _, v := range infos {
		if v.UID == uid {
			deleted = true
			continue
		}
		tmp = append(tmp, v)
	}
	return tmp, deleted
}

// Equal Check whether equal with target
func (info *UserInfo) Equal(target *UserInfo) bool {
	if info.Name != target.Name || info.UID != target.UID ||
		info.Desktop != target.Desktop || info.Display != target.Display {
		return false
	}
	return true
}
