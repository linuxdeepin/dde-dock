package main

import (
	. "launchpad.net/gocheck"
)

func init() {
	manager := &Manager{}
	Suite(manager)
}

func (m *Manager) TestRemoveDevice(c *C) {
	devs := make([]*device, 0)
	devs = append(devs, &device{"path1", 0})
	devs = append(devs, &device{"path2", 0})
	devs, _ = m.removeDevice(devs, "path1")
	c.Check(len(devs), Equals, 1)
	logger.Info(devs)
}
