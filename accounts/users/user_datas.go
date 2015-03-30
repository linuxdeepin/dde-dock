/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
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
	"os"
	"path"
	"pkg.linuxdeepin.com/lib/archive"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"regexp"
	"strings"
)

const (
	defaultLangName = "en_US"
	defaultLangFile = "/etc/default/locale"
)

/**
 * Copy user resource datas to their home directory
 **/
func CopyUserDatas(user string) {
	info, err := GetUserInfoByName(user)
	if err != nil {
		return
	}
	lang := getDefaultLang()

	copyXDGDirConfig(info.Home, lang)
	renameXDGDirs(info.Home, lang)
	copyDeepinManuals(info.Home, lang)
	copySoundThemeData(info.Home, lang)
	copyBroswerConfig(info.Home, lang)
	changeDirOwner(user, info.Home)
}

func copyDeepinManuals(home, lang string) error {
	var (
		langDesc = map[string]string{
			"zh_CN": "用户手册",
		}

		langDoc = map[string]string{
			"zh_CN": "文档",
			"zh_TW": "文件",
			"en_US": "Documents",
		}
	)

	src := path.Join("/usr/share/doc/deepin-manuals", lang)
	if !dutils.IsFileExist(src) {
		return fmt.Errorf("Not found the file or directiry: %v", src)
	}

	destName, ok := langDesc[lang]
	if !ok {
		return fmt.Errorf("The language '%s' does not support", lang)
	}

	docName, ok := langDoc[lang]
	if !ok {
		docName = "Documents"
	}
	doc := path.Join(home, docName)
	if !dutils.IsFileExist(doc) {
		err := os.MkdirAll(doc, 0755)
		if err != nil {
			return err
		}
	}

	dest := path.Join(doc, destName)
	if dutils.IsFileExist(dest) {
		return nil
	}

	return dutils.SymlinkFile(src, dest)
}

func copySoundThemeData(home, lang string) error {
	src := "/usr/share/deepin-sample-music/playlist.m3u"
	if !dutils.IsFileExist(src) {
		return fmt.Errorf("Not found the file: %v", src)
	}

	dir := path.Join(home, ".sample-music")
	if !dutils.IsFileExist(dir) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	dest := path.Join(dir, "太歌·四季.m3u")
	if dutils.IsFileExist(dest) {
		return nil
	}

	return dutils.SymlinkFile(src, dest)
}

// Default broswer: google-chrome
func copyBroswerConfig(home, lang string) error {
	dest := path.Join(home, ".config/google-chrome")
	if dutils.IsFileExist(dest) {
		return nil
	}

	var (
		override   = "/usr/share/deepin-default-settings/google-chrome/override-chrome-config.tar"
		configLang = fmt.Sprintf("/usr/share/deepin-default-settings/google-chrome/chrome-config-%s.tar", lang)
		config     = "/usr/share/deepin-default-settings/google-chrome/chrome-config.tar"

		broswerConfig string
	)
	switch {
	case dutils.IsFileExist(override):
		broswerConfig = override
	case dutils.IsFileExist(configLang):
		broswerConfig = configLang
	case dutils.IsFileExist(config):
		broswerConfig = config
	}
	if len(broswerConfig) == 0 {
		return fmt.Errorf("Not found broswer configure file")
	}

	_, err := archive.Extracte(broswerConfig, path.Join(home, ".config"))
	return err
}

func renameXDGDirs(home, lang string) {
	var (
		desktop   = path.Join(home, "Desktop")
		templates = path.Join(home, "Templates")
	)

	switch lang {
	case "zh_CN":
		if dutils.IsFileExist(desktop) {
			os.Rename(desktop, path.Join(home, "桌面"))
		}

		if dutils.IsFileExist(templates) {
			os.Rename(templates, path.Join(home, "模板"))
			//dutils.CreateFile(path.Join(home, "模板", "文本文件"))
		}
	case "zh_TW":
		if dutils.IsFileExist(desktop) {
			os.Rename(desktop, path.Join(home, "桌面"))
		}

		if dutils.IsFileExist(templates) {
			os.Rename(templates, path.Join(home, "模板"))
			dutils.CreateFile(path.Join(home, "模板", "新增檔案"))
		}
	default:
		if dutils.IsFileExist(templates) {
			dutils.CreateFile(path.Join(templates, "New file"))
		}
	}
}

func copyXDGDirConfig(home, lang string) error {
	src := path.Join("/etc/skel.locale", lang, "user-dirs.dirs")
	if !dutils.IsFileExist(src) {
		return fmt.Errorf("Not found this file: %s", src)
	}

	dest := path.Join(home, ".config", "user-dirs.dirs")
	return dutils.CopyFile(src, dest)
}

func changeDirOwner(user, dir string) error {
	cmd := fmt.Sprintf("chown -hR %s:%s %s", user, user, dir)
	return doAction(cmd)
}

func getDefaultLang() string {
	fp, err := os.Open(defaultLangFile)
	if err != nil {
		return defaultLangName
	}
	defer fp.Close()

	var (
		locale  string
		match   = regexp.MustCompile(`^LANG=(.*)`)
		scanner = bufio.NewScanner(fp)
	)

	for scanner.Scan() {
		line := scanner.Text()
		fields := match.FindStringSubmatch(line)
		if len(fields) < 2 {
			continue
		}

		locale = fields[1]
		break
	}

	return strings.Split(locale, ".")[0]
}
