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

package checkers

import (
	"fmt"
	"io/ioutil"
	"os/user"
	. "pkg.deepin.io/lib/gettext"
	"regexp"
	"strconv"
	"strings"
)

const (
	userNameMaxLength = 32

	passwdFile = "/etc/passwd"
	groupFile  = "/etc/group"
)

type ErrorCode int32

type ErrorInfo struct {
	Code  ErrorCode
	Error error
}

const (
	ErrCodeEmpty ErrorCode = iota + 1
	ErrCodeInvalidChar
	ErrCodeFirstNotLower
	ErrCodeExist
	ErrCodeSystemUsed
	ErrCodeLenMoreThen
)

func (code ErrorCode) Error() *ErrorInfo {
	var err error
	switch code {
	case ErrCodeEmpty:
		err = fmt.Errorf(Tr("Username can not be empty."))
	case ErrCodeInvalidChar:
		err = fmt.Errorf(Tr("Username must comprise a~z, 0~9, - or _."))
	case ErrCodeFirstNotLower:
		err = fmt.Errorf(Tr("The first character must be in lower case."))
	case ErrCodeExist:
		err = fmt.Errorf(Tr("The username exists."))
	case ErrCodeSystemUsed:
		err = fmt.Errorf(Tr("The username has been used by system."))
	case ErrCodeLenMoreThen:
		err = fmt.Errorf(Tr("The username's length exceeds the limit"))
	default:
		return nil
	}

	return &ErrorInfo{
		Code:  code,
		Error: err,
	}
}

func CheckUsernameValid(name string) *ErrorInfo {
	if len(name) == 0 {
		return ErrCodeEmpty.Error()
	}

	if len(name) > userNameMaxLength {
		return ErrCodeLenMoreThen.Error()
	}

	if Username(name).isNameExist() {
		id, err := Username(name).getUid()
		if err != nil || id >= 1000 {
			return ErrCodeExist.Error()
		} else {
			return ErrCodeSystemUsed.Error()
		}
	}

	if !Username(name).isLowerCharStart() {
		return ErrCodeFirstNotLower.Error()
	}

	if !Username(name).isStringValid() {
		return ErrCodeInvalidChar.Error()
	}

	return nil
}

type Username string

type UsernameList []string

func (name Username) isNameExist() bool {
	names, err := getAllUsername(passwdFile)
	if err != nil {
		return false
	}

	if !isStrInArray(string(name), names) {
		return false
	}

	return true
}

func (name Username) isLowerCharStart() bool {
	match := regexp.MustCompile(`^[a-z]`)
	if !match.MatchString(string(name)) {
		return false
	}

	return true
}

func (name Username) isStringValid() bool {
	match := regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
	if !match.MatchString(string(name)) {
		return false
	}

	return true
}

func (name Username) getUid() (int64, error) {
	u, err := user.Lookup(string(name))
	if err != nil {
		return -1, err
	}

	id, err := strconv.ParseInt(u.Uid, 10, 64)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func getAllUsername(file string) (UsernameList, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var names UsernameList
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		items := strings.Split(line, ":")
		if len(items) < 3 {
			continue
		}

		names = append(names, items[0])
	}

	return names, nil
}
