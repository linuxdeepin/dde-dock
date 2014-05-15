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

package main

import (
	"dlib/utils"
	"os"
	"path"
)

const (
	DMRC_FILE      = ".dmrc"
	DMRC_KEY_GROUP = "Desktop"
	PAM_ENV_FILE   = ".pam_environment"
)

func listenLocaleChange() {
	setDate.ConnectGenLocaleStatus(func(ok bool, locale string) {
		if ok && changeLocaleFlag {
			//setLocaleDmrc(locale)
			setLocalePamEnv(locale)
			changeLocaleFlag = false
		}
	})
}

func setLocaleDmrc(locale string) {
	homeDir, ok := getHomeDir()
	if !ok {
		return
	}
	utilsObj := utils.NewUtils()
	filePath := path.Join(homeDir, DMRC_FILE)
	utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
		"LANG", locale)
	utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
		"LANGUAGE", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_CTYPE", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_NUMERIC", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_TIME", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_COLLATE", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_MONETARY", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_MESSAGES", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_PAPER", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_NAME", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_ADDRESS", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_TELEPHONE", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_MEASUREMENT", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_IDENTIFICATION", locale)
	//utilsObj.WriteKeyToKeyFile(filePath, DMRC_KEY_GROUP,
	//"LC_ALL", "")
}

func genPamContents(locale string) string {
	contents := ""
	tmp := "LANG=" + locale + "\n"
	contents += tmp
	tmp = "LANGUAGE=" + locale + "\n"
	contents += tmp
	//tmp = "LC_CTYPE=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_NUMERIC=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_TIME=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_COLLATE=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_MONETARY=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_MESSAGES=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_PAPER=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_NAME=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_ADDRESS=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_TELEPHONE=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_MEASUREMENT=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_IDENTIFICATION=" + locale + "\n"
	//contents += tmp
	//tmp = "LC_ALL=\n"
	//contents += tmp

	return contents
}

func setLocalePamEnv(locale string) {
	homeDir, ok := getHomeDir()
	if !ok {
		return
	}
	filePath := path.Join(homeDir, PAM_ENV_FILE)

	fp, err := os.Create(filePath + "~")
	if err != nil {
		logger.Errorf("Create '%s' failed: %v", filePath+"~", err)
		return
	}
	defer fp.Close()

	if _, err = fp.WriteString(genPamContents(locale)); err != nil {
		logger.Errorf("Write '%s' failed: %v", filePath+"~", err)
		return
	}
	fp.Sync()
	os.Rename(filePath+"~", filePath)
}
