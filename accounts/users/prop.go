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

package users

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	errInvalidParam = fmt.Errorf("Invalid or empty parameter")
)

func ModifyName(newname, username string) error {
	if len(newname) == 0 {
		return errInvalidParam
	}

	var cmd = fmt.Sprintf("%s -l %s %s", userCmdModify, newname, username)
	return doAction(cmd)
}

func ModifyHome(dir, username string) error {
	if len(dir) == 0 {
		return errInvalidParam
	}

	var cmd = fmt.Sprintf("%s -m -d %s %s", userCmdModify, dir, username)
	return doAction(cmd)
}

func ModifyShell(shell, username string) error {
	if len(shell) == 0 {
		return errInvalidParam
	}

	var cmd = fmt.Sprintf("%s -s %s %s", userCmdModify, shell, username)
	return doAction(cmd)
}

func ModifyPasswd(words, username string) error {
	if len(words) == 0 {
		return errInvalidParam
	}

	return updatePasswd(EncodePasswd(words), username)
}

// passwd -S username
func IsUserLocked(username string) bool {
	var cmd = fmt.Sprintf("passwd -S %s", username)

	output, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return true
	}

	items := strings.Split(string(output), " ")
	if items[1] == "L" {
		return true
	}

	return false
}

func IsAutoLoginUser(username string) bool {
	name, _ := GetAutoLoginUser()
	if name == username {
		return true
	}

	return false
}

func IsAdminUser(username string) bool {
	admins, err := getAdminUserList(userFileGroup, userFileSudoers)
	if err != nil {
		return false
	}

	return isStrInArray(username, admins)
}

func getAdminUserList(fileGroup, fileSudoers string) ([]string, error) {
	adms, err := getAdmGroup(fileSudoers)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(fileGroup)
	if err != nil {
		return nil, err
	}

	var tmp string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		items := strings.Split(line, ":")
		if len(items) != itemLenGroup {
			continue
		}

		if !isStrInArray(items[0], adms) {
			continue
		}

		tmp = items[3]
	}

	return strings.Split(tmp, ","), nil
}

// get adm group from '/etc/sudoers'
func getAdmGroup(file string) ([]string, error) {
	fr, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fr.Close()

	var (
		adms    []string
		scanner = bufio.NewScanner(fr)
	)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || !strings.Contains(line, "%") {
			continue
		}

		if !strings.Contains(line, `ALL=(ALL`) {
			continue
		}

		// deepin: %sudo\tALL=(ALL:ALL) ALL
		// archlinux: %wheel ALL=(ALL) ALL
		array := strings.Split(line, "ALL")
		array = strings.Split(strings.TrimSpace(array[0]), "%")
		adms = append(adms, strings.TrimRight(array[1], "\t"))
	}

	return adms, nil
}
