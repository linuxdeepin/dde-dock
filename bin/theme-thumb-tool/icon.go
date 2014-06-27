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
	"os"
	"path"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"regexp"
	"strings"
)

type iconTypeInfo struct {
	appsList    []string
	placesList  []string
	devicesList []string
	actionsList []string
	statusList  []string
	emblemsList []string
	typeCnt     int
}

const (
	ICON_CONFIG = "index.theme"

	NAME_APP_DIR     = "apps"
	NAME_PLACE_DIR   = "places"
	NAME_DEVICES_DIR = "devices"
	NAME_ACTIONS_DIR = "actions"
	NAME_STATUS_DIR  = "status"
	NAME_EMBLEMS_DIR = "emblems"

	FOLDER          = "folder.png"
	USER_TRASH      = "user-trash.png"
	USER_TRASH_FULL = "user-trash-full.png"
	FILE_MANAGER    = "system-file-manager.png"
)

func getIconDirList(iconPath string) []string {
	list := []string{}
	if !dutils.IsFileExist(iconPath) {
		return list
	}

	filename := path.Join(iconPath, ICON_CONFIG)
	dirList, ok := dutils.ReadKeyFromKeyFile(filename,
		"Icon Theme", "Directories", list)
	if !ok {
		return list
	}

	var tList []string
	if tList, ok = dirList.([]string); !ok {
		return list
	}

	for _, n := range tList {
		strs := strings.Split(n, ",")
		for _, t := range strs {
			if strings.Contains(t, "scalable") {
				continue
			}

			name := path.Join(iconPath, t)
			if !dutils.IsFileExist(name) {
				continue
			}

			list = append(list, name)
		}
	}

	return list
}

func getIconTypeInfo(iconPath string) iconTypeInfo {
	typeInfo := iconTypeInfo{}
	if !dutils.IsFileExist(iconPath) {
		return typeInfo
	}

	dirList := getIconDirList(iconPath)
	if len(dirList) < 1 {
		return typeInfo
	}

	for _, l := range dirList {
		if strings.Contains(l, NAME_APP_DIR) {
			if len(typeInfo.appsList) < 1 {
				typeInfo.typeCnt++
			}

			typeInfo.appsList = append(typeInfo.appsList, l)
		} else if strings.Contains(l, NAME_PLACE_DIR) {
			if len(typeInfo.placesList) < 1 {
				typeInfo.typeCnt++
			}

			typeInfo.placesList = append(typeInfo.placesList, l)
		} else if strings.Contains(l, NAME_DEVICES_DIR) {
			if len(typeInfo.devicesList) < 1 {
				typeInfo.typeCnt++
			}

			typeInfo.devicesList = append(typeInfo.devicesList, l)
		} else if strings.Contains(l, NAME_ACTIONS_DIR) {
			if len(typeInfo.actionsList) < 1 {
				typeInfo.typeCnt++
			}

			typeInfo.actionsList = append(typeInfo.actionsList, l)
		} else if strings.Contains(l, NAME_STATUS_DIR) {
			if len(typeInfo.statusList) < 1 {
				typeInfo.typeCnt++
			}

			typeInfo.statusList = append(typeInfo.statusList, l)
		} else if strings.Contains(l, NAME_EMBLEMS_DIR) {
			if len(typeInfo.emblemsList) < 1 {
				typeInfo.typeCnt++
			}

			typeInfo.emblemsList = append(typeInfo.emblemsList, l)
		}
	}

	return typeInfo
}

func getPngList(dir string) []string {
	list := []string{}

	if !dutils.IsFileExist(dir) {
		return list
	}

	fp, err := os.Open(dir)
	if err != nil {
		return list
	}
	defer fp.Close()

	names, err1 := fp.Readdirnames(0)
	if err1 != nil {
		return list
	}

	for _, n := range names {
		ok, _ := regexp.MatchString(`\.png$`, n)
		if !ok {
			continue
		}

		filename := path.Join(dir, n)
		list = append(list, filename)
	}

	return list
}

func getAppsFiles(list []string) (aFile string) {
	if d, ok := isStrInList("48", list); ok {
		filename := path.Join(d, FILE_MANAGER)
		if dutils.IsFileExist(filename) {
			aFile = filename
			return
		}
	}

	if d, ok := isStrInList("32", list); ok {
		filename := path.Join(d, FILE_MANAGER)
		if dutils.IsFileExist(filename) {
			aFile = filename
			return
		}
	}

	if d, ok := isStrInList("24", list); ok {
		filename := path.Join(d, FILE_MANAGER)
		if dutils.IsFileExist(filename) {
			aFile = filename
			return
		}
	}

	if d, ok := isStrInList("22", list); ok {
		filename := path.Join(d, FILE_MANAGER)
		if dutils.IsFileExist(filename) {
			aFile = filename
			return
		}
	}

	imgList := getPngList(list[0])
	if len(imgList) > 0 {
		aFile = imgList[0]
	}
	return
}

