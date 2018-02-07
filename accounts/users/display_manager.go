/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
)

const (
	defaultDMFile         = "/etc/X11/default-display-manager"
	defaultDisplayService = "/etc/systemd/system/display-manager.service"
	lightdmConfig         = "/etc/lightdm/lightdm.conf"
	kdmConfig             = "/usr/share/config/kdm/kdmrc"
	gdmConfig             = "/etc/gdm/custom.conf"
	sddmConfig            = "/etc/sddm.conf"
	slimCOnfig            = "/etc/slim.conf"
	lxdmConfig            = "/etc/lxdm/lxdm.conf"

	kfGroupLightdmSeat        = "Seat:*"
	kfKeyLightdmAutoLoginUser = "autologin-user"
	kfKeyLightdmUserSession   = "user-session"
	kfGroupKDMXCore           = "X-:0-Core"
	kfKeyKDMAutoLoginEnable   = "AutoLoginEnable"
	kfKeyKDMAutoLoginUser     = "AutoLoginUser"
	kfGroupGDM3Daemon         = "daemon"
	kfKeyGDM3AutomaticEnable  = "AutomaticLoginEnable"
	kfKeyGDM3AutomaticLogin   = "AutomaticLogin"
	kfGroupDmrcDesktop        = "Desktop"
	kfKeyDmrcSession          = "Session"
	kfGroupSDDMAutologin      = "Autologin"
	kfKeySDDMUser             = "User"
	kfKeySDDMSession          = "Session"
	kfGroupLXDMBase           = "base"
	kfKeyLXDMAutologin        = "autologin"
	kfKeyLXDMSession          = "Session"

	// values: 'yes', 'no'
	slimKeyAutoLogin   = "auto_login"
	slimKeyDefaultUser = "default_user"
)

// SetAutoLoginUser set the autologin user,
// if disable autologin, set the 'username' to empty string
func SetAutoLoginUser(username, session string) error {
	dm, err := getDefaultDM(defaultDMFile)
	if err != nil {
		dm, err = getDMFromSystemService(defaultDisplayService)
		if err != nil {
			return err
		}
	}

	name, _ := GetAutoLoginUser()
	if name == username {
		return nil
	}

	// if user not in group 'autologin', lightdm autologin will no effect
	// detail see archlinux wiki for lightdm
	if username != "" {
		if !isGroupExists("autologin") {
			doAction("groupadd", []string{"-r", "autologin"})
		}

		if !isUserInGroup(username, "autologin") {
			err := doAction(userCmdGroup, []string{"-a", username, "autologin"})
			if err != nil {
				return err
			}
		}
	}

	switch dm {
	case "lightdm":
		keys := []string{kfKeyLightdmAutoLoginUser}
		values := []string{username}
		if session != "" {
			keys = append(keys, kfKeyLightdmUserSession)
			values = append(values, session)
		}
		return setIniKeys(lightdmConfig, kfGroupLightdmSeat,
			keys, values)
	case "kdm":
		values := []string{"true", username}
		if username == "" {
			values[0] = "false"
		}
		return setIniKeys(kdmConfig, kfGroupKDMXCore,
			[]string{
				kfKeyKDMAutoLoginEnable,
				kfKeyKDMAutoLoginUser}, values)
	case "gdm", "gdm3":
		values := []string{"True", username}
		if username == "" {
			values[0] = "False"
		}
		return setIniKeys(gdmConfig, kfGroupGDM3Daemon,
			[]string{
				kfKeyGDM3AutomaticEnable,
				kfKeyGDM3AutomaticLogin}, values)
	case "sddm":
		keys := []string{kfKeySDDMUser}
		values := []string{username}
		if session != "" {
			keys = append(keys, kfKeySDDMSession)
			values = append(values, session)
		}
		return setIniKeys(sddmConfig, kfGroupSDDMAutologin,
			keys, values)
	case "lxdm":
		keys := []string{kfKeySDDMUser}
		values := []string{username}
		// TODO: get session binary file
		// if session != "" {
		// 	keys = append(keys, kfKeySDDMSession)
		// 	values = append(values, session)
		// }
		return setIniKeys(lxdmConfig, kfGroupLXDMBase,
			keys, values)
	case "slim":
		// TODO
	}
	return fmt.Errorf("Not supported or invalid display manager: %q", dm)
}

// GetAutoLoginUser get the autologin user, if no, return empty string
func GetAutoLoginUser() (string, error) {
	dm, err := getDefaultDM(defaultDMFile)
	if err != nil {
		dm, err = getDMFromSystemService(defaultDisplayService)
		if err != nil {
			return "", err
		}
	}

	switch dm {
	case "lightdm":
		return getIniKeys(lightdmConfig, kfGroupLightdmSeat,
			[]string{kfKeyLightdmAutoLoginUser}, []string{""})
	case "kdm":
		return getIniKeys(kdmConfig, kfGroupKDMXCore,
			[]string{
				kfKeyKDMAutoLoginEnable,
				kfKeyKDMAutoLoginUser}, []string{"true", ""})
	case "gdm", "gdm3":
		return getIniKeys(gdmConfig, kfGroupGDM3Daemon,
			[]string{
				kfKeyGDM3AutomaticEnable,
				kfKeyGDM3AutomaticLogin}, []string{"True", ""})
	case "sddm":
		return getIniKeys(sddmConfig, kfGroupSDDMAutologin,
			[]string{kfKeySDDMUser}, []string{""})
	case "lxdm":
		return getIniKeys(lxdmConfig, kfGroupLXDMBase,
			[]string{kfKeyLXDMAutologin}, []string{""})
	case "slim":
		// TODO
	}
	return "", fmt.Errorf("Not supported or invalid display manager: %q", dm)
}

