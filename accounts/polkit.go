/*
 * Copyright (C) 2013 ~ 2017 Deepin Technology Co., Ltd.
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
	"dbus/org/freedesktop/policykit1"
	"fmt"
	"pkg.deepin.io/lib/dbus"
)

const (
	polkitSender = "org.freedesktop.PolicyKit1"
	polkitPath   = "/org/freedesktop/PolicyKit1/Authority"
)

type polkitSubject struct {
	/*
	 * The following kinds of subjects are known:
	 * Unix Process: should be set to unix-process with keys
	 *                  pid (of type uint32) and
	 *                  start-time (of type uint64)
	 * Unix Session: should be set to unix-session with the key
	 *                  session-id (of type string)
	 * System Bus Name: should be set to system-bus-name with the key
	 *                  name (of type string)
	 */
	SubjectKind    string
	SubjectDetails map[string]dbus.Variant
}

func polkitAuthWithPid(actionId string, pid uint32) (bool, error) {
	polkit, err := policykit1.NewAuthority(polkitSender,
		polkitPath)
	if err != nil {
		return false, err
	}

	var subject polkitSubject
	subject.SubjectKind = "unix-process"
	subject.SubjectDetails = make(map[string]dbus.Variant)
	subject.SubjectDetails["pid"] = dbus.MakeVariant(uint32(pid))
	subject.SubjectDetails["start-time"] = dbus.MakeVariant(uint64(0))

	var details = make(map[string]string)
	details[""] = ""
	var flg uint32 = 1
	var cancelId string

	ret, err := polkit.CheckAuthorization(subject,
		actionId, details, flg, cancelId)
	if err != nil {
		return false, err
	}

	//If ret[0].(bool) == true, successful.
	if len(ret) == 0 {
		return false, fmt.Errorf("No results returned")
	}

	v, ok := ret[0].(bool)
	if !ok {
		return false, fmt.Errorf("Invalid result type")
	}

	return v, nil
}
