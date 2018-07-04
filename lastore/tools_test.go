/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

package lastore

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (*MySuite) TestStrSliceSetEqual(c *C) {
	c.Assert(strSliceSetEqual([]string{}, []string{}), Equals, true)
	c.Assert(strSliceSetEqual([]string{"a"}, []string{"a"}), Equals, true)
	c.Assert(strSliceSetEqual([]string{"a", "b"}, []string{"a"}), Equals, false)
	c.Assert(strSliceSetEqual([]string{"a", "b", "d"}, []string{"b", "d", "a"}), Equals, true)
}
