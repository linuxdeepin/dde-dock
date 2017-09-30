/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

package shortcuts

import (
	"os"
	"path/filepath"

	"pkg.deepin.io/lib/keyfile"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	kfKeyName   = "Name"
	kfKeyAccels = "Accels"
	kfKeyAction = "Action"
)

type CustomShortcut struct {
	BaseShortcut
	manager *CustomShortcutManager
	Cmd     string `json:"Exec"`
}

func (cs *CustomShortcut) getAccelStrv() []string {
	accels := cs.GetAccels()
	strv := make([]string, len(accels))
	for i, accel := range accels {
		strv[i] = accel.String()
	}
	return strv
}

func (cs *CustomShortcut) SaveAccels() error {
	section := cs.GetId()
	csm := cs.manager
	csm.kfile.SetStringList(section, kfKeyAccels, cs.getAccelStrv())
	return csm.Save()
}

// after Reset, custom shortcut accels should be empty
func (cs *CustomShortcut) ReloadAccels() bool {
	oldAccels := cs.GetAccels()
	cs.setAccels(nil)

	if len(oldAccels) > 0 {
		return true
	}
	return false
}

func (cs *CustomShortcut) Save() error {
	section := cs.GetId()
	kfile := cs.manager.kfile
	kfile.SetString(section, kfKeyName, cs.Name)
	kfile.SetString(section, kfKeyAction, cs.Cmd)
	kfile.SetStringList(section, kfKeyAccels, cs.getAccelStrv())
	return cs.manager.Save()
}

func (cs *CustomShortcut) GetAction() *Action {
	a := &Action{
		Type: ActionTypeExecCmd,
		Arg: &ActionExecCmdArg{
			Cmd: cs.Cmd,
		},
	}
	return a
}

type CustomShortcutManager struct {
	file  string
	kfile *keyfile.KeyFile
}

func NewCustomShortcutManager(file string) *CustomShortcutManager {
	kfile := keyfile.NewKeyFile()
	kfile.LoadFromFile(file)

	m := &CustomShortcutManager{
		file:  file,
		kfile: kfile,
	}
	return m
}

func (csm *CustomShortcutManager) List() []Shortcut {
	kfile := csm.kfile
	sections := kfile.GetSections()
	ret := make([]Shortcut, 0, len(sections))
	for _, section := range sections {
		id := section
		name, _ := kfile.GetString(section, kfKeyName)
		cmd, _ := kfile.GetString(section, kfKeyAction)
		accels, _ := kfile.GetStringList(section, kfKeyAccels)

		shortcut := &CustomShortcut{
			BaseShortcut: BaseShortcut{
				Id:     id,
				Type:   ShortcutTypeCustom,
				Accels: ParseStandardAccels(accels),
				Name:   name,
			},
			manager: csm,
			Cmd:     cmd,
		}

		ret = append(ret, shortcut)
	}
	return ret
}

func (csm *CustomShortcutManager) Save() error {
	os.MkdirAll(filepath.Dir(csm.file), 0755)
	return csm.kfile.SaveToFile(csm.file)
}

func (csm *CustomShortcutManager) Add(name, action string, accels []ParsedAccel) (Shortcut, error) {
	id := dutils.GenUuid()
	csm.kfile.SetString(id, kfKeyName, name)
	csm.kfile.SetString(id, kfKeyAction, action)
	// accels
	accelStrv := make([]string, 0, len(accels))
	for _, accel := range accels {
		accelStrv = append(accelStrv, accel.String())
	}
	csm.kfile.SetStringList(id, kfKeyAccels, accelStrv)

	shortcut := &CustomShortcut{
		BaseShortcut: BaseShortcut{
			Id:     id,
			Type:   ShortcutTypeCustom,
			Accels: accels,
			Name:   name,
		},
		manager: csm,
		Cmd:     action,
	}
	return shortcut, csm.Save()
}

func (csm *CustomShortcutManager) Delete(id string) error {
	if _, err := csm.kfile.GetSection(id); err != nil {
		return err
	}

	csm.kfile.DeleteSection(id)
	return csm.Save()
}

func (csm *CustomShortcutManager) DisableAll() error {
	kfile := csm.kfile
	sections := kfile.GetSections()
	for _, section := range sections {
		kfile.SetValue(section, kfKeyAccels, "")
	}
	return csm.Save()
}