func getPlacesFiles(list []string) (pFile1, pFile2 string) {
	flag1 := false
	flag2 := false

	if d, ok := isStrInList("48", list); ok {
		filename := path.Join(d, FOLDER)
		if dutils.IsFileExist(filename) {
			pFile1 = filename
			flag1 = true
		}

		filename = path.Join(d, USER_TRASH_FULL)
		if dutils.IsFileExist(filename) {
			pFile2 = filename
			flag2 = true
		} else {
			filename = path.Join(d, USER_TRASH)
			if dutils.IsFileExist(filename) {
				pFile2 = filename
				flag2 = true
			}
		}

		if flag1 && flag2 {
			return
		}
	}

	flag1 = false
	flag2 = false
	pFile1 = ""
	pFile2 = ""
	if d, ok := isStrInList("32", list); ok {
		filename := path.Join(d, FOLDER)
		if dutils.IsFileExist(filename) {
			pFile1 = filename
			flag1 = true
		}

		filename = path.Join(d, USER_TRASH_FULL)
		if dutils.IsFileExist(filename) {
			pFile2 = filename
			flag2 = true
		} else {
			filename = path.Join(d, USER_TRASH)
			if dutils.IsFileExist(filename) {
				pFile2 = filename
				flag2 = true
			}
		}

		if flag1 && flag2 {
			return
		}
	}

	flag1 = false
	flag2 = false
	pFile1 = ""
	pFile2 = ""
	if d, ok := isStrInList("24", list); ok {
		filename := path.Join(d, FOLDER)
		if dutils.IsFileExist(filename) {
			pFile1 = filename
			flag1 = true
		}

		filename = path.Join(d, USER_TRASH_FULL)
		if dutils.IsFileExist(filename) {
			pFile2 = filename
			flag2 = true
		} else {
			filename = path.Join(d, USER_TRASH)
			if dutils.IsFileExist(filename) {
				pFile2 = filename
				flag2 = true
			}
		}

		if flag1 && flag2 {
			return
		}
	}

	flag1 = false
	flag2 = false
	pFile1 = ""
	pFile2 = ""
	if d, ok := isStrInList("22", list); ok {
		filename := path.Join(d, FOLDER)
		if dutils.IsFileExist(filename) {
			pFile1 = filename
			flag1 = true
		}

		filename = path.Join(d, USER_TRASH_FULL)
		if dutils.IsFileExist(filename) {
			pFile2 = filename
			flag2 = true
		} else {
			filename = path.Join(d, USER_TRASH)
			if dutils.IsFileExist(filename) {
				pFile2 = filename
				flag2 = true
			}
		}

		if flag1 && flag2 {
			return
		}
	}

	imgList := getPngList(list[0])
	if len(imgList) > 1 {
		pFile1 = imgList[0]
		pFile2 = imgList[1]
	}

	return
}

func getTypePngFiles(list []string, num int) []string {
	if num < 1 {
		return []string{}
	}

	if d, ok := isStrInList("48", list); ok {
		imgList := getPngList(d)
		if len(imgList) >= num {
			return imgList[0:num]
		}
	}

	if d, ok := isStrInList("32", list); ok {
		imgList := getPngList(d)
		if len(imgList) >= num {
			return imgList[0:num]
		}
	}

	if d, ok := isStrInList("24", list); ok {
		imgList := getPngList(d)
		if len(imgList) >= num {
			return imgList[0:num]
		}
	}

	if d, ok := isStrInList("22", list); ok {
		imgList := getPngList(d)
		if len(imgList) >= num {
			return imgList[0:num]
		}
	}

	imgList := getPngList(list[0])
	if len(imgList) >= num {
		return imgList[0:num]
	}

	return []string{}
}

func getIconFiles(iconPath string) (f1, f2, f3 string) {
	if !dutils.IsFileExist(iconPath) {
		return
	}

	iconNum := 3
	iconInfo := getIconTypeInfo(iconPath)
	if len(iconInfo.appsList) > 0 {
		filename := getAppsFiles(iconInfo.appsList)
		if len(filename) > 0 {
			f1 = filename
			iconNum -= 1
		}
	}

	if len(iconInfo.placesList) > 0 {
		file1, file2 := getPlacesFiles(iconInfo.placesList)
		if len(file1) > 0 && len(file2) > 0 {
			f2 = file1
			f3 = file2
			iconNum -= 2
		}
	}

	if iconNum > 0 {
		if len(iconInfo.actionsList) > 0 {
			list := getTypePngFiles(iconInfo.actionsList, 3)
			if len(list) > 2 {
				f1 = list[0]
				f2 = list[1]
				f3 = list[2]
				return
			}
		}

		if len(iconInfo.devicesList) > 0 {
			list := getTypePngFiles(iconInfo.devicesList, 3)
			if len(list) > 2 {
				f1 = list[0]
				f2 = list[1]
				f3 = list[2]
				return
			}
		}

		if len(iconInfo.statusList) > 0 {
			list := getTypePngFiles(iconInfo.statusList, 3)
			if len(list) > 2 {
				f1 = list[0]
				f2 = list[1]
				f3 = list[2]
				return
			}
		}

		if len(iconInfo.emblemsList) > 0 {
			list := getTypePngFiles(iconInfo.emblemsList, 3)
			if len(list) > 2 {
				f1 = list[0]
				f2 = list[1]
				f3 = list[2]
				return
			}
		}
	}
	return
}
