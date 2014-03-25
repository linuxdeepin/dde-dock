/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

package main

import (
        "io/ioutil"
        "os/user"
        "strconv"
        "strings"
)

const (
        CMD_USERADD = "/usr/sbin/useradd"
        CMD_USERDEL = "/usr/sbin/userdel"

        ACCOUNT_TYPE_STANDARD      = 0
        ACCOUNT_TYPE_ADMINISTACTOR = 1
)

func (op *AccountManager) CreateUser(name, fullname string, accountTyte int32) string {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In CreateUser:%v",
                                err)
                }
        }()

        //authWithPolkit(POLKIT_MANAGER_USER)

        args := []string{}

        args = append(args, "-m")
        args = append(args, "-s")
        args = append(args, "/bin/bash")
        args = append(args, "-c")
        args = append(args, fullname)
        args = append(args, name)
        execCommand(CMD_USERADD, args)

        info, _ := getInfoViaName(name)
        newUser := newUserManager(info.Uid)
        //newUser.AccountType = accountTyte
        newUser.applyPropertiesChanged("AccountType", accountTyte)
        newUser.updateUserInfo()

        path := op.FindUserByName(name)
        //op.UserAdded(path)
        //op.setPropName("UserList")
        return path
}

func (op *AccountManager) DeleteUser(name string, removeFiles bool) {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In DeleteUser:%v",
                                err)
                }
        }()

        //authWithPolkit(POLKIT_MANAGER_USER)

        args := []string{}

        if removeFiles {
                args = append(args, "-r")
        }
        args = append(args, name)

        //path := op.FindUserByName(name)
        execCommand(CMD_USERDEL, args)
        //op.UserDeleted(path)
        //op.setPropName("UserList")
}

func (op *AccountManager) FindUserById(id string) string {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In FindUserById:%v",
                                err)
                }
        }()

        path := USER_MANAGER_PATH + id
        op.setPropName("UserList")

        for _, v := range op.UserList {
                if path == v {
                        return path
                }
        }

        return ""
}

func (op *AccountManager) FindUserByName(name string) string {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In FindUserByName:%v", err)
                }
        }()

        userInfo, err := user.Lookup(name)
        if err != nil {
                logObject.Warningf("Lookup By Name Failed:%v", err)
                return ""
        }

        return op.FindUserById(userInfo.Uid)
}

func (op *AccountManager) RandUserIcon() (string, bool) {
        if icon := getRandUserIcon(); len(icon) > 0 {
                return icon, true
        }

        return "", false
}

func getInfoViaUid(uid string) (UserInfo, bool) {
        infos := getUserInfoList()

        for _, info := range infos {
                if info.Uid == uid {
                        return info, true
                }
        }

        return UserInfo{}, false
}

func getInfoViaName(name string) (UserInfo, bool) {
        infos := getUserInfoList()

        for _, info := range infos {
                if info.Name == name {
                        return info, true
                }
        }

        return UserInfo{}, false
}

func getUserInfoList() []UserInfo {
        contents, err := ioutil.ReadFile(ETC_PASSWD)
        if err != nil {
                logObject.Warningf("ReadFile '%s' failed: %s", ETC_PASSWD, err)
                panic(err)
        }

        infos := []UserInfo{}
        lines := strings.Split(string(contents), "\n")
        for _, line := range lines {
                strs := strings.Split(line, ":")

                /* len of each line in /etc/passwd by spliting ':' is 7 */
                if len(strs) != PASSWD_SPLIT_LEN {
                        continue
                }

                info := newUserInfo(strs[0], strs[2], strs[3],
                        strs[5], strs[6])
                if userIsHuman(&info) {
                        infos = append(infos, info)
                }
        }

        return infos
}

func newUserInfo(name, uid, gid, home, shell string) UserInfo {
        info := UserInfo{}

        info.Name = name
        info.Uid = uid
        info.Gid = gid
        info.Home = home
        info.Shell = shell

        return info
}

func userIsFilterList(name string) bool {
        return opUtils.IsElementExist(name, filterList)
}

func userIsHuman(info *UserInfo) bool {
        if userIsFilterList(info.Name) {
                return false
        }

        shells := strings.Split(info.Shell, "/")
        tmpShell := shells[len(shells)-1]
        if SHELL_END_FALSE == tmpShell ||
                SHELL_END_NOLOGIN == tmpShell {
                return false
        }

        if !detetedViaShadowFile(info) {
                id, _ := strconv.ParseInt(info.Uid, 10, 64)
                if id < 1000 {
                        return false
                }
        }

        return true
}

func detetedViaShadowFile(info *UserInfo) bool {
        contents, err := ioutil.ReadFile(ETC_SHADOW)
        if err != nil {
                logObject.Warningf("ReadFile '%s' failed: %s", ETC_SHADOW, err)
                panic(err)
        }

        isHuman := false
        info.Locked = false
        lines := strings.Split(string(contents), "\n")
        for _, line := range lines {
                strs := strings.Split(line, ":")
                if len(strs) != SHADOW_SPLIT_LEN {
                        continue
                }

                if strs[0] != info.Name {
                        continue
                }
                pw := strs[1]
                /*
                   // modern hashes start with "$n$" && len is 98
                   if pw[0] == '$' {
                           if len(pw) < 4 {
                                   continue
                           }
                   } else if pw[0] == '!' {
                           info.Locked = true
                           id, _ := strconv.ParseInt(info.Uid, 10, 64)
                           if id < 1000 {
                                   continue
                           }
                   } else if pw[0] != '.' || pw[0] != '/' ||
                           !charIsAlNum(pw[0]) {
                           // DES crypt is base64 encoded [./A-Za-z0-9]
                           continue
                   }
                */
                //加盐密码最短为13
                if len(pw) < 13 {
                        continue
                }

                if pw[0] == '!' {
                        info.Locked = true
                }

                isHuman = true
        }

        return isHuman
}
