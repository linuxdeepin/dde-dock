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

package mime

import (
	"fmt"
	"os"
	"path"

	"pkg.deepin.io/gir/gio-2.0"
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
	Exec      string
	CanDelete bool

	fileName string
}

type AppInfos []*AppInfo

func GetDefaultAppInfo(mimeType string) (*AppInfo, error) {
	id, err := mime.GetDefaultApp(mimeType, false)
	if err != nil {
		return nil, err
	}

	info, err := newAppInfoById2(id, mimeType)
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

func GetAppInfos(mimeType string) AppInfos {
	var infos AppInfos
	for _, id := range mime.GetAppList(mimeType) {
		appInfo, err := newAppInfoById2(id, mimeType)
		if err != nil {
			logger.Warning(err)
			continue
		}
		infos = append(infos, appInfo)
	}
	return infos
}

func newAppInfoByIdAux(id string, fn func(dai *gio.DesktopAppInfo, appInfo *AppInfo)) (*AppInfo, error) {
	dai := gio.NewDesktopAppInfo(id)
	if dai == nil {
		id = "kde4-" + id
		dai = gio.NewDesktopAppInfo(id)
	}
	if dai == nil {
		return nil, fmt.Errorf("gio.NewDesktopAppInfo failed: id %v", id)
	}
	defer dai.Unref()
	if !dai.ShouldShow() {
		return nil, fmt.Errorf("app %q should not show", id)
	}
	var appInfo = &AppInfo{
		Id:          id,
		Name:        dai.GetName(),
		DisplayName: dai.GetGenericName(),
		Description: dai.GetDescription(),
		Exec:        dai.GetCommandline(),
		fileName:    dai.GetFilename(),
	}
	iconObj := dai.GetIcon()
	if iconObj != nil {
		appInfo.Icon = iconObj.ToString()
		iconObj.Unref()
	}

	if fn != nil {
		fn(dai, appInfo)
	}

	return appInfo, nil
}

func newAppInfoById2(id string, mimeType string) (*AppInfo, error) {
	// 可以填写 CanDelete 字段
	gInfo, err := newAppInfoByIdAux(id, func(dai *gio.DesktopAppInfo, appInfo *AppInfo) {
		appInfo.CanDelete = canDeleteAssociation(dai, mimeType)
	})
	return gInfo, err
}

func newAppInfoById(id string) (*AppInfo, error) {
	return newAppInfoByIdAux(id, nil)
}

func canDeleteAssociation(appInfo *gio.DesktopAppInfo, mimeType string) bool {
	mimeTypes := appInfo.GetSupportedTypes()
	for _, mt := range mimeTypes {
		if mt == mimeType {
			return false
		}
	}
	return true
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
