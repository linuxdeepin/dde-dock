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

package libtouchpad

import (
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

func (touchpad *Touchpad) handleGSettings() {
	touchpad.settings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case tpadKeyEnabled:
			touchpad.enable(touchpad.TPadEnable.Get())
		case tpadKeyLeftHanded:
			touchpad.leftHanded(touchpad.LeftHanded.Get())
			touchpad.tapToClick(touchpad.TapClick.Get(),
				touchpad.LeftHanded.Get())
		case tpadKeyTapClick:
			touchpad.tapToClick(touchpad.TapClick.Get(),
				touchpad.LeftHanded.Get())
		case tpadKeyNaturalScroll, tpadKeyScrollDelta:
			touchpad.naturalScroll(touchpad.NaturalScroll.Get(),
				touchpad.DeltaScroll.Get())
		case tpadKeyEdgeScroll:
			touchpad.edgeScroll(touchpad.EdgeScroll.Get())
		case tpadKeyVertScroll, tpadKeyHorizScroll:
			touchpad.twoFingerScroll(touchpad.VertScroll.Get(),
				touchpad.HorizScroll.Get())
		case tpadKeyWhileTyping:
			touchpad.disableTpadWhileTyping(touchpad.DisableIfTyping.Get())
		case tpadKeyDoubleClick:
			touchpad.doubleClick(touchpad.DoubleClick.Get())
		case tpadKeyDragThreshold:
			touchpad.dragThreshold(touchpad.DragThreshold.Get())
		case tpadKeyAcceleration:
			touchpad.motionAcceleration(touchpad.MotionAcceleration.Get())
		case tpadKeyThreshold:
			touchpad.motionThreshold(touchpad.MotionThreshold.Get())
		}
	})
}
