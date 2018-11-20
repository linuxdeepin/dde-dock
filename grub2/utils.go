/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package grub2

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"pkg.deepin.io/dde/daemon/grub_common"

	polkit "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.policykit1"
	"pkg.deepin.io/lib/dbus1"
)

func quoteString(str string) string {
	return strconv.Quote(str)
}

func checkGfxmode(v string) error {
	if v == "auto" {
		return nil
	}

	_, err := grub_common.ParseGfxmode(v)
	return err
}

func getStringIndexInArray(a string, list []string) int {
	for i, b := range list {
		if b == a {
			return i
		}
	}
	return -1
}

var noCheckAuth bool

func allowNoCheckAuth() {
	if os.Getenv("NO_CHECK_AUTH") == "1" {
		noCheckAuth = true
		return
	}
}

func checkAuthWithPid(pid uint32, actionId string) (bool, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return false, err
	}
	authority := polkit.NewAuthority(systemBus)
	subject := polkit.MakeSubject(polkit.SubjectKindUnixProcess)
	subject.SetDetail("pid", pid)
	subject.SetDetail("start-time", uint64(0))
	result, err := authority.CheckAuthorization(0, subject, actionId, nil,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}

	return result.IsAuthorized, nil
}

var errAuthFailed = errors.New("authentication failed")

func getBytesMD5Sum(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

func getFileMD5sum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, nil
}
