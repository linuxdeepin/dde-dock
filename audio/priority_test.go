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

package audio

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_contains(t *testing.T) {
	Convey("contains", t, func(c C) {
		c.So(contains("hbc.abcd.1234", "world.abcd.1234", "hbc"), ShouldBeTrue)
		c.So(contains("hello.abcd.1234", "hbc.abcd.1234", "hbc"), ShouldBeTrue)
		c.So(contains("HBC.abcd.1234", "world.abcd.1234", "hbc"), ShouldBeTrue)
		c.So(contains("hello.abcd.1234", "HBC.abcd.1234", "hbc"), ShouldBeTrue)
		c.So(contains("hello.abcd.1234", "world.abcd.1234", "hbc"), ShouldBeFalse)
	})
}

func Test_GetPortType(t *testing.T) {
	Convey("GetPortType", t, func(c C) {
		c.So(GetPortType("hbc.abcd.1234", "world.abcd.1234"), ShouldEqual, PortTypeHeadset)
		c.So(GetPortType("bluez.abcd.1234", "world.abcd.1234"), ShouldEqual, PortTypeBluetooth)
		c.So(GetPortType("hbc.abcd.1234", "bluez.abcd.1234"), ShouldEqual, PortTypeBluetooth)
		c.So(GetPortType("usb.abcd.1234", "world.abcd.1234"), ShouldEqual, PortTypeHeadset)
		c.So(GetPortType("hbc.abcd.1234", "usb.abcd.1234"), ShouldEqual, PortTypeHeadset)
		c.So(GetPortType("hello.abcd.speaker", "world.abcd.1234"), ShouldEqual, PortTypeSpeaker)
		c.So(GetPortType("hdmi.abcd.speaker", "world.abcd.1234"), ShouldEqual, PortTypeHdmi)
	})
}

func Test_IsInputTypeAfter(t *testing.T) {
	Convey("IsInputTypeAfter", t, func(c C) {
		pr := NewPriorities()
		pr.defaultInit(CardList{})
		c.So(pr.IsInputTypeAfter(PortTypeHeadset, PortTypeBluetooth), ShouldBeFalse)
		c.So(pr.IsInputTypeAfter(PortTypeSpeaker, PortTypeBluetooth), ShouldBeFalse)
		c.So(pr.IsInputTypeAfter(PortTypeHdmi, PortTypeBluetooth), ShouldBeFalse)
		c.So(pr.IsInputTypeAfter(PortTypeSpeaker, PortTypeHeadset), ShouldBeFalse)
		c.So(pr.IsInputTypeAfter(PortTypeHdmi, PortTypeSpeaker), ShouldBeFalse)

		c.So(pr.IsInputTypeAfter(PortTypeBluetooth, PortTypeHeadset), ShouldBeTrue)
		c.So(pr.IsInputTypeAfter(PortTypeBluetooth, PortTypeSpeaker), ShouldBeTrue)
		c.So(pr.IsInputTypeAfter(PortTypeBluetooth, PortTypeHdmi), ShouldBeTrue)
		c.So(pr.IsInputTypeAfter(PortTypeHeadset, PortTypeSpeaker), ShouldBeTrue)
		c.So(pr.IsInputTypeAfter(PortTypeSpeaker, PortTypeHdmi), ShouldBeTrue)
	})
}

func Test_IsOutputTypeAfter(t *testing.T) {
	Convey("IsOutputTypeAfter", t, func(c C) {
		pr := NewPriorities()
		pr.defaultInit(CardList{})
		c.So(pr.IsOutputTypeAfter(PortTypeHeadset, PortTypeBluetooth), ShouldBeFalse)
		c.So(pr.IsOutputTypeAfter(PortTypeSpeaker, PortTypeBluetooth), ShouldBeFalse)
		c.So(pr.IsOutputTypeAfter(PortTypeHdmi, PortTypeBluetooth), ShouldBeFalse)
		c.So(pr.IsOutputTypeAfter(PortTypeSpeaker, PortTypeHeadset), ShouldBeFalse)
		c.So(pr.IsOutputTypeAfter(PortTypeHdmi, PortTypeSpeaker), ShouldBeFalse)

		c.So(pr.IsOutputTypeAfter(PortTypeBluetooth, PortTypeHeadset), ShouldBeTrue)
		c.So(pr.IsOutputTypeAfter(PortTypeBluetooth, PortTypeSpeaker), ShouldBeTrue)
		c.So(pr.IsOutputTypeAfter(PortTypeBluetooth, PortTypeHdmi), ShouldBeTrue)
		c.So(pr.IsOutputTypeAfter(PortTypeHeadset, PortTypeSpeaker), ShouldBeTrue)
		c.So(pr.IsOutputTypeAfter(PortTypeSpeaker, PortTypeHdmi), ShouldBeTrue)
	})
}
