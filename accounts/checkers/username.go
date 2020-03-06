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
	"regexp"
	"strconv"
	"strings"
)

func Tr(text string) string {
	return text
}

const (
	userNameMaxLength = 32
	userNameMinLength = 3

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
	ErrCodeFirstNotAlphabetic
	ErrCodeExist
	ErrCodeSystemUsed
	ErrCodeLen
)

func (code ErrorCode) Error() *ErrorInfo {
	var err error
	switch code {
	case ErrCodeEmpty:
		err = fmt.Errorf(Tr("Username cannot be empty"))
	case ErrCodeInvalidChar:
		err = fmt.Errorf(Tr("Username must only contain a~z, A-Z, 0~9, - or _"))
	case ErrCodeFirstNotAlphabetic:
		err = fmt.Errorf(Tr("Username must begin with an alphabetic character"))
	case ErrCodeExist:
		err = fmt.Errorf(Tr("The username already exists"))
	case ErrCodeSystemUsed:
		err = fmt.Errorf(Tr("The username has been used by system"))
	case ErrCodeLen:
		err = fmt.Errorf(Tr("Username must be between 3 and 32 characters"))
	default:
		return nil
	}

	return &ErrorInfo{
		Code:  code,
		Error: err,
	}
}

func CheckUsernameValid(name string) *ErrorInfo {
	length := len(name)
	if length == 0 {
		return ErrCodeEmpty.Error()
	}

	if length > userNameMaxLength || length < userNameMinLength {
		return ErrCodeLen.Error()
	}

	if Username(name).isNameExist() {
		id, err := Username(name).getUid()
		if err != nil || id >= 1000 {
			return ErrCodeExist.Error()
		} else {
			return ErrCodeSystemUsed.Error()
		}
	}

	if !Username(name).isAlphabeticCharStart() {
		return ErrCodeFirstNotAlphabetic.Error()
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

func (name Username) isAlphabeticCharStart() bool {
	match := regexp.MustCompile(`^[a-zA-Z]`)
	if !match.MatchString(string(name)) {
		return false
	}

	return true
}

func (name Username) isStringValid() bool {
	match := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
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
