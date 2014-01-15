/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

package main

import (
	"launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

var op *Manager

func init() {
	objectMap = make(map[int32]*ObjectInfo)

	op = &Manager{}
	op.setPropName("DiskList")
	op.listenSignalChanged()
	gocheck.Suite(op)
}

func (op *Manager) TestMount(c *gocheck.C) {
	for _, info := range op.DiskList {
		if !info.CanUnmount {
			op.DeviceMount(info.Id)
		}
	}
	if c.Failed() {
                c.Error("Test Mount Failed")
        }
}

func (op *Manager) TestUnmount(c *gocheck.C) {
	for _, info := range op.DiskList {
		if info.CanUnmount {
			op.DeviceUnmount(info.Id)
		}
	}
	if c.Failed() {
                c.Error("Test Mount Failed")
        }
}

func (op *Manager) TestEject(c *gocheck.C) {
	for _, info := range op.DiskList {
		if info.CanEject {
			op.DeviceEject(info.Id)
		}
	}
	if c.Failed() {
                c.Error("Test Mount Failed")
        }
}
