/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
	cmd     string
}

func (cs *CustomShortcut) SaveAccels() error {
	section := cs.GetId()
	accels := cs.GetAccels()
	csm := cs.manager
	accelStrv := make([]string, 0, len(accels))
	for _, accel := range accels {
		accelStrv = append(accelStrv, accel.String())
	}
	csm.kfile.SetStringList(section, kfKeyAccels, accelStrv)
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

func (cs *CustomShortcut) SetAction(newAction *Action) error {
	// TODO
	return nil
}

func (cs *CustomShortcut) SetName(newName string) error {
	return nil
}

func (cs *CustomShortcut) GetAction() *Action {
	a := &Action{
		Type: ActionTypeExecCmd,
		Arg: &ActionExecCmdArg{
			Cmd: cs.cmd,
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
			cmd:     cmd,
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
		cmd:     action,
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
