/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package accounts

import (
	"fmt"
	"os/exec"
	"path"
	"pkg.deepin.io/dde/daemon/accounts/users"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
)

const (
	defaultLang       = "en_US"
	defaultLocaleFile = "/etc/default/locale"

	userDataCommon = "deepin-default-settings/skel.common"
	userDataLang   = "deepin-default-settings/skel.%s"
)

func (m *Manager) copyUserDatas(uPath string) {
	uid := getUidFromUserPath(uPath)
	info, err := users.GetUserInfoByUid(uid)
	if err != nil {
		logger.Warningf("Find user by uid '%s' failed: %v", uid, err)
		return
	}

	locale := getLocaleFromFile(defaultLocaleFile)
	lang := strings.Split(locale, ".")[0]
	if len(lang) == 0 {
		lang = defaultLang
	}

	err = copyCommonDatas(info.Home)
	if err != nil {
		logger.Debugf("Copy common datas for '%s' failed: %v",
			info.Name, err)
	}
	err = copyDatasByLang(info.Home, lang)
	if err != nil {
		logger.Debugf("Copy user datas for '%s' - '%s' failed: %v",
			info.Name, lang, err)
	}

	err = changeFileOwner(info.Home, info.Name, info.Name)
	if err != nil {
		logger.Warningf("Chown for '%s' failed: %v", info.Name, err)
	}
}

func copyCommonDatas(home string) error {
	data, err := findDatasPath(userDataCommon)
	if err != nil {
		return err
	}

	return dutils.CopyDir(data, home)
}

func copyDatasByLang(home, lang string) error {
	data, err := findDatasPath(fmt.Sprintf(userDataLang, lang))
	if err != nil {
		return err
	}

	return dutils.CopyDir(data, home)
}

func changeFileOwner(file, owner, group string) error {
	out, err := exec.Command("chown",
		"-hR",
		owner+":"+group,
		file).CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(out))
	}
	return nil
}

func findDatasPath(config string) (string, error) {
	data := path.Join("/usr/local/share", config)
	if dutils.IsFileExist(data) {
		return data, nil
	}

	data = path.Join("/usr/share", config)
	if dutils.IsFileExist(data) {
		return data, nil
	}

	return "", fmt.Errorf("Not found user datas '%s'", data)
}
