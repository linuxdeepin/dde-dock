package bluetooth

import (
	C "launchpad.net/gocheck"
	"pkg.linuxdeepin.com/lib/dbus"
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
