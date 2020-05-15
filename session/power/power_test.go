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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWarnLevelConfig(t *testing.T) {
	Convey("WarnLevelConfigManager isValid", t, func(c C) {
		conf := &warnLevelConfig{
			UsePercentageForPolicy: true,

			LowTime:      1200,
			DangerTime:   900,
			CriticalTime: 600,
			ActionTime:   300,

			LowPercentage:      20,
			DangerPercentage:   15,
			CriticalPercentage: 10,
			ActionPercentage:   5,
		}
		c.So(conf.isValid(), ShouldBeTrue)
		conf.LowTime = 599
		c.So(conf.isValid(), ShouldBeFalse)

		conf.LowTime = 1200
		conf.LowPercentage = 9
		c.So(conf.isValid(), ShouldBeFalse)
	})
}

func Test_getWarnLevel(t *testing.T) {
	Convey("_getWarnLevel", t, func(c C) {
		config := &warnLevelConfig{
			UsePercentageForPolicy: true,

			LowTime:      1200,
			DangerTime:   900,
			CriticalTime: 600,
			ActionTime:   300,

			LowPercentage:      20,
			DangerPercentage:   15,
			CriticalPercentage: 10,
			ActionPercentage:   5,
		}

		onBattery := false
		c.So(getWarnLevel(config, onBattery, 1.0, 0), ShouldEqual, WarnLevelNone)

		onBattery = true
		config.UsePercentageForPolicy = true

		c.So(getWarnLevel(config, onBattery, 0.0, 0), ShouldEqual, WarnLevelNone)
		c.So(getWarnLevel(config, onBattery, 1.1, 0), ShouldEqual, WarnLevelAction)
		c.So(getWarnLevel(config, onBattery, 5.0, 0), ShouldEqual, WarnLevelAction)
		c.So(getWarnLevel(config, onBattery, 5.1, 0), ShouldEqual, WarnLevelCritical)
		c.So(getWarnLevel(config, onBattery, 10.0, 0), ShouldEqual, WarnLevelCritical)
		c.So(getWarnLevel(config, onBattery, 10.1, 0), ShouldEqual, WarnLevelDanger)
		c.So(getWarnLevel(config, onBattery, 15.0, 0), ShouldEqual, WarnLevelDanger)
		c.So(getWarnLevel(config, onBattery, 15.1, 0), ShouldEqual, WarnLevelLow)
		c.So(getWarnLevel(config, onBattery, 20.0, 0), ShouldEqual, WarnLevelLow)
		c.So(getWarnLevel(config, onBattery, 20.1, 0), ShouldEqual, WarnLevelNone)
		c.So(getWarnLevel(config, onBattery, 50.0, 0), ShouldEqual, WarnLevelNone)

		config.UsePercentageForPolicy = false
		// use time to empty
		c.So(getWarnLevel(config, onBattery, 0, 0), ShouldEqual, WarnLevelNone)
		c.So(getWarnLevel(config, onBattery, 0, 61), ShouldEqual, WarnLevelAction)
		c.So(getWarnLevel(config, onBattery, 0, 300), ShouldEqual, WarnLevelAction)
		c.So(getWarnLevel(config, onBattery, 0, 301), ShouldEqual, WarnLevelCritical)
		c.So(getWarnLevel(config, onBattery, 0, 600), ShouldEqual, WarnLevelCritical)
		c.So(getWarnLevel(config, onBattery, 0, 601), ShouldEqual, WarnLevelDanger)
		c.So(getWarnLevel(config, onBattery, 0, 900), ShouldEqual, WarnLevelDanger)
		c.So(getWarnLevel(config, onBattery, 0, 901), ShouldEqual, WarnLevelLow)
		c.So(getWarnLevel(config, onBattery, 0, 1200), ShouldEqual, WarnLevelLow)
		c.So(getWarnLevel(config, onBattery, 0, 12001), ShouldEqual, WarnLevelNone)
	})
}

func TestMetaTasksMin(t *testing.T) {
	Convey("metaTasks.min", t, func(c C) {
		tasks := metaTasks{
			metaTask{
				name:  "n1",
				delay: 10,
			},
			metaTask{
				name:  "n2",
				delay: 30,
			},
			metaTask{
				name:  "n3",
				delay: 20,
			},
		}
		c.So(tasks.min(), ShouldEqual, 10)

		tasks = metaTasks{}
		c.So(tasks.min(), ShouldEqual, 0)

		tasks = metaTasks{
			metaTask{
				name:  "n1",
				delay: 10,
			},
		}
		c.So(tasks.min(), ShouldEqual, 10)
	})
}
