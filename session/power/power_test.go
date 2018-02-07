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

package power

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func _setWarnLevelConfigDefault(c *WarnLevelConfig) {
	c.UsePercentageForPolicy = false
	c.LowTime = 1200
	c.CriticalTime = 600
	c.ActionTime = 300

	c.LowPercentage = 20
	c.CriticalPercentage = 10
	c.ActionPercentage = 5
}

func TestWarnLevelConfig(t *testing.T) {
	Convey("WarnLevelConfig isValid", t, func() {
		c := NewWarnLevelConfig()
		_setWarnLevelConfigDefault(c)
		So(c.isValid(), ShouldBeTrue)

		c.LowTime = 599
		So(c.isValid(), ShouldBeFalse)

		_setWarnLevelConfigDefault(c)
		c.LowPercentage = 9
		So(c.isValid(), ShouldBeFalse)
	})
}

func Test_getWarnLevel(t *testing.T) {
	Convey("_getWarnLevel", t, func() {
		config := NewWarnLevelConfig()
		_setWarnLevelConfigDefault(config)
		onBattery := false
		So(_getWarnLevel(config, onBattery, 1.0, 0), ShouldEqual, WarnLevelNone)

		onBattery = true
		config.UsePercentageForPolicy = true
		So(_getWarnLevel(config, onBattery, 0.0, 0), ShouldEqual, WarnLevelNone)
		So(_getWarnLevel(config, onBattery, 1.0, 0), ShouldEqual, WarnLevelAction)
		So(_getWarnLevel(config, onBattery, 5.0, 0), ShouldEqual, WarnLevelAction)
		So(_getWarnLevel(config, onBattery, 5.1, 0), ShouldEqual, WarnLevelCritical)
		So(_getWarnLevel(config, onBattery, 10.0, 0), ShouldEqual, WarnLevelCritical)
		So(_getWarnLevel(config, onBattery, 10.1, 0), ShouldEqual, WarnLevelLow)
		So(_getWarnLevel(config, onBattery, 20.0, 0), ShouldEqual, WarnLevelLow)
		So(_getWarnLevel(config, onBattery, 20.1, 0), ShouldEqual, WarnLevelNone)
		So(_getWarnLevel(config, onBattery, 50.1, 0), ShouldEqual, WarnLevelNone)

		config.UsePercentageForPolicy = false
		// use time to empty
		So(_getWarnLevel(config, onBattery, 0, 0), ShouldEqual, WarnLevelNone)
		So(_getWarnLevel(config, onBattery, 0, 1), ShouldEqual, WarnLevelAction)
		So(_getWarnLevel(config, onBattery, 0, 300), ShouldEqual, WarnLevelAction)
		So(_getWarnLevel(config, onBattery, 0, 301), ShouldEqual, WarnLevelCritical)
		So(_getWarnLevel(config, onBattery, 0, 600), ShouldEqual, WarnLevelCritical)
		So(_getWarnLevel(config, onBattery, 0, 601), ShouldEqual, WarnLevelLow)
		So(_getWarnLevel(config, onBattery, 0, 1200), ShouldEqual, WarnLevelLow)
		So(_getWarnLevel(config, onBattery, 0, 1201), ShouldEqual, WarnLevelNone)
		So(_getWarnLevel(config, onBattery, 0, 12000), ShouldEqual, WarnLevelNone)
	})
}
