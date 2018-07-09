/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package x_event_monitor

import (
	"testing"

	C "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	C.TestingT(t)
}

type mySuite struct{}

var _ = C.Suite(&mySuite{})

func (s *mySuite) TestUtils(c *C.C) {
	c.Check(hasButtonFlag(2), C.Equals, true)
	c.Check(hasButtonFlag(0), C.Equals, false)

	c.Check(hasKeyFlag(4), C.Equals, true)
	c.Check(hasKeyFlag(0), C.Equals, false)
}
