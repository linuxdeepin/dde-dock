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

package screenedge

import (
	gio "gir/gio-2.0"
)

type Settings struct {
	gsettings *gio.Settings
}

func NewSettings() *Settings {
	s := new(Settings)
	s.gsettings = gio.NewSettings("com.deepin.dde.zone")
	return s
}

func (s *Settings) GetDelay() int32 {
	return s.gsettings.GetInt("delay")
}

func (s *Settings) SetEdgeAction(name, value string) {
	s.gsettings.SetString(name, value)
}

func (s *Settings) GetEdgeAction(name string) string {
	return s.gsettings.GetString(name)
}

func (s *Settings) GetWhiteList() []string {
	return s.gsettings.GetStrv("white-list")
}

func (s *Settings) GetBlackList() []string {
	return s.gsettings.GetStrv("black-list")
}
