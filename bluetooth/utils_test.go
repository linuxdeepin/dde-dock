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

package bluetooth

import (
	"testing"

	C "gopkg.in/check.v1"
	"pkg.deepin.io/lib/dbus1"
)

func TestT(t *testing.T) { C.TestingT(t) }

type testWrapper struct{}

var _ = C.Suite(&testWrapper{})

func (*testWrapper) TestIsDBusObjectKeyExists(c *C.C) {
	data := make(dbusObjectData)
	c.Check(isDBusObjectKeyExists(data, "key"), C.Equals, false)
	data["key"] = dbus.MakeVariant(int32(8))
	c.Check(isDBusObjectKeyExists(data, "key"), C.Equals, true)
}

func (*testWrapper) TestGetDBusObjectValueString(c *C.C) {
	data := make(dbusObjectData)
	c.Check(getDBusObjectValueString(data, "key"), C.Equals, "")
	data["key"] = dbus.MakeVariant("value")
	c.Check(getDBusObjectValueString(data, "key"), C.Equals, "value")
}

func (*testWrapper) TestGetDBusObjectValueInt16(c *C.C) {
	data := make(dbusObjectData)
	c.Check(getDBusObjectValueInt16(data, "key"), C.Equals, int16(0))
	data["key"] = dbus.MakeVariant(int16(8))
	c.Check(getDBusObjectValueInt16(data, "key"), C.Equals, int16(8))
}
