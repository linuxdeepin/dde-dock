/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package background

import (
	"os"
	"path/filepath"
)

var (
	systemWallpapersDir = []string{
		"/usr/share/wallpapers/deepin",
	}
)

// ListDirs list all background dirs
func ListDirs() []string {
	var result []string
	for _, value := range systemWallpapersDir {
		result = append(result, value)
	}
	result = append(result, CustomWallpapersConfigDir)
	return result
}

func getSysBgFiles() []string {
	var files []string
	for _, dir := range systemWallpapersDir {
		files = append(files, getBgFilesInDir(dir)...)
	}
	return files
}

func getCustomBgFiles() []string {
	return getBgFilesInDir(CustomWallpapersConfigDir)
}

func getBgFilesInDir(dir string) []string {
	fr, err := os.Open(dir)
	if err != nil {
		return []string{}
	}
	defer fr.Close()

	names, err := fr.Readdirnames(0)
	if err != nil {
		return []string{}
	}

	var walls []string
	for _, name := range names {
		path := filepath.Join(dir, name)
		if !IsBackgroundFile(path) {
			continue
		}
		walls = append(walls, path)
	}
	return walls
}

func isFileInDirs(file string, dirs []string) bool {
	for _, dir := range dirs {
		if filepath.Dir(file) == dir {
			return true
		}
	}
	return false
}
