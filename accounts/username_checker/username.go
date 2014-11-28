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

package username_checker

import (
	"fmt"
	"io/ioutil"
	. "pkg.linuxdeepin.com/lib/gettext"
	"regexp"
	"strconv"
	"strings"
)

const (
	passwdFilePath = "/etc/passwd"
	groupFilePath  = "/etc/group"
)

type usernameInfo struct {
	name string
	id   string
}

type ErrorInfo struct {
	Message error
	Code    int32
}

var (
	ErrCodeEmpty         int32 = 1
	ErrCodeInvalidChar   int32 = 2
	ErrCodeFirstNotLower int32 = 3
	ErrCodeExist         int32 = 4
	ErrCodeSystemUsed    int32 = 5
)

var (
	ErrMSGEmpty         error
	ErrMSGInvalidChar   error
	ErrMSGFirstNotLower error
	ErrMSGExist         error
	ErrMSGSystemUsed    error
)

func initErrorInfo() {
	ErrMSGEmpty = fmt.Errorf(Tr("Username can not be empty."))
	ErrMSGInvalidChar = fmt.Errorf(Tr("Username must comprise a~z, 0~9, - or _."))
	ErrMSGFirstNotLower = fmt.Errorf(Tr("The first character must be in lower case."))
	ErrMSGExist = fmt.Errorf(Tr("The username exists."))
	ErrMSGSystemUsed = fmt.Errorf(Tr("The username has been used by system."))
}

func CheckUsernameValid(username string) *ErrorInfo {
	if ErrMSGEmpty == nil {
		initErrorInfo()
	}

	if len(username) == 0 {
		return &ErrorInfo{
			Message: ErrMSGEmpty,
			Code:    ErrCodeEmpty,
		}
	}

	ok, err := checkNameExist(username, passwdFilePath)
	if ok {
		return err
	}

	/**
	 * The username is allowed only started with a letter,
	 * and is composed of letters and numbers
	 */
	match := regexp.MustCompile(`^[a-z]`)
	if !match.MatchString(username) {
		return &ErrorInfo{
			Message: ErrMSGFirstNotLower,
			Code:    ErrCodeFirstNotLower,
		}
	}

	match = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
	if !match.MatchString(username) {
		return &ErrorInfo{
			Message: ErrMSGInvalidChar,
			Code:    ErrCodeInvalidChar,
		}
	}

	return nil
}

func checkNameExist(name, config string) (bool, *ErrorInfo) {
	infos, err := getNameListFromFile(config)
	if err != nil {
		return false, nil
	}

	exist, info := isNameInInfoList(name, infos)
	if !exist {
		return false, nil
	}

	interval, _ := strconv.ParseInt(info.id, 10, 64)
	if interval < 1000 {
		return true, &ErrorInfo{
			Message: ErrMSGSystemUsed,
			Code:    ErrCodeSystemUsed,
		}
	} else {
		return true, &ErrorInfo{
			Message: ErrMSGExist,
			Code:    ErrCodeExist,
		}
	}

	return false, nil
}

func getNameListFromFile(config string) ([]usernameInfo, error) {
	contents, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var infos []usernameInfo
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		tmp := strings.Split(line, ":")
		if len(tmp) < 3 {
			continue
		}

		info := usernameInfo{name: tmp[0], id: tmp[2]}
		infos = append(infos, info)
	}

	return infos, nil
}

func isNameInInfoList(name string, infos []usernameInfo) (bool, usernameInfo) {
	for _, info := range infos {
		if name == info.name {
			return true, info
		}
	}

	return false, usernameInfo{}
}
