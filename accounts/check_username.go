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
	"io/ioutil"
	"regexp"
	"strings"
)

func getListFromFile(filename, sep string) []string {
	list := []string{}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Warningf("ReadFile '%s' failed: %v", filename, err)
		return list
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if len(line) < 1 {
			continue
		}

		strs := strings.Split(line, sep)
		if len(strs) < 1 {
			continue
		}

		list = append(list, strs[0])
	}

	return list
}

func getUsernameList() []string {
	return getListFromFile(ETC_PASSWD, ":")
}

func getGroupnameList() []string {
	return getListFromFile(ETC_GROUP, ":")
}

func isUsernameValid(username string) bool {
	/**
	 * The user name is allowed only started with a letter,
	 * and is composed of letters and numbers
	 */
	//match, err := regexp.Compile(`^[A-Za-z][A-Za-z0-9]+$`)
	match, err := regexp.Compile(`^[a-z][a-z0-9_-]+$`)
	if err != nil {
		logger.Warning("New Compile Failed:", err)
		return false
	}

	if !match.MatchString(username) {
		return false
	}

	return true
}

func isUserExist(username string) bool {
	userList := getUsernameList()
	groupList := getGroupnameList()

	if strIsInList(username, userList) ||
		strIsInList(username, groupList) ||
		strings.ToLower(username) == "guest" {
		return false
	}

	return true
}

func isPasswordValid(passwd string) bool {
	return true
	match, err := regexp.Compile(`^[!-~]+$`)
	if err != nil {
		logger.Warning("Check passwd.New match failed:", err)
		return false
	}

	if !match.MatchString(passwd) {
		return false
	}

	return true
}
