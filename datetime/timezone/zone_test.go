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

package timezone

import (
	C "launchpad.net/gocheck"
	"os"
	"testing"
)

type testWrapper struct{}

func Test(t *testing.T) {
	C.TestingT(t)
}

func init() {
	C.Suite(&testWrapper{})
}

func (*testWrapper) TestGetZoneList(c *C.C) {
	zoneList := GetZoneInfoList()
	c.Check(len(zoneList), C.Equals, len(zoneWhiteList))
}

type testZoneSummary struct {
	zone      string
	dstOffset string
	ret       bool
}

var infos = []testZoneSummary{
	{
		zone:      "America/Atka",
		dstOffset: "−09:00",
		ret:       true,
	},
	{
		zone:      "US/Aleutian",
		dstOffset: "−09:00",
		ret:       true,
	},
	{
		zone:      "Pacific/Niue",
		dstOffset: "",
		ret:       true,
	},
	{
		zone:      "Pacific/Niueiii",
		dstOffset: "",
		ret:       false,
	},
}

func (*testWrapper) TestZoneValid(c *C.C) {
	for _, info := range infos {
		c.Check(IsZoneValid(info.zone), C.Equals, info.ret)
	}
}

func (*testWrapper) TestNewDSTInfo(c *C.C) {
	type testDST struct {
		zone string
		dst  DSTInfo
	}

	var infos = []testDST{
		{
			zone: "Asia/Shanghai",
			dst: DSTInfo{
				Enter:     0,
				Leave:     0,
				DSTOffset: 0,
			},
		},
		{
			zone: "America/New_York",
			dst: DSTInfo{
				Enter:     1394348400,
				Leave:     1414907999,
				DSTOffset: -14400,
			},
		},
	}

	for _, info := range infos {
		e, l, ok := getDSTTime(info.zone, 2014)
		c.Check(e, C.Equals, info.dst.Enter)
		c.Check(l, C.Equals, info.dst.Leave)

		if !ok {
			continue
		}
		off := getOffsetByTimestamp(info.zone, e)
		c.Check(off, C.Equals, info.dst.DSTOffset)
	}
}

func (*testWrapper) TestZoneDesc(c *C.C) {
	var infos = []zoneDesc{
		{
			zone: "Asia/Beijing",
			desc: "Beijing",
		},
		{
			zone: "Pacific/Johnston",
			desc: "Pacific/Johnston",
		},
	}

	lang := os.Getenv("LANG")
	os.Setenv("LANG", "en_US.UTF-8")
	for _, info := range infos {
		c.Check(getZoneDesc(info.zone), C.Equals, info.desc)
	}
	os.Setenv("LANG", lang)
}

func (*testWrapper) TestFindDSTInfo(c *C.C) {
	type testDSTData struct {
		data dstData
		ret  bool
	}

	// dst info for 2014
	var infos = []testDSTData{
		{
			data: dstData{
				zone: "US/Alaska",
				dst: DSTInfo{
					Enter:     1394362800,
					Leave:     1414922399,
					DSTOffset: -28800,
				},
			},
			ret: true,
		},
		{
			data: dstData{
				zone: "Atlantic/Azores",
				dst: DSTInfo{
					Enter:     1396141200,
					Leave:     1414285199,
					DSTOffset: 0,
				},
			},
			ret: true,
		},
		{
			data: dstData{
				zone: "Pacific/Apia",
				dst: DSTInfo{
					Enter:     1396706399,
					Leave:     1411826400,
					DSTOffset: 50400,
				},
			},
			ret: true,
		},
		{
			data: dstData{
				zone: "Asia/Shanghai",
				dst:  DSTInfo{},
			},
			ret: false,
		},
	}

	for _, info := range infos {
		dst, err := findDSTInfo(info.data.zone, "testdata/dst_data")
		if info.ret {
			c.Check(err, C.IsNil)
			c.Check(info.data.dst.Enter, C.Equals, dst.Enter)
			c.Check(info.data.dst.Leave, C.Equals, dst.Leave)
			c.Check(info.data.dst.DSTOffset, C.Equals, dst.DSTOffset)
		} else {
			c.Check(dst, C.IsNil)
			c.Check(err, C.NotNil)
		}
	}
}
