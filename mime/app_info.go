/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package mime

import (
	"fmt"
	"gir/gio-2.0"
	"os"
	"path"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

type AppInfo struct {
	// Desktop id
	Id string
	// App name
	Name string
	// Display name
	DisplayName string
	// Comment
	Description string
	// Icon
	Icon string
	// Commandline
	Exec string

	fileName string
}

type AppInfos []*AppInfo

func GetAppInfo(ty string) (*AppInfo, error) {
	id, err := mime.GetDefaultApp(ty, false)
	if err != nil {
		return nil, err
	}

	info, err := newAppInfoById(id)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (infos AppInfos) Add(id string) (AppInfos, error) {
	for _, info := range infos {
		if info.Id == id {
			return infos, nil
		}
	}
	tmp, err := newAppInfoById(id)
	if err != nil {
		return nil, err
	}
	infos = append(infos, tmp)
	return infos, nil
}

func (infos AppInfos) Delete(id string) AppInfos {
	var ret AppInfos
	for _, info := range infos {
		if info.Id == id {
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func SetAppInfo(ty, id string) error {
	return mime.SetDefaultApp(ty, id)
}

func GetAppInfos(ty string) AppInfos {
	var infos AppInfos
	for _, id := range mime.GetAppList(ty) {
		appInfo, err := newAppInfoById(id)
		if err != nil {
			logger.Warning(err)
			continue
		}
		infos = append(infos, appInfo)
	}
	return infos
}

func newAppInfoById(id string) (*AppInfo, error) {
	ginfo := gio.NewDesktopAppInfo(id)
	if ginfo == nil {
		id = "kde4-" + id
		ginfo = gio.NewDesktopAppInfo(id)
	}
	if ginfo == nil {
		return nil, fmt.Errorf("gio.NewDesktopAppInfo failed: id %v", id)
	}

	defer ginfo.Unref()
	if !ginfo.ShouldShow() {
		return nil, fmt.Errorf("app %q should not show", id)
	}

	var info = &AppInfo{
		Id:          id,
		Name:        ginfo.GetName(),
		DisplayName: ginfo.GetGenericName(),
		Description: ginfo.GetDescription(),
		Exec:        ginfo.GetCommandline(),
		fileName:    ginfo.GetFilename(),
	}
	iconObj := ginfo.GetIcon()
	if iconObj != nil {
		info.Icon = iconObj.ToString()
		iconObj.Unref()
	}

	return info, nil
}

func findFilePath(file string) string {
	data := path.Join(os.Getenv("HOME"), ".local/share", file)
	if dutils.IsFileExist(data) {
		return data
	}

	data = path.Join("/usr/local/share", file)
	if dutils.IsFileExist(data) {
		return data
	}

	return path.Join("/usr/share", file)
}
