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

package touchpad

import (
	"pkg.deepin.io/lib/gio-2.0"
)

func (tpad *Touchpad) handleGSettings() {
	tpad.settings.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case tpadKeyEnabled:
			tpad.enable(tpad.TPadEnable.Get())
		case tpadKeyLeftHanded:
			tpad.leftHanded(tpad.LeftHanded.Get())
			tpad.tapToClick(tpad.TapClick.Get(),
				tpad.LeftHanded.Get())
		case tpadKeyTapClick:
			tpad.tapToClick(tpad.TapClick.Get(),
				tpad.LeftHanded.Get())
		case tpadKeyNaturalScroll, tpadKeyScrollDelta:
			tpad.naturalScroll(tpad.NaturalScroll.Get(),
				tpad.DeltaScroll.Get())
		case tpadKeyEdgeScroll:
			tpad.edgeScroll(tpad.EdgeScroll.Get())
		case tpadKeyVertScroll, tpadKeyHorizScroll:
			tpad.twoFingerScroll(tpad.VertScroll.Get(),
				tpad.HorizScroll.Get())
		case tpadKeyWhileTyping:
			tpad.disableTpadWhileTyping(tpad.DisableIfTyping.Get())
		case tpadKeyDoubleClick:
			tpad.doubleClick(tpad.DoubleClick.Get())
		case tpadKeyDragThreshold:
			tpad.dragThreshold(tpad.DragThreshold.Get())
		case tpadKeyAcceleration:
			tpad.motionAcceleration(tpad.MotionAcceleration.Get())
		case tpadKeyThreshold:
			tpad.motionThreshold(tpad.MotionThreshold.Get())
		}
	})
}
