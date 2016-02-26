/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
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
	groups, users, err := getAdmGroupAndUser(fileSudoers)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(fileGroup)
	if err != nil {
		return nil, err
	}

	var list []string = users
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		items := strings.Split(line, ":")
		if len(items) != itemLenGroup {
			continue
		}

		if !isStrInArray(items[0], groups) {
			continue
		}

		list = append(list, strings.Split(items[3], ",")...)
	}

	return list, nil
}

// get adm group and user from '/etc/sudoers'
func getAdmGroupAndUser(file string) ([]string, []string, error) {
	fr, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer fr.Close()

	var (
		groups  []string
		users   []string
		scanner = bufio.NewScanner(fr)
	)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		line = strings.TrimSpace(line)
		if line[0] == '#' || !strings.Contains(line, `ALL=(ALL`) {
			continue
		}

		array := strings.Split(line, "ALL")
		// admin group
		if line[0] == '%' {
			// deepin: %sudo\tALL=(ALL:ALL) ALL
			// archlinux: %wheel ALL=(ALL) ALL
			array = strings.Split(array[0], "%")
			tmp := strings.TrimRight(array[1], "\t")
			groups = append(groups, strings.TrimSpace(tmp))
		} else {
			// admin user
			// deepin: root\tALL=(ALL:ALL) ALL
			// archlinux: root ALL=(ALL) ALL
			tmp := strings.TrimRight(array[0], "\t")
			users = append(users, strings.TrimSpace(tmp))
		}
	}
	return groups, users, nil
}
