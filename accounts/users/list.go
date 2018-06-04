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
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	userFilePasswd    = "/etc/passwd"
	userFileShadow    = "/etc/shadow"
	userFileGroup     = "/etc/group"
	userFileLoginDefs = "/etc/login.defs"
	userFileSudoers   = "/etc/sudoers"

	itemLenPasswd    = 7
	itemLenShadow    = 9
	itemLenGroup     = 4
	itemLenLoginDefs = 2
)

var (
	invalidShells = []string{
		"false",
		"nologin",
	}
)

type UserInfo struct {
	Name    string
	Uid     string
	Gid     string
	comment string
	Home    string
	Shell   string
}

func (u *UserInfo) Comment() *CommentInfo {
	return newCommentInfo(u.comment)
}

type UserInfos []UserInfo

func GetAllUserInfos() (UserInfos, error) {
	return getUserInfosFromFile(userFilePasswd)
}

func GetHumanUserInfos() (UserInfos, error) {
	infos, err := getUserInfosFromFile(userFilePasswd)
	if err != nil {
		return nil, err
	}

	infos = infos.filterUserInfos()

	return infos, nil
}

func GetUserInfoByName(name string) (UserInfo, error) {
	return getUserInfo(UserInfo{Name: name}, userFilePasswd)
}

func GetUserInfoByUid(uid string) (UserInfo, error) {
	return getUserInfo(UserInfo{Uid: uid}, userFilePasswd)
}

func getUserInfo(condition UserInfo, file string) (UserInfo, error) {
	infos, err := getUserInfosFromFile(file)
	if err != nil {
		return UserInfo{}, err
	}

	for _, info := range infos {
		if info.Name == condition.Name ||
			info.Uid == condition.Uid {
			return info, nil
		}
	}

	return UserInfo{}, fmt.Errorf("Invalid username or uid")
}

func getUserInfosFromFile(file string) (UserInfos, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var infos UserInfos
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		items := strings.Split(line, ":")
		if len(items) != itemLenPasswd {
			continue
		}

		info := UserInfo{
			Name:    items[0],
			Uid:     items[2],
			Gid:     items[3],
			comment: items[4],
			Home:    items[5],
			Shell:   items[6],
		}

		infos = append(infos, info)
	}

	return infos, nil
}

func (infos UserInfos) GetUserNames() []string {
	var names []string
	for _, info := range infos {
		names = append(names, info.Name)
	}

	return names
}

func (infos UserInfos) filterUserInfos() UserInfos {
	var tmp UserInfos
	for _, info := range infos {
		if !info.isHumanUser(userFileShadow, userFileLoginDefs) {
			continue
		}

		tmp = append(tmp, info)
	}

	return tmp
}

func (info UserInfo) isHumanUser(configShadow string, configLoginDefs string) bool {
	if info.Name == "root" {
		return false
	}

	if CanNoPasswdLogin(info.Name) {
		return true
	}

	if !info.isHumanViaShell() {
		return false
	}

	if !info.isHumanViaShadow(configShadow) {
		return false
	}

	if !info.isHumanViaLoginDefs(configLoginDefs) {
		return false
	}

	return true
}

func (info UserInfo) isHumanViaShell() bool {
	items := strings.Split(info.Shell, "/")
	if len(items) == 0 {
		return true
	}

	if isStrInArray(items[len(items)-1], invalidShells) {
		return false
	}

	return true
}

func (info UserInfo) isHumanViaShadow(config string) bool {
	content, err := ioutil.ReadFile(config)
	if err != nil {
		return false
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		items := strings.Split(line, ":")
		if len(items) != itemLenShadow {
			continue
		}

		if items[0] != info.Name {
			continue
		}

		pw := items[1]
		if pw == "" {
			break
		}

		// user was locked
		if pw[0] == '!' {
			return true
		}

		// 加盐密码最短为13
		if pw[0] == '*' || len(pw) < 13 {
			break
		}

		return true
	}

	return false
}

func (info UserInfo) isHumanViaLoginDefs(config string) bool {
	fr, err := os.Open(config)
	if err != nil {
		return false
	}
	defer fr.Close()
	var (
		found  int
		uidMin string
		uidMax string

		scanner = bufio.NewScanner(fr)
	)

	for scanner.Scan() {
		if found == 2 {
			break
		}

		var line = scanner.Text()

		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		items := strings.Fields(line)
		if len(items) != itemLenLoginDefs {
			continue
		}

		if items[0] == "UID_MIN" {
			uidMin = items[1]
			found += 1
			continue
		}

		if items[0] == "UID_MAX" {
			uidMax = items[1]
			found += 1
		}
	}

	if len(uidMax) == 0 || len(uidMin) == 0 {
		return false
	}

	uidMinInt, err := strconv.Atoi(uidMin)
	if err != nil {
		return false
	}

	uidMaxInt, err := strconv.Atoi(uidMax)
	if err != nil {
		return false
	}

	uidInt, err := strconv.Atoi(info.Uid)
	if err != nil {
		return false
	}

	if uidInt > uidMaxInt || uidInt < uidMinInt {
		return false
	}

	return true
}
