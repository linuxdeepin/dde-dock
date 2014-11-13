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

package accounts

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

func (*testWrap) TestNewUserShell(c *C.C) {
	shell, err := getNewUserDefaultShell("testdata/adduser.conf")
	c.Check(err, C.Not(C.NotNil))
	c.Check(shell, C.Equals, "/bin/zsh")

	shell, err = getNewUserDefaultShell("testdata/adduser1.conf")
	c.Check(err, C.Not(C.NotNil))
	c.Check(shell, C.Equals, "")

	shell, err = getNewUserDefaultShell("testdata/xxxxxxx")
	c.Check(err, C.NotNil)
	c.Check(shell, C.Equals, "")
}
