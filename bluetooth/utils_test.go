/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package bluetooth

import (
	C "gopkg.in/check.v1"
	"pkg.deepin.io/lib/dbus"
	"testing"
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
