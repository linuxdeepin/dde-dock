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

package accounts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"pkg.linuxdeepin.com/lib/graphic"
	"strings"
)

const (
	polkitManagerUser    = "com.deepin.daemon.accounts.user-administration"
	polkitChangeOwnData  = "com.deepin.daemon.accounts.change-own-user-data"
	polkitSetLoginOption = "com.deepin.daemon.accounts.set-login-option"
)

type ErrCodeType int32

const (
	ErrCodeUnkown ErrCodeType = iota
	ErrCodeAuthFailed
	ErrCodeExecFailed
	ErrCodeParamInvalid
)

func (code ErrCodeType) String() string {
	switch code {
	case ErrCodeUnkown:
		return "Unkown error"
	case ErrCodeAuthFailed:
		return "Policykit authentication failed"
	case ErrCodeExecFailed:
		return "Exec command failed"
	case ErrCodeParamInvalid:
		return "Invalid parameters"
	}

	return "Unkown error"
}

func clearUserDatas(name string) {
	icons := getUserCustomIcons(name)
	config := path.Join(userConfigDir, name)

	icons = append(icons, config)
	for _, v := range icons {
		os.Remove(v)
	}
}

func getUserStandardIcons() []string {
	imgs, err := graphic.GetImagesInDir(userIconsDir)
	if err != nil {
		return nil
	}

	var icons []string
	for _, img := range imgs {
		if strings.Contains(img, "guest") {
			continue
		}

		icons = append(icons, img)
	}

	return icons
}

func getUserCustomIcons(name string) []string {
	return getUserIconsFromDir(userCustomIconsDir, name+"-")
}

func getUserIconsFromDir(dir, condition string) []string {
	imgs, err := graphic.GetImagesInDir(dir)
	if err != nil {
		return nil
	}

	var icons []string
	for _, img := range imgs {
		if !strings.Contains(img, condition) {
			continue
		}

		icons = append(icons, img)
	}

	return icons
}

func isStrInArray(str string, array []string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}

	return false
}

func polkitAuthManagerUser(pid uint32) error {
	return polkitAuthentication(polkitManagerUser, pid)
}

func polkitAuthChangeOwnData(pid uint32) error {
	return polkitAuthentication(polkitChangeOwnData, pid)
}

func polkitAuthentication(action string, pid uint32) error {
	success, err := polkitAuthWithPid(action, pid)
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf(ErrCodeAuthFailed.String())
	}

	return nil
}

const (
	pidFileStatus = "/proc/%v/status"
)

func getUidByPid(pid uint32) (string, error) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Recover error in getUidByPid:", err)
		}
	}()

	var file = fmt.Sprintf(pidFileStatus, pid)
	datas, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	var lines = strings.Split(string(datas), "\n")
	for _, line := range lines {
		if !strings.Contains(line, "Uid:") {
			continue
		}

		strv := strings.Split(line, "\t")
		return strv[1], nil
	}

	return "", fmt.Errorf("Invalid file: %s", file)
}
