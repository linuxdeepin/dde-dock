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
	"path"
	dutils "pkg.deepin.io/lib/utils"
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

	defaultDMFile         = "/etc/X11/default-display-manager"
	defaultDisplayService = "/etc/systemd/system/display-manager.service"
	lightdmConfig         = "/etc/lightdm/lightdm.conf"
	kdmConfig             = "/usr/share/config/kdm/kdmrc"
	gdmConfig             = "/etc/gdm/custom.conf"
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

func SetAutoLoginUser(username string) error {
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
		return setLightdmAutoLoginUser(username, lightdmConfig)
	case "kdm":
		return setKDMAutoLoginUser(username, kdmConfig)
	case "gdm":
		return setGDMAutoLoginUser(username, gdmConfig)
	default:
		return fmt.Errorf("Not supported or invalid display manager: %q", dm)
	}
}

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
		return getLightdmAutoLoginUser(lightdmConfig)
	case "kdm":
		return getKDMAutoLoginUser(kdmConfig)
	case "gdm":
		return getGDMAutoLoginUser(gdmConfig)
	default:
		return "", fmt.Errorf("Not supported or invalid display manager: %q", dm)
	}
}

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
	case "gdm":
		return gdmConfig, nil
	}
	return "", fmt.Errorf("Not supported the display manager: %q", dm)
}

//Default config: /etc/lightdm/lightdm.conf
func getLightdmAutoLoginUser(file string) (string, error) {
	if !dutils.IsFileExist(file) {
		return "", fmt.Errorf("Not found this file: %s", file)
	}

	v, exist := dutils.ReadKeyFromKeyFile(file,
		"Seat:*", "autologin-user", "")
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
		"Seat:*", "autologin-user", name)
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
