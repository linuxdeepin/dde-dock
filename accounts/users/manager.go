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

package users

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

const (
	UserTypeStandard int32 = iota
	UserTypeAdmin
)

const (
	userCmdAdd    = "useradd"
	userCmdDelete = "userdel"
	userCmdModify = "usermod"
	userCmdGroup  = "gpasswd"

	defaultConfigShell = "/etc/adduser.conf"
)

func CreateUser(username, fullname, shell string, ty int32) error {
	if len(username) == 0 {
		return errInvalidParam
	}

	if len(shell) == 0 {
		shell, _ = getDefaultShell(defaultConfigShell)
	}

	var args = []string{"-m"}
	if len(shell) != 0 {
		args = append(args, "-s", shell)
	}

	if len(fullname) != 0 {
		args = append(args, "-c", fullname)
	}

	args = append(args, username)
	return doAction(userCmdAdd, args)
}

func DeleteUser(rmFiles bool, username string) error {
	var args = []string{"-f"}
	if rmFiles {
		args = append(args, "-r")
	}
	args = append(args, username)

	return doAction(userCmdDelete, args)
}

func LockedUser(locked bool, username string) error {
	var arg string
	if locked {
		arg = "-L"
	} else {
		arg = "-U"
	}
	return doAction(userCmdModify, []string{arg, username})
}

func SetUserType(ty int32, username string) error {
	groups, _, _ := getAdmGroupAndUser(userFileSudoers)
	if len(groups) == 0 {
		return fmt.Errorf("No privilege user group exists")
	}
	var args []string
	switch ty {
	case UserTypeStandard:
		if !IsAdminUser(username) {
			return nil
		}

		// TODO: remove user from all privilege groups
		args = []string{"-d", username, groups[0]}
	case UserTypeAdmin:
		if IsAdminUser(username) {
			return nil
		}

		args = []string{"-a", username, groups[0]}
	default:
		return errInvalidParam
	}

	return doAction(userCmdGroup, args)
}

// Default config: /etc/adduser.conf
func getDefaultShell(config string) (string, error) {
	fp, err := os.Open(config)
	if err != nil {
		return "", err
	}
	defer fp.Close()

	var (
		shell   string
		match   = regexp.MustCompile(`^DSHELL=(.*)`)
		scanner = bufio.NewScanner(fp)
	)

	for scanner.Scan() {
		line := scanner.Text()
		fields := match.FindStringSubmatch(line)
		if len(fields) < 2 {
			continue
		}

		shell = fields[1]
		break
	}

	return shell, nil
}
