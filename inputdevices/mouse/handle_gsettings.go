/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

package mouse

import (
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

func (m *Mouse) handleGSettings() {
	m.settings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case mouseKeyLeftHanded:
			m.leftHanded(m.LeftHanded.Get())
		case mouseKeyDisableTouchpad:
			m.disableTouchpad(m.DisableTpad.Get())
		case mouseKeyNaturalScroll:
			m.naturalScroll(m.NaturalScroll.Get())
		case mouseKeyMiddleButton:
			m.middleButtonEmulation(m.MiddleButtonEmulation.Get())
		case mouseKeyAcceleration:
			m.motionAcceleration(m.MotionAcceleration.Get())
		case mouseKeyThreshold:
			m.motionThreshold(m.MotionThreshold.Get())
		case mouseKeyDoubleClick:
			m.doubleClick(uint32(m.DoubleClick.Get()))
		case mouseKeyDragThreshold:
			m.dragThreshold(uint32(m.DragThreshold.Get()))
		}
	})
}