// GetDefaultXSession return the default user session
func GetDefaultXSession() (string, error) {
	dm, err := getDefaultDM(defaultDMFile)
	if err != nil {
		dm, err = getDMFromSystemService(defaultDisplayService)
		if err != nil {
			return "", err
		}
	}

	switch dm {
	case "lightdm":
		return getIniKeys(lightdmConfig, kfGroupLightdmSeat,
			[]string{"user-session"}, []string{""})
	case "kdm", "gdm", "gdm3":
		//return getIniKeys(userHome+"/.dmrc", kfGroupDmrcDesktop,
		//[]string{kfKeyDmrcSession}, []string{""})
		// no default session
		return "", nil
	case "sddm":
		v, err := getIniKeys(sddmConfig, kfGroupSDDMAutologin,
			[]string{kfKeySDDMSession}, []string{""})
		if err != nil {
			return "", err
		}
		return strings.TrimRight(v, ".desktop"), nil
	case "lxdm":
		// TODO: the session value is the binary file path
		// such as: session=/usr/bin/startlxde
		return "", nil
	case "slim":
		// no default session
		return "", nil
	}
	return "", fmt.Errorf("Not supported or invalid display manager: %q", dm)
}

// GetDMConfig return the current display manager
func GetDMConfig() (string, error) {
	dm, err := getDefaultDM(defaultDMFile)
	if err != nil {
		return "", err
	}

	switch dm {
	case "lightdm":
		return lightdmConfig, nil
	case "kdm":
		return kdmConfig, nil
	case "gdm", "gdm3":
		return gdmConfig, nil
	}
	return "", fmt.Errorf("Not supported the display manager: %q", dm)
}

func getIniKeys(filename, group string, keys, expected []string) (string, error) {
	if !dutils.IsFileExist(filename) {
		return "", fmt.Errorf("Not found the file: %s", filename)
	}

	kf, err := dutils.NewKeyFileFromFile(filename)
	if err != nil {
		return "", err
	}
	defer kf.Free()

	var username = ""
	for i := 0; i < len(keys); i++ {
		v, err := kf.GetString(group, keys[i])
		if err != nil {
			// ignore error, if no key exists
			return "", nil
		}

		if v == "" {
			return "", nil
		}

		if expected[i] != "" && v != expected[i] {
			return "", nil
		}
		username = v
	}
	return username, nil
}

func setIniKeys(filename, group string, keys, values []string) error {
	if !dutils.IsFileExist(filename) {
		err := dutils.CreateFile(filename)
		if err != nil {
			return err
		}
	}

	kf, err := dutils.NewKeyFileFromFile(filename)
	if err != nil {
		return err
	}
	defer kf.Free()

	for i := 0; i < len(keys); i++ {
		kf.SetString(group, keys[i], values[i])
	}

	_, err = kf.SaveToFile(filename)
	return err
}

//Default config: /etc/X11/default-display-manager
func getDefaultDM(file string) (string, error) {
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

// Default service: /etc/systemd/system/display-manager.service
func getDMFromSystemService(service string) (string, error) {
	if !dutils.IsFileExist(service) {
		return "", fmt.Errorf("Not found this file: %s", service)
	}

	name, err := os.Readlink(service)
	if err != nil {
		return "", err
	}

	base := path.Base(name)
	switch {
	case base == "lightdm.service":
		return "lightdm", nil
	case base == "gdm.service" || base == "gdm3.service":
		return "gdm", nil
	}
	return "", fmt.Errorf("Unsupported the login manager: %s", base)
}

// enable autologin: set 'auto_login' to 'yes', and 'default_user' to 'username'
func parseSlimConfig(filename, username string, isWirte bool) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	set := make(map[string]int)
	set[slimKeyAutoLogin] = -1
	set[slimKeyDefaultUser] = -1
	lines := strings.Split(string(content), "\n")
	for idx, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}
		line = strings.TrimSpace(line)
		list := strings.Split(line, " ")
		if list[0] == slimKeyAutoLogin {
			if isWirte {
				set[slimKeyAutoLogin] = idx
			} else {
				if len(list) < 2 || list[len(list)-1] != "yes" {
					return "", nil
				}
			}
			continue
		}

		if list[0] == slimKeyDefaultUser {
			if isWirte {
				line = slimKeyDefaultUser + " " + username
				set[slimKeyDefaultUser] = idx
			} else {
				if len(list) >= 2 {
					return list[len(list)-1], nil
				}
			}
		}
	}

	if !isWirte {
		return "", nil
	}

	autoLogin := ""
	defaultUser := ""
	sync := false
	idx, _ := set[slimKeyAutoLogin]
	if username != "" {
		autoLogin = slimKeyAutoLogin + " yes"
		defaultUser = slimKeyDefaultUser + " " + username
	}
	if idx == -1 && autoLogin != "" {
		lines = append(lines, autoLogin)
		sync = true
	} else {
		lines[idx] = autoLogin
		sync = true
	}

	idx, _ = set[slimKeyDefaultUser]
	if idx == -1 && defaultUser != "" {
		lines = append(lines, defaultUser)
		sync = true
	} else {
		lines[idx] = defaultUser
		sync = true
	}

	if !sync {
		return "", nil
	}

	data := strings.Join(lines, "\n")
	return "", ioutil.WriteFile(filename, []byte(data), 0644)
}
