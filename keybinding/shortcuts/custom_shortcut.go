/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

	"pkg.deepin.io/dde/daemon/keybinding/util"
	"pkg.deepin.io/lib/keyfile"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	kfKeyName       = "Name"
	kfKeyKeystrokes = "Accels"
	kfKeyAction     = "Action"
)

type CustomShortcut struct {
	BaseShortcut
	manager *CustomShortcutManager
	Cmd     string `json:"Exec"`
}

func (cs *CustomShortcut) Marshal() (string, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return util.MarshalJSON(cs)
}

func (cs *CustomShortcut) SaveKeystrokes() error {
	section := cs.GetId()
	csm := cs.manager
	csm.kfile.SetStringList(section, kfKeyKeystrokes, cs.getKeystrokesStrv())
	return csm.Save()
}

// 经过 Reset 重置后， 自定义快捷键的 keystrokes 被设置为空，始终返回 false
// 是为了另外计算改变的自定义快捷键项目。
func (cs *CustomShortcut) ReloadKeystrokes() bool {
	cs.setKeystrokes(nil)
	return false
}

func (cs *CustomShortcut) ShouldEmitSignalChanged() bool {
	return true
}

func (cs *CustomShortcut) Save() error {
	section := cs.GetId()
	kfile := cs.manager.kfile
	kfile.SetString(section, kfKeyName, cs.Name)
	kfile.SetString(section, kfKeyAction, cs.Cmd)
	kfile.SetStringList(section, kfKeyKeystrokes, cs.getKeystrokesStrv())
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

func (cs *CustomShortcut) SetName(name string) {
	cs.Name = name
	if cs.manager.pinyinEnabled {
		cs.nameBlocksInit = false
		cs.initNameBlocks()
	}
}

type CustomShortcutManager struct {
	file          string
	kfile         *keyfile.KeyFile
	pinyinEnabled bool
}

func NewCustomShortcutManager(file string) *CustomShortcutManager {
	kfile := keyfile.NewKeyFile()
	err := kfile.LoadFromFile(file)
	if err != nil && !os.IsNotExist(err) {
		logger.Warning(err)
	}

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
		keystrokes, _ := kfile.GetStringList(section, kfKeyKeystrokes)

		shortcut := &CustomShortcut{
			BaseShortcut: BaseShortcut{
				Id:         id,
				Type:       ShortcutTypeCustom,
				Keystrokes: ParseKeystrokes(keystrokes),
				Name:       name,
			},
			manager: csm,
			Cmd:     cmd,
		}

		ret = append(ret, shortcut)
	}
	return ret
}

func (csm *CustomShortcutManager) Save() error {
	err := os.MkdirAll(filepath.Dir(csm.file), 0755)
	if err != nil {
		return err
	}
	return csm.kfile.SaveToFile(csm.file)
}

func (csm *CustomShortcutManager) Add(name, action string, keystrokes []*Keystroke) (Shortcut, error) {
	id := dutils.GenUuid()
	csm.kfile.SetString(id, kfKeyName, name)
	csm.kfile.SetString(id, kfKeyAction, action)

	keystrokesStrv := make([]string, 0, len(keystrokes))
	for _, ks := range keystrokes {
		keystrokesStrv = append(keystrokesStrv, ks.String())
	}
	csm.kfile.SetStringList(id, kfKeyKeystrokes, keystrokesStrv)

	shortcut := &CustomShortcut{
		BaseShortcut: BaseShortcut{
			Id:         id,
			Type:       ShortcutTypeCustom,
			Keystrokes: keystrokes,
			Name:       name,
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
