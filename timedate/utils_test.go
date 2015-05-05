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

package timedate

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

var list = []string{
	"home",
	"hello",
	"world",
	"goodbye",
}

func (*testWrapper) TestItemInList(c *C.C) {
	var infos = []struct {
		name  string
		exist bool
	}{
		{
			name:  "hello",
			exist: true,
		},
		{
			name:  "goodbye",
			exist: true,
		},
		{
			name:  "helloxxx",
			exist: false,
		},
	}

	for _, info := range infos {
		c.Check(isItemInList(info.name, list), C.Equals, info.exist)
	}
}

func (*testWrapper) TestAddItem(c *C.C) {
	var infos = []struct {
		name   string
		length int
		added  bool
	}{
		{
			name:   "hello",
			length: 4,
			added:  false,
		},
		{
			name:   "helloxxx",
			length: 5,
			added:  true,
		},
	}

	for _, info := range infos {
		tList, added := addItemToList(info.name, list)
		c.Check(len(tList), C.Equals, info.length)
		c.Check(added, C.Equals, info.added)
		c.Check(isItemInList(info.name, tList), C.Equals, true)
	}
}

func (*testWrapper) TestDeleteItem(c *C.C) {
	var infos = []struct {
		name    string
		length  int
		deleted bool
	}{
		{
			name:    "hello",
			length:  3,
			deleted: true,
		},
		{
			name:    "helloxxx",
			length:  4,
			deleted: false,
		},
	}

	for _, info := range infos {
		tList, deleted := deleteItemFromList(info.name, list)
		c.Check(len(tList), C.Equals, info.length)
		c.Check(deleted, C.Equals, info.deleted)
		c.Check(isItemInList(info.name, tList), C.Equals, false)
	}
}
