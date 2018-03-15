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

package inputdevices

import "pkg.deepin.io/lib/dbus1"

func (m *Mouse) Reset() *dbus.Error {
	for _, key := range m.setting.ListKeys() {
		m.setting.Reset(key)
	}
	return nil
}

func (tp *TrackPoint) Reset() *dbus.Error {
	for _, key := range tp.setting.ListKeys() {
		tp.setting.Reset(key)
	}
	return nil
}

func (tpad *Touchpad) Reset() *dbus.Error {
	for _, key := range tpad.setting.ListKeys() {
		tpad.setting.Reset(key)
	}
	return nil
}

func (w *Wacom) Reset() *dbus.Error {
	for _, key := range w.setting.ListKeys() {
		w.setting.Reset(key)
	}
	for _, key := range w.stylusSetting.ListKeys() {
		w.stylusSetting.Reset(key)
	}
	for _, key := range w.eraserSetting.ListKeys() {
		w.eraserSetting.Reset(key)
	}
	return nil
}

func (kbd *Keyboard) Reset() *dbus.Error {
	for _, key := range kbd.setting.ListKeys() {
		kbd.setting.Reset(key)
	}
	return nil
}

func (kbd *Keyboard) LayoutList() (map[string]string, *dbus.Error) {
	return kbd.layoutDescMap, nil
}

func (kbd *Keyboard) GetLayoutDesc(layout string) (string, *dbus.Error) {
	if len(layout) == 0 {
		return "", nil
	}

	desc, ok := kbd.layoutDescMap[layout]
	if !ok {
		return "", nil
	}

	return desc, nil
}

func (kbd *Keyboard) AddUserLayout(layout string) *dbus.Error {
	kbd.addUserLayout(layout)
	return nil
}

func (kbd *Keyboard) DeleteUserLayout(layout string) *dbus.Error {
	kbd.delUserLayout(layout)
	return nil
}

func (kbd *Keyboard) AddLayoutOption(option string) *dbus.Error {
	kbd.addUserOption(option)
	return nil
}

func (kbd *Keyboard) DeleteLayoutOption(option string) *dbus.Error {
	kbd.delUserOption(option)
	return nil
}

func (kbd *Keyboard) ClearLayoutOption() *dbus.Error {
	kbd.UserOptionList.Set([]string{})
	return nil
}
