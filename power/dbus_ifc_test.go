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

package power

import (
	libupower "dbus/org/freedesktop/upower"
	"encoding/xml"
	C "launchpad.net/gocheck"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/dbus/introspect"
	"pkg.linuxdeepin.com/lib/log"
	"testing"
)

type DBusInterfaceTest struct{}

func init() {
	if logger == nil {
		logger = log.NewLogger("dde-daemon/power")
	}

	C.Suite(&DBusInterfaceTest{})
}

func (dbusIfc *DBusInterfaceTest) SetUpSuite(c *C.C) {
	_, err := dbus.SystemBus()
	if err != nil {
		c.Skip(err.Error())
	}
}

func Test(t *testing.T) {
	C.TestingT(t)
}

type interfaceInfo struct {
	name    string
	ifcType string
	exist   bool
}

func (dbusIfc *DBusInterfaceTest) TestUPowerInterfaceExist(c *C.C) {
	root, err := getNodeInfo("org.freedesktop.UPower",
		"/org/freedesktop/UPower")
	if err != nil {
		c.Error(err)
		return
	}

	IfcInfos := []interfaceInfo{
		interfaceInfo{"LidIsPresent", "property", true},
		interfaceInfo{"LidIsClosed", "property", true},
		interfaceInfo{"OnBattery", "property", true},
		interfaceInfo{"EnumerateDevices", "method", true},
		interfaceInfo{"DeviceAdded", "signal", true},
		interfaceInfo{"DeviceRemoved", "signal", true},
		interfaceInfo{"LidIsPresent11", "property", false},
		interfaceInfo{"Changed", "signal", false},
	}

	for _, info := range IfcInfos {
		c.Check(isInterfaceNameFound(info.name, info.ifcType, root),
			C.Equals, info.exist)
	}
}

func (dbusIfc *DBusInterfaceTest) TestDeviceInterfaceExist(c *C.C) {
	upowerObj, err := libupower.NewUpower("org.freedesktop.UPower",
		"/org/freedesktop/UPower")
	if err != nil {
		c.Error(err)
		return
	}

	devs, err := upowerObj.EnumerateDevices()
	if err != nil {
		c.Error(err)
		return
	}

	if len(devs) == 0 {
		return
	}

	root, err := getNodeInfo("org.freedesktop.UPower", devs[0])
	if err != nil {
		c.Error(err)
		return
	}

	IfcInfos := []interfaceInfo{
		interfaceInfo{"Percentage", "property", true},
		interfaceInfo{"State", "property", true},
		interfaceInfo{"IsPresent", "property", true},
		interfaceInfo{"Type", "property", true},
		interfaceInfo{"Changed", "signal", false},
	}

	for _, info := range IfcInfos {
		c.Check(isInterfaceNameFound(info.name, info.ifcType, root),
			C.Equals, info.exist)
	}
}

func getNodeInfo(dest string, path dbus.ObjectPath) (*introspect.NodeInfo, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	var xmlString string
	dbusObj := conn.Object(dest, path)
	dbusObj.Call("org.freedesktop.DBus.Introspectable.Introspect",
		dbus.FlagNoAutoStart).Store(&xmlString)

	var node introspect.NodeInfo
	err = xml.Unmarshal([]byte(xmlString), &node)
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func isInterfaceNameFound(name, t string, root *introspect.NodeInfo) bool {
	if isInterfaceNameFoundNoChild(name, t, root) {
		return true
	}

	for _, node := range root.Children {
		if isInterfaceNameFoundNoChild(name, t, &node) {
			return true
		}
	}

	return false
}

func isInterfaceNameFoundNoChild(name, t string, node *introspect.NodeInfo) bool {
	for _, ifcInfo := range node.Interfaces {
		switch t {
		case "property":
			for _, info := range ifcInfo.Properties {
				if name == info.Name {
					return true
				}
			}
		case "method":
			for _, info := range ifcInfo.Methods {
				if name == info.Name {
					return true
				}
			}
		case "signal":
			for _, info := range ifcInfo.Signals {
				if name == info.Name {
					return true
				}
			}
		}
	}

	return false
}
