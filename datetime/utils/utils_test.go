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

package utils

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

func (*testWrapper) TestStrInList(c *C.C) {
	list := []string{
		"1111",
		"2222",
		"aaaa",
		"bbbb",
	}

	c.Check(IsStrInList("1111", list), C.Equals, true)
	c.Check(IsStrInList("aaaa", list), C.Equals, true)
	c.Check(IsStrInList("xxxxxx", list), C.Equals, false)
}

type dateInfo struct {
	d   int32
	ret bool
}

func (*testWrapper) TestYearValid(c *C.C) {
	var infos = []dateInfo{
		{
			d:   2014,
			ret: true,
		},
		{
			d:   1970,
			ret: true,
		},
		{
			d:   1,
			ret: false,
		},
		{
			d:   -2014,
			ret: false,
		},
	}

	for _, info := range infos {
		c.Check(IsYearValid(info.d), C.Equals, info.ret)
	}
}

func (*testWrapper) TestMonthValid(c *C.C) {
	var infos = []dateInfo{
		{
			d:   1,
			ret: true,
		},
		{
			d:   12,
			ret: true,
		},
		{
			d:   13,
			ret: false,
		},
		{
			d:   0,
			ret: false,
		},
	}

	for _, info := range infos {
		c.Check(IsMonthValid(info.d), C.Equals, info.ret)
	}
}

func (*testWrapper) TestDayValid(c *C.C) {
	type mdayInfo struct {
		year  int32
		month int32
		mday  int32
		ret   bool
	}

	var infos = []mdayInfo{
		{
			year:  2014,
			month: 3,
			mday:  1,
			ret:   true,
		},
		{
			year:  2014,
			month: 3,
			mday:  31,
			ret:   true,
		},
		{
			year:  2014,
			month: 3,
			mday:  32,
			ret:   false,
		},
		{
			year:  2014,
			month: 3,
			mday:  0,
			ret:   false,
		},
		{
			year:  2014,
			month: 2,
			mday:  1,
			ret:   true,
		},
		{
			year:  2014,
			month: 2,
			mday:  28,
			ret:   true,
		},
		{
			year:  2014,
			month: 2,
			mday:  29,
			ret:   false,
		},
		{
			year:  2014,
			month: 2,
			mday:  0,
			ret:   false,
		},
		{
			year:  2014,
			month: 4,
			mday:  1,
			ret:   true,
		},
		{
			year:  2014,
			month: 4,
			mday:  30,
			ret:   true,
		},
		{
			year:  2014,
			month: 4,
			mday:  0,
			ret:   false,
		},
		{
			year:  2014,
			month: 4,
			mday:  31,
			ret:   false,
		},
		{
			year:  2000,
			month: 2,
			mday:  29,
			ret:   true,
		},
	}

	for _, info := range infos {
		c.Check(IsDayValid(info.year, info.month, info.mday),
			C.Equals, info.ret)
	}
}

func (*testWrapper) TestHourValid(c *C.C) {
	var infos = []dateInfo{
		{
			d:   0,
			ret: true,
		},
		{
			d:   23,
			ret: true,
		},
		{
			d:   24,
			ret: false,
		},
		{
			d:   -1,
			ret: false,
		},
	}

	for _, info := range infos {
		c.Check(IsHourValid(info.d), C.Equals, info.ret)
	}
}

func (*testWrapper) TestMinuteValid(c *C.C) {
	var infos = []dateInfo{
		{
			d:   0,
			ret: true,
		},
		{
			d:   59,
			ret: true,
		},
		{
			d:   60,
			ret: false,
		},
		{
			d:   -1,
			ret: false,
		},
	}

	for _, info := range infos {
		c.Check(IsMinuteValid(info.d), C.Equals, info.ret)
	}
}

func (*testWrapper) TestSecondValid(c *C.C) {
	var infos = []dateInfo{
		{
			d:   0,
			ret: true,
		},
		{
			d:   60,
			ret: true,
		},
		{
			d:   61,
			ret: false,
		},
		{
			d:   -1,
			ret: false,
		},
	}

	for _, info := range infos {
		c.Check(IsSecondValid(info.d), C.Equals, info.ret)
	}
}

func (*testWrapper) TestLeapYear(c *C.C) {
	var infos = []dateInfo{
		{
			d:   2000,
			ret: true,
		},
		{
			d:   2004,
			ret: true,
		},
		{
			d:   2100,
			ret: false,
		},
		{
			d:   2014,
			ret: false,
		},
	}

	for _, info := range infos {
		c.Check(IsLeapYear(info.d), C.Equals, info.ret)
	}
}
