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

package fonts

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

var infos = []StyleInfo{
	{
		Id:          "WenQuanYi Micro Hei",
		Families:    []string{"WenQuanYi Micro Hei", "文泉译微米黑"},
		FamilyLangs: []string{"en", "zh"},
		StyleList:   []string{"Regular", "Bold"},
	},
	{
		Id:          "Source Code Pro",
		Families:    []string{"Source Code Pro"},
		FamilyLangs: []string{"en"},
		StyleList:   []string{"Regular", "Bold"},
	},
	{
		Id:          "Source San Hans",
		Families:    []string{"Source San Hans", "思源黑体"},
		FamilyLangs: []string{"en", "zh"},
		StyleList:   []string{"Regular", "Bold"},
	},
	{
		Id:          "WenQuanYi Micro Mono Hei",
		Families:    []string{"WenQuanYi Micro Mono Hei", "文泉译等宽微米黑"},
		FamilyLangs: []string{"en", "zh"},
		StyleList:   []string{"Regular", "Bold"},
	},
}

func (*testWrapper) TestInfoList(c *C.C) {
	standList, monoList := getStyleInfoList()
	c.Check(len(standList), C.Not(C.Equals), 0)
	c.Check(len(monoList), C.Not(C.Equals), 0)
}

func (*testWrapper) TestGetNameList(c *C.C) {
	list := getNameStrList(infos)
	tmp := []string{
		"WenQuanYi Micro Hei",
		"Source Code Pro",
		"Source San Hans",
		"WenQuanYi Micro Mono Hei",
	}

	c.Check(isStrListEqual(list, tmp), C.Equals, true)
}

func (*testWrapper) TestGetStyleList(c *C.C) {
	list := getStyleList("Source Code Pro", infos)
	tmp := []string{"Regular", "Bold"}
	c.Check(isStrListEqual(list, tmp), C.Equals, true)
}

func isStrListEqual(list1, list2 []string) bool {
	l1 := len(list1)
	l2 := len(list2)

	if l1 != l2 {
		return false
	}

	for i := 0; i < l1; i++ {
		if list1[i] != list2[i] {
			return false
		}
	}

	return true
}
