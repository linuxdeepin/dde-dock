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

package utils

import (
	"fmt"
	"os"
	"path"
)

const (
	FileFlagSystemOwned = 0
	FileFlagUserOwned   = 1
)

var (
	errNotFound = fmt.Errorf("Not found in list")
)

type PathInfo struct {
	BaseName string
	FilePath string
	FileFlag int32
}

func IsNameInInfoList(name string, infos []PathInfo) bool {
	for _, info := range infos {
		if name == info.BaseName {
			return true
		}
	}

	return false
}

func GetInfoByName(name string, infos []PathInfo) (PathInfo, error) {
	for _, info := range infos {
		if info.BaseName == name {
			return info, nil
		}
	}

	return PathInfo{}, errNotFound
}

func GetBaseNameList(infos []PathInfo) []string {
	var list []string
	for _, info := range infos {
		list = append(list, info.BaseName)
	}

	return list
}

func GetFileFlagByName(name string, infos []PathInfo) int32 {
	for _, info := range infos {
		if info.BaseName == name {
			return info.FileFlag
		}
	}

	return -1
}

func GetInfoListFromDirs(dirInfos []PathInfo, conditions []string) []PathInfo {
	var infos []PathInfo

	for _, dirInfo := range dirInfos {
		tmpList, err := getInfoListFromDir(dirInfo, conditions)
		if err != nil {
			continue
		}
		infos = append(infos, tmpList...)
	}

	return infos
}

func getInfoListFromDir(dirInfo PathInfo, conditions []string) ([]PathInfo, error) {
	var infos []PathInfo

	fp, err := os.Open(dirInfo.FilePath)
	if err != nil {
		return infos, err
	}

	fileInfos, err := fp.Readdir(0)
	fp.Close()
	if err != nil {
		return infos, err
	}

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			continue
		}

		filename := path.Join(dirInfo.FilePath, fileInfo.Name())
		if !isMatchingConditions(filename, conditions) {
			continue
		}

		info := PathInfo{
			fileInfo.Name(),
			filename,
			dirInfo.FileFlag,
		}
		infos = append(infos, info)
	}

	return infos, nil
}

func isMatchingConditions(dir string, conditions []string) bool {
	fp, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer fp.Close()

	names, err := fp.Readdirnames(0)
	if err != nil {
		return false
	}

	var found int
	for _, name := range names {
		if isStrInList(name, conditions) {
			found += 1
		}
	}

	if found == len(conditions) {
		return true
	}

	return false
}

func isStrInList(value string, vList []string) bool {
	for _, v := range vList {
		if value == v {
			return true
		}
	}

	return false
}
