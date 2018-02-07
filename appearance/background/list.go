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
	dutils "pkg.deepin.io/lib/utils"
	"syscall"
)

var (
	home               = filepath.Clean(os.Getenv("HOME"))
	dirsCache          []string
	systemWallpaperDir = []string{
		"/usr/share/backgrounds",
		"/usr/share/wallpapers/deepin",
	}
)

// ListDirs list all background dirs
func ListDirs() []string {
	if len(dirsCache) != 0 {
		return dirsCache
	}

	userWallpapersDir := filepath.Join(getUserPictureDir(), "Wallpapers")
	var dirs = []string{userWallpapersDir}
	for _, dir := range systemWallpaperDir {
		dirs = append(dirs, dir)
	}

	dirsCache = dirs
	return dirsCache
}

func getBgFiles() []string {
	var walls []string
	for _, dir := range ListDirs() {
		walls = append(walls, getBgFilesInDir(dir)...)
	}
	return walls
}

const W_OK = 0x2

func isDeletable(file string) bool {
	dir := filepath.Dir(file)
	if strvContains(systemWallpaperDir, dir) {
		// directory is system wallpapers directory
		return false
	}
	if err := syscall.Access(dir, W_OK); err != nil {
		// directory is not writable
		return false
	}
	return true
}

func strvContains(strv []string, str string) bool {
	for _, v := range strv {
		if v == str {
			return true
		}
	}
	return false
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

func isFileInSpecialDir(file string, dirs []string) bool {
	file = dutils.DecodeURI(file)
	for _, dir := range dirs {
		if filepath.Dir(file) == dir {
			return true
		}
	}
	return false
}
