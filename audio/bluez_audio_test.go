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
	"pkg.deepin.io/lib/pulse"
)

func Test_isBluezAudio(t *testing.T) {
	Convey("isBluezAudio", t, func(c C) {
		c.So(isBluezAudio("bluez.abcd.1234"), ShouldBeTrue)
		c.So(isBluezAudio("hbc.bluez.1234"), ShouldBeTrue)
		c.So(isBluezAudio("BLUEZ.abcd.1234"), ShouldBeTrue)
		c.So(isBluezAudio("bluz.abcd.1234"), ShouldBeFalse)
		c.So(isBluezAudio("hbc.bluz.1234"), ShouldBeFalse)
		c.So(isBluezAudio("BLUZ.abcd.1234"), ShouldBeFalse)
	})
}

func Test_createBluezVirtualCardPorts(t *testing.T) {
	Convey("createBluezVirtualCardPorts", t, func(c C) {
		name := "bluez.abcd.1234"
		desc := "headset"
		var cards pulse.CardPortInfos
		cards = append(cards, pulse.CardPortInfo{
			PortInfo: pulse.PortInfo{
				Name:        name,
				Description: desc,
				Priority:    0,
				Available:   pulse.AvailableTypeUnknow,
			},
			Direction: pulse.DirectionSink,
			Profiles:  pulse.ProfileInfos2{},
		})
		portinfos := createBluezVirtualCardPorts(cards)
		c.So(len(portinfos), ShouldEqual, 0)

		cards[0].Profiles = append(cards[0].Profiles, pulse.ProfileInfo2{
			Name:        "a2dp_sink",
			Description: "A2DP",
			Priority:    0,
			NSinks:      1,
			NSources:    0,
			Available:   pulse.AvailableTypeUnknow,
		})
		portinfos = createBluezVirtualCardPorts(cards)
		c.So(len(portinfos), ShouldEqual, 1)

		cards[0].Profiles = append(cards[0].Profiles, pulse.ProfileInfo2{
			Name:        "headset_head_unit",
			Description: "Headset",
			Priority:    0,
			NSinks:      1,
			NSources:    0,
			Available:   pulse.AvailableTypeUnknow,
		})
		portinfos = createBluezVirtualCardPorts(cards)
		c.So(len(portinfos), ShouldEqual, 2)
		c.So(portinfos[0].Name, ShouldEqual, name+"(headset_head_unit)")
		c.So(portinfos[1].Name, ShouldEqual, name+"(a2dp_sink)")
	})
}

func Test_createBluezVirtualSinkPorts(t *testing.T) {
	Convey("createBluezVirtualSinkPorts", t, func(c C) {
		name := "bluez.abcd.1234"
		var ports []Port
		ports = append(ports, Port{
			Name:        name,
			Description: "headsert",
			Available:   byte(pulse.AvailableTypeUnknow),
		})

		ret := createBluezVirtualSinkPorts(ports)
		c.So(len(ret), ShouldEqual, 2)
		c.So(ret[0].Name, ShouldEqual, name+"(headset_head_unit)")
		c.So(ret[1].Name, ShouldEqual, name+"(a2dp_sink)")
	})
}

func Test_createBluezVirtualSourcePorts(t *testing.T) {
	Convey("createBluezVirtualSourcePorts", t, func(c C) {
		name := "bluez.abcd.1234"
		var ports []Port
		ports = append(ports, Port{
			Name:        name,
			Description: "headset",
			Available:   byte(pulse.AvailableTypeUnknow),
		})

		ret := createBluezVirtualSourcePorts(ports)
		c.So(len(ret), ShouldEqual, 1)
		c.So(ret[0].Name, ShouldEqual, name+"(headset_head_unit)")
	})
}
