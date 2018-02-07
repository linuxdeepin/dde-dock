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

func (m *Mouse) Reset() {
	for _, key := range m.setting.ListKeys() {
		m.setting.Reset(key)
	}
}

func (tp *TrackPoint) Reset() {
	for _, key := range tp.setting.ListKeys() {
		tp.setting.Reset(key)
	}
}

func (tpad *Touchpad) Reset() {
	for _, key := range tpad.setting.ListKeys() {
		tpad.setting.Reset(key)
	}
}

func (w *Wacom) Reset() {
	for _, key := range w.setting.ListKeys() {
		w.setting.Reset(key)
	}
	for _, key := range w.stylusSetting.ListKeys() {
		w.stylusSetting.Reset(key)
	}
	for _, key := range w.eraserSetting.ListKeys() {
		w.eraserSetting.Reset(key)
	}
}

func (kbd *Keyboard) Reset() {
	for _, key := range kbd.setting.ListKeys() {
		kbd.setting.Reset(key)
	}
}

func (kbd *Keyboard) LayoutList() map[string]string {
	return kbd.layoutDescMap
}

func (kbd *Keyboard) GetLayoutDesc(layout string) string {
	if len(layout) == 0 {
		return ""
	}

	desc, ok := kbd.layoutDescMap[layout]
	if !ok {
		return ""
	}

	return desc
}

func (kbd *Keyboard) AddUserLayout(layout string) {
	kbd.addUserLayout(layout)
}

func (kbd *Keyboard) DeleteUserLayout(layout string) {
	kbd.delUserLayout(layout)
}

func (kbd *Keyboard) AddLayoutOption(option string) {
	kbd.addUserOption(option)
}

func (kbd *Keyboard) DeleteLayoutOption(option string) {
	kbd.delUserOption(option)
}

func (kbd *Keyboard) ClearLayoutOption() {
	kbd.UserOptionList.Set([]string{})
}
