/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
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

package zoneinfo

import (
	C "launchpad.net/gocheck"
	"testing"
)

type testWrapper struct{}

func init() {
	C.Suite(&testWrapper{})
}

func Test(t *testing.T) {
	C.TestingT(t)
}

func (*testWrapper) TestGetZoneList(c *C.C) {
	var ret = []string{
		"Europe/Andorra",
		"Asia/Dubai",
		"Asia/Kabul",
		"Europe/Tirane",
		"Asia/Yerevan",
	}

	list, err := getZoneListFromFile("testdata/zone1970.tab")
	c.Check(err, C.Equals, nil)
	for i, _ := range list {
		c.Check(list[i], C.Equals, ret[i])
	}
}

func (*testWrapper) TestZoneValid(c *C.C) {
	var infos = []struct {
		zone  string
		valid bool
	}{
		{
			zone:  "Asia/Shanghai",
			valid: true,
		},
		//{
		//zone:  "Asia/Beijing",
		//valid: true,
		//},
		{
			zone:  "Asia/xxxx",
			valid: false,
		},
	}

	for _, info := range infos {
		c.Check(IsZoneValid(info.zone), C.Equals, info.valid)
	}
}
