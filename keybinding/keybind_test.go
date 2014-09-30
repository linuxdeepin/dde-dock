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

package keybinding

import (
	C "launchpad.net/gocheck"
	"testing"
)

type testWrap struct{}

func Test(t *testing.T) {
	C.TestingT(t)
}

func init() {
	C.Suite(&testWrap{})
}

type weirdTestData struct {
	shortcut string
	weird    string
}

func (t *testWrap) TestKeysym2Weird(c *C.C) {
	datas := []weirdTestData{
		weirdTestData{"control-'", "control-apostrophe"},
		weirdTestData{"control-\"", "control-quotedbl"},
		weirdTestData{"control-\\", "control-backslash"},
		weirdTestData{"control-_", "control-underscore"},
		weirdTestData{"control-|", "control-bar"},
		weirdTestData{"control--", "control-minus"},
		weirdTestData{"control-s", "control-s"},
		weirdTestData{"control-mod4-s", "control-mod4-s"},
		weirdTestData{"control-mod2-mod4-s", "control-mod2-mod4-s"},
	}

	for _, info := range datas {
		c.Check(convertKeysym2Weird(info.shortcut), C.Equals, info.weird)
	}
}
