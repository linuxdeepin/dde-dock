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

package shortcuts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"pkg.deepin.io/lib/glib-2.0"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	// $USER_CONFIG/customKeyConfig
	customKeyConfig = "deepin/dde-daemon/keybinding/custom.ini"

	kfKeyName   = "Name"
	kfKeyAccels = "Accels"
	kfKeyAction = "Action"
)

type customKeyInfo struct {
	Shortcut
	Action string
}
type customKeyInfos []*customKeyInfo

var (
	kfileLocker sync.Mutex
)

func ListCustomKey() customKeyInfos {
	file, _ := getCustomConfig()
	return readCustomConfig(file)
}

func AddCustomKey(name, action string, accels []string) (string, error) {
	id := dutils.GenUuid()
	err := writeCustomKeyInfo(&customKeyInfo{
		Shortcut: Shortcut{
			Id:     id,
			Name:   name,
			Accels: filterNilStr(accels),
			Type:   KeyTypeCustom,
		},
		Action: action,
	})
	if err != nil {
		return "", err
	}

	return id, nil
}

func DeleteCustomKey(id string) error {
	file, err := getCustomConfig()
	if err != nil {
		return err
	}

	return doDelKeyInfo(id, file)
}

func (infos customKeyInfos) GetShortcuts() Shortcuts {
	var ss Shortcuts
	for _, info := range infos {
		ss = append(ss, &info.Shortcut)
	}
	return ss
}

func (infos customKeyInfos) Get(id string) *customKeyInfo {
	for _, info := range infos {
		if info.Id == id {
			return info
		}
	}
	return nil
}

func resetCustomKeys() {
	for _, info := range ListCustomKey() {
		DeleteCustomKey(info.Id)
	}
}

func disableCustomKey(id string) error {
	infos := ListCustomKey()
	if infos == nil {
		return fmt.Errorf("Invalid custom id '%s'", id)
	}

	info := infos.Get(id)
	if info == nil {
		return fmt.Errorf("Invalid custom id '%s'", id)
	}
	info.Accels = nil

	return writeCustomKeyInfo(info)
}

func modifyCustomName(id, value string) error {
	return setCustomKey(id, kfKeyName, value, false)
}

func modifyCustomAction(id, value string) error {
	return setCustomKey(id, kfKeyAction, value, false)
}

func modifyCustomAccels(id, value string, deleted bool) error {
	if len(value) == 0 {
		return nil
	}
	return setCustomKey(id, kfKeyAccels, value, deleted)
}

func setCustomKey(id, prop, value string, deleted bool) error {
	info := ListCustomKey().Get(id)
	if info == nil {
		return fmt.Errorf("Invalid custom id '%s'", id)
	}

	switch prop {
	case kfKeyName:
		info.Name = value
	case kfKeyAction:
		info.Action = value
	case kfKeyAccels:
		var ret bool
		if deleted {
			info.Accels, ret = delAccelFromList(value,
				info.Accels)
		} else {
			info.Accels, ret = addAccelToList(value,
				info.Accels)
		}
		if !ret {
			return nil
		}
	}
	return writeCustomKeyInfo(info)
}

func readCustomConfig(file string) customKeyInfos {
	kfile, err := dutils.NewKeyFileFromFile(file)
	if err != nil {
		return nil
	}
	defer kfile.Free()

	_, groups := kfile.GetGroups()
	var infos customKeyInfos
	for _, group := range groups {
		info, err := newKeyInfoByGroup(kfile, group)
		if err != nil {
			continue
		}
		infos = append(infos, info)
	}
	return infos
}

func writeCustomConfig(kfile *glib.KeyFile) error {
	file, err := getCustomConfig()
	if err != nil {
		return err
	}

	return saveKeyFile(kfile, file)
}

func writeCustomKeyInfo(info *customKeyInfo) error {
	file, err := getCustomConfig()
	if err != nil {
		return err
	}

	kfile, err := dutils.NewKeyFileFromFile(file)
	if err != nil {
		return err
	}
	defer kfile.Free()

	kfile.SetString(info.Id, kfKeyName, info.Name)
	kfile.SetString(info.Id, kfKeyAction, info.Action)
	kfile.SetStringList(info.Id, kfKeyAccels, info.Accels)

	return saveKeyFile(kfile, file)
}

func doDelKeyInfo(id, file string) error {
	kfile, err := dutils.NewKeyFileFromFile(file)
	if err != nil {
		return err
	}
	defer kfile.Free()

	err = delKeyFileGroup(kfile, id)
	if err != nil {
		return err
	}

	return saveKeyFile(kfile, file)
}

func newKeyInfoByGroup(kfile *glib.KeyFile, group string) (*customKeyInfo, error) {
	var (
		info customKeyInfo
		core Shortcut
		err  error
	)

	core.Id = group
	core.Type = KeyTypeCustom

	core.Name, err = kfile.GetString(group, kfKeyName)
	if err != nil {
		return nil, err
	}
	_, core.Accels, err = kfile.GetStringList(group, kfKeyAccels)
	if err != nil {
		return nil, err
	}
	info.Action, err = kfile.GetString(group, kfKeyAction)
	if err != nil {
		return nil, err
	}
	info.Shortcut = core

	return &info, nil
}

func delKeyFileGroup(kfile *glib.KeyFile, group string) error {
	_, groups := kfile.GetGroups()
	if !isStrInList(group, groups) {
		return nil
	}

	_, err := kfile.RemoveGroup(group)
	return err
}

func saveKeyFile(kfile *glib.KeyFile, file string) error {
	kfileLocker.Lock()
	defer kfileLocker.Unlock()

	_, content, err := kfile.ToData()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, []byte(content), 0644)
}

func getCustomConfig() (string, error) {
	file := path.Join(glib.GetUserConfigDir(), customKeyConfig)
	err := os.MkdirAll(path.Dir(file), 0755)
	if err != nil {
		return "", err
	}

	if !dutils.IsFileExist(file) {
		err := dutils.CreateFile(file)
		if err != nil {
			return "", err
		}
	}
	return file, nil
}
