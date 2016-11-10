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
	"gir/gio-2.0"
)

type GSettingsShortcut struct {
	BaseShortcut
	gsettings *gio.Settings
}

func NewGSettingsShortcut(gsettings *gio.Settings, id string, _type int32,
	accels []string, name string) *GSettingsShortcut {
	gs := &GSettingsShortcut{
		BaseShortcut: BaseShortcut{
			Id:     id,
			Type:   _type,
			Accels: ParseStandardAccels(accels),
			Name:   name,
		},
		gsettings: gsettings,
	}

	return gs
}

func (gs *GSettingsShortcut) SaveAccels() error {
	accelStrv := make([]string, 0, len(gs.Accels))
	for _, pa := range gs.Accels {
		accelStrv = append(accelStrv, pa.String())
	}
	gs.gsettings.SetStrv(gs.Id, accelStrv)
	logger.Debugf("GSettingsShortcut.SaveAccels id: %v, accels: %v", gs.Id, accelStrv)
	return nil
}

func parsedAccelsEqual(s1 []ParsedAccel, s2 []ParsedAccel) bool {
	l1 := len(s1)
	l2 := len(s2)

	if l1 != l2 {
		return false
	}

	for i := 0; i < l1; i++ {
		pa1 := s1[i]
		pa2 := s2[i]

		if pa1 != pa2 {
			return false
		}
	}
	return true
}

func (gs *GSettingsShortcut) ReloadAccels() bool {
	oldAccels := gs.GetAccels()
	id := gs.GetId()
	accelStrv := gs.gsettings.GetStrv(id)
	newAccels := ParseStandardAccels(accelStrv)
	gs.setAccels(newAccels)
	return !parsedAccelsEqual(oldAccels, newAccels)
}
