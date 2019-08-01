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
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
)

var (
	errInvalidParam = fmt.Errorf("Invalid or empty parameter")
)

var (
	groupFileTimestamp int64 = 0
	groupFileInfo            = make(map[string][]string)
	groupFileLocker    sync.Mutex

	groupNameNoPasswdLogin = "nopasswdlogin"
)

const CommentFieldsLen = 5

// CommentInfo is passwd file user comment info
type CommentInfo [CommentFieldsLen]string

func newCommentInfo(comment string) *CommentInfo {
	var ci CommentInfo
	parts := strings.Split(comment, ",")

	// length is min(CommentFieldsLen, len(parts))
	length := len(parts)
	if length > CommentFieldsLen {
		length = CommentFieldsLen
	}

	copy(ci[:], parts[:length])
	return &ci
}

func (ci *CommentInfo) String() string {
	return strings.Join(ci[:], ",")
}

func (ci *CommentInfo) FullName() string {
	return ci[0]
}

func (ci *CommentInfo) SetFullName(value string) {
	ci[0] = value
}

func isCommentFieldValid(name string) bool {
	if strings.ContainsAny(name, ",=:\n") {
		return false
	}
	return true
}

func ModifyFullName(newname, username string) error {
	if !isCommentFieldValid(newname) {
		return errors.New("invalid nickname")
	}

	user, err := GetUserInfoByName(username)
	if err != nil {
		return err
	}
	comment := user.Comment()
	comment.SetFullName(newname)
	return modifyComment(comment.String(), username)
}

func modifyComment(comment, username string) error {
	cmd := exec.Command(userCmdModify, "-c", comment, username)
	return cmd.Run()
}

func ModifyHome(dir, username string) error {
	if len(dir) == 0 {
		return errInvalidParam
	}

	return doAction(userCmdModify, []string{"-m", "-d", dir, username})
}

func ModifyShell(shell, username string) error {
	if len(shell) == 0 {
		return errInvalidParam
	}

	return doAction(userCmdModify, []string{"-s", shell, username})
}

func ModifyPasswd(words, username string) error {
	if len(words) == 0 {
		return errInvalidParam
	}

	return updatePasswd(words, username)
}

const (
	// Same as the abbreviation in `passwd --status`
	PasswordStatusUsable     = "P"
	PasswordStatusNoPassword = "NP"
	PasswordStatusLocked     = "L"
)

func GetUserPasswordStatus(username string) (string, error) {
	content, err := ioutil.ReadFile(userFileShadow)
	if err != nil {
		return "", err
	}
	lines := bytes.Split(content, []byte{'\n'})
	for _, line := range lines {
		fields := bytes.Split(line, []byte{':'})
		if len(fields) != itemLenShadow {
			continue
		}

		if string(fields[0]) == username {
			pw := fields[1]
			if len(pw) == 0 {
				return PasswordStatusNoPassword, nil
			}
			if pw[0] == '!' || pw[0] == '*' {
				return PasswordStatusLocked, nil
			}
			return PasswordStatusUsable, nil
		}
	}

	return "", errors.New("user not found")
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

func CanNoPasswdLogin(username string) bool {
	return isUserInGroup(username, groupNameNoPasswdLogin)
}

func EnableNoPasswdLogin(username string, enabled bool) error {
	if !isGroupExists(groupNameNoPasswdLogin) {
		doAction("groupadd", []string{"-r", groupNameNoPasswdLogin})
	}

	var err error
	exists := isUserInGroup(username, groupNameNoPasswdLogin)
	if enabled {
		if !exists {
			err = doAction(userCmdGroup, []string{"-a", username, groupNameNoPasswdLogin})
		}
	} else {
		if exists {
			err = doAction(userCmdGroup, []string{"-d", username, groupNameNoPasswdLogin})
		}
	}
	return err
}

func getAdminUserList(fileGroup, fileSudoers string) ([]string, error) {
	groups, users, err := getAdmGroupAndUser(fileSudoers)
	if err != nil {
		return nil, err
	}

	groupFileLocker.Lock()
	defer groupFileLocker.Unlock()
	infos, err := getGroupInfoWithCache(fileGroup)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		v, ok := infos[group]
		if !ok {
			continue
		}
		users = append(users, v...)
	}
	return users, nil
}

var (
	_admGroups       []string
	_admUsers        []string
	_admTimestampMap = make(map[string]int64)
)

// get adm group and user from '/etc/sudoers'
func getAdmGroupAndUser(file string) ([]string, []string, error) {
	finfo, err := os.Stat(file)
	if err != nil {
		return nil, nil, err
	}
	timestamp := finfo.ModTime().Unix()
	if t, ok := _admTimestampMap[file]; ok && t == timestamp {
		return _admGroups, _admUsers, nil
	}

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
	_admGroups, _admUsers = groups, users
	_admTimestampMap[file] = timestamp
	return groups, users, nil
}

func isGroupExists(group string) bool {
	groupFileLocker.Lock()
	defer groupFileLocker.Unlock()
	infos, err := getGroupInfoWithCache(userFileGroup)
	if err != nil {
		return false
	}
	_, ok := infos[group]
	return ok
}

func isUserInGroup(user, group string) bool {
	groupFileLocker.Lock()
	defer groupFileLocker.Unlock()
	infos, err := getGroupInfoWithCache(userFileGroup)
	if err != nil {
		return false
	}
	v, ok := infos[group]
	if !ok {
		return false
	}
	return isStrInArray(user, v)
}

func GetUserGroups(user string) ([]string, error) {
	groupFileLocker.Lock()
	defer groupFileLocker.Unlock()
	infos, err := getGroupInfoWithCache(userFileGroup)
	if err != nil {
		return nil, err
	}

	var result []string
	for groupName, users := range infos {
		if groupName == user {
			result = append(result, groupName)
			continue
		}
		for _, u := range users {
			if u == user {
				result = append(result, groupName)
				break
			}
		}
	}
	sort.Strings(result)
	return result, nil
}

func GetAllGroups() ([]string, error) {
	groupFileLocker.Lock()
	defer groupFileLocker.Unlock()
	infos, err := getGroupInfoWithCache(userFileGroup)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(infos))
	idx := 0
	for groupName := range infos {
		result[idx] = groupName
		idx++
	}
	sort.Strings(result)
	return result, nil
}

func getGroupInfoWithCache(file string) (map[string][]string, error) {
	info, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	if groupFileTimestamp == info.ModTime().UnixNano() &&
		len(groupFileInfo) != 0 {
		return groupFileInfo, nil
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	groupFileTimestamp = info.ModTime().UnixNano()
	groupFileInfo = parseGroup(content)
	return groupFileInfo, nil
}

func parseGroup(data []byte) map[string][]string {
	result := make(map[string][]string)
	lines := bytes.Split(data, []byte{'\n'})
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		items := bytes.Split(line, []byte{':'})
		if len(items) != itemLenGroup {
			continue
		}

		groupName := string(items[0])
		users0 := bytes.Split(items[3], []byte{','})
		// convert users0 to users
		var users []string
		if len(users0) > 0 {
			users = make([]string, len(users0))
			for idx := range users {
				users[idx] = string(users0[idx])
			}
		}
		result[groupName] = users
	}

	return result
}
