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

package libkeyboard

import (
	C "launchpad.net/gocheck"
	"testing"
)

type testWrapper struct{}

func Test(t *testing.T) {
	C.TestingT(t)
}

func init() {
	C.Suite(&testWrapper{})
}

func (*testWrapper) TestGetLayout(c *C.C) {
	layout := getLayoutFromFile("testdata/keyboard")
	c.Check(layout, C.Equals, "us;")

	layout = getLayoutFromFile("testdata/xxxxxx")
	c.Check(layout, C.Equals, "us;")
}

func (*testWrapper) TestQTCursorBlink(c *C.C) {
	c.Check(setQtCursorBlink(1200, "testdata/Trolltech.conf"),
		C.Not(C.NotNil))
}

func (*testWrapper) TestLayoutList(c *C.C) {
	v, err := getLayoutListByFile("testdata/base.xml")
	c.Check(err, C.Not(C.NotNil))
	c.Check(v, C.NotNil)

	v, err = getLayoutListByFile("testdata/xxxxxx.xml")
	c.Check(err, C.NotNil)
	c.Check(v, C.Not(C.NotNil))
}
