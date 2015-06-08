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
	"path"
	dutils "pkg.linuxdeepin.com/lib/utils"
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

	displayManagerDefaultConfig = "/etc/X11/default-display-manager"
	lightdmDefaultConfig        = "/etc/lightdm/lightdm.conf"
	kdmDefaultConfig            = "/usr/share/config/kdm/kdmrc"
	gdmDefaultConfig            = "/etc/gdm/custom.conf"
)

func CreateUser(username, fullname, shell string, ty int32) error {
	if len(username) == 0 {
		return errInvalidParam
	}

	if len(shell) == 0 {
		shell, _ = getDefaultShell(defaultConfigShell)
	}

	var cmd = fmt.Sprintf("%s -m ", userCmdAdd)
	if len(shell) != 0 {
		cmd = fmt.Sprintf("%s -s %s", cmd, shell)
	}

	if len(fullname) != 0 {
		cmd = fmt.Sprintf("%s -c %s", cmd, fullname)
	}

	cmd = fmt.Sprintf("%s %s", cmd, username)
	return doAction(cmd)
}

func DeleteUser(rmFiles bool, username string) error {
	var cmd string
	if rmFiles {
		cmd = fmt.Sprintf("%s -rf %s", userCmdDelete, username)
	} else {
		cmd = fmt.Sprintf("%s -f %s", userCmdDelete, username)
	}

	return doAction(cmd)
}

func LockedUser(locked bool, username string) error {
	var cmd string
	if locked {
		cmd = fmt.Sprintf("%s -L %s", userCmdModify, username)
	} else {
		cmd = fmt.Sprintf("%s -U %s", userCmdModify, username)
	}

	return doAction(cmd)
}

func SetUserType(ty int32, username string) error {
	var cmd string
	switch ty {
	case UserTypeStandard:
		if !IsAdminUser(username) {
			return nil
		}

		cmd = fmt.Sprintf("%s -d %s sudo", userCmdGroup, username)
	case UserTypeAdmin:
		if IsAdminUser(username) {
			return nil
		}

		cmd = fmt.Sprintf("%s -a %s sudo", userCmdGroup, username)
	default:
		return errInvalidParam
	}

	return doAction(cmd)
}

func SetAutoLoginUser(username string) error {
	dm, err := getDefaultDisplayManager(displayManagerDefaultConfig)
	if err != nil {
		return err
	}

	name, _ := GetAutoLoginUser()
	if name == username {
		return nil
	}

	switch dm {
	case "lightdm":
		return setLightdmAutoLoginUser(username, lightdmDefaultConfig)
	case "kdm":
		return setKDMAutoLoginUser(username, kdmDefaultConfig)
	case "gdm":
		return setGDMAutoLoginUser(username, gdmDefaultConfig)
	default:
		return fmt.Errorf("Not supported or invalid display manager: %q", dm)
	}

	return nil
}

func GetAutoLoginUser() (string, error) {
	dm, err := getDefaultDisplayManager(displayManagerDefaultConfig)
	if err != nil {
		return "", err
	}

	switch dm {
	case "lightdm":
		return getLightdmAutoLoginUser(lightdmDefaultConfig)
	case "kdm":
		return getKDMAutoLoginUser(kdmDefaultConfig)
	case "gdm":
		return getGDMAutoLoginUser(gdmDefaultConfig)
	default:
		return "", fmt.Errorf("Not supported or invalid display manager: %q", dm)
	}

	return "", nil
}

//Default config: /etc/lightdm/lightdm.conf
func getLightdmAutoLoginUser(file string) (string, error) {
	if !dutils.IsFileExist(file) {
		return "", fmt.Errorf("Not found this file: %s", file)
	}

	v, exist := dutils.ReadKeyFromKeyFile(file,
		"SeatDefaults", "autologin-user", "")
	if !exist {
		return "", nil
	}

	name, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("The value's type error.")
	}

	return name, nil
}

func setLightdmAutoLoginUser(name, file string) error {
	success := dutils.WriteKeyToKeyFile(file,
		"SeatDefaults", "autologin-user", name)
	if !success {
		return fmt.Errorf("Set autologin user for %q failed!", name)
	}

	return nil
}

//Default config: /usr/share/config/kdm/kdmrc
func getKDMAutoLoginUser(file string) (string, error) {
	if !dutils.IsFileExist(file) {
		return "", fmt.Errorf("Not found this file: %s", file)
	}

	v, exist := dutils.ReadKeyFromKeyFile(file,
		"X-:0-Core", "AutoLoginEnable", true)
	if !exist {
		return "", nil
	}

	enable, ok := v.(bool)
	if !ok {
		return "", fmt.Errorf("The value's type error.")
	}

	if !enable {
		return "", nil
	}

	v, exist = dutils.ReadKeyFromKeyFile(file,
		"X-:0-Core", "AutoLoginUser", "")
	if !exist {
		return "", nil
	}

	var name string
	name, ok = v.(string)
	if !ok {
		return "", fmt.Errorf("The value's type error.")
	}

	return name, nil
}

func setKDMAutoLoginUser(name, file string) error {
	success := dutils.WriteKeyToKeyFile(file,
		"X-:0-Core", "AutoLoginEnable", true)
	if !success {
		return fmt.Errorf("Set 'AutoLoginEnable' to 'true' failed!")
	}

	success = dutils.WriteKeyToKeyFile(file,
		"X-:0-Core", "AutoLoginUser", name)
	if !success {
		return fmt.Errorf("Set autologin user for %q failed!", name)
	}

	return nil
}

//Default config: /etc/gdm/custom.conf
func getGDMAutoLoginUser(file string) (string, error) {
	if !dutils.IsFileExist(file) {
		return "", fmt.Errorf("Not found this file: %s", file)
	}

	v, exist := dutils.ReadKeyFromKeyFile(file,
		"daemon", "AutomaticLogin", "")
	if !exist {
		return "", nil
	}

	name, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("The value's type error.")
	}

	return name, nil
}

func setGDMAutoLoginUser(name, file string) error {
	success := dutils.WriteKeyToKeyFile(file,
		"daemon", "AutomaticLogin", name)
	if !success {
		return fmt.Errorf("Set autologin user for %q failed!", name)
	}

	return nil
}

//Default config: /etc/X11/default-display-manager
func getDefaultDisplayManager(file string) (string, error) {
	if !dutils.IsFileExist(file) {
		return "", fmt.Errorf("Not found this file: %s", file)
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	var tmp string
	for _, b := range content {
		if b == '\n' {
			continue
		}

		tmp += string(b)
	}

	return path.Base(tmp), nil
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
