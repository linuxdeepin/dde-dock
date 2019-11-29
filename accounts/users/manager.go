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

package users

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	userCmdAdd    = "useradd"
	userCmdDelete = "userdel"
	userCmdModify = "usermod"
	userCmdGroup  = "gpasswd"

	defaultConfigShell = "/etc/adduser.conf"
)

const (
	UserTypeStandard = iota
	UserTypeAdmin
)

func CreateUser(username, fullname, shell string) error {
	if len(username) == 0 {
		return errInvalidParam
	}

	if len(shell) == 0 {
		shell, _ = getDefaultShell(defaultConfigShell)
	}

	mockUserInfo := UserInfo{
		Name:    username,
		Uid:     "10000",
		Gid:     "10000",
		comment: fullname,
		Home:    filepath.Join("/home/", username),
		Shell:   shell,
	}
	err := mockUserInfo.checkLength()
	if err != nil {
		return err
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

var commonGroups = []string{
	"lp",
	"lpadmin",
	"netdev",
	"network",
	"sambashare",
	"scanner",
	"storage",
	"users",
}

func GetPresetGroups(userType int) []string {
	var groups []string
	switch userType {
	case UserTypeStandard:
		groups = commonGroups
	case UserTypeAdmin:
		groups = make([]string, len(commonGroups))
		copy(groups, commonGroups)

		adminGroups, _, _ := getAdmGroupAndUser(userFileSudoers)
		groups = append(groups, adminGroups...)
	default:
		return nil
	}

	return filterInvalidGroup(groups)
}

func filterInvalidGroup(groups []string) []string {
	result := make([]string, 0, len(groups))
	for _, group := range groups {
		if isGroupExists(group) {
			result = append(result, group)
		}
	}
	return result
}

func SetGroupsForUser(groups []string, user string) error {
	return doAction(userCmdModify, []string{"-G", strings.Join(groups, ","), user})
}

func AddGroupForUser(group, user string) error {
	if group == user {
		return nil
	}
	return doAction(userCmdGroup, []string{"-a", user, group})
}

func DeleteGroupForUser(group, user string) error {
	if group == user {
		return errors.New("not allowed to delete the same name group")
	}
	return doAction(userCmdGroup, []string{"-d", user, group})
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
