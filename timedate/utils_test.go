/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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

package timedate

import (
	C "gopkg.in/check.v1"
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

func (*testWrapper) TestFilterNilStr(c *C.C) {
	var infos = []struct {
		list   []string
		hasNil bool
		ret    []string
	}{
		{
			list:   []string{"abs", "apt", "", "pacman"},
			hasNil: true,
			ret:    []string{"abs", "apt", "pacman"},
		},
		{
			list:   []string{"c", "go", "python"},
			hasNil: false,
			ret:    []string{"c", "go", "python"},
		},
	}

	for _, info := range infos {
		list, hasNil := filterNilString(info.list)
		c.Check(hasNil, C.Equals, info.hasNil)
		c.Check(len(list), C.Equals, len(info.ret))
	}
}
