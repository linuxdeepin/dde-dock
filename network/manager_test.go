package network

import (
	. "launchpad.net/gocheck"
)

func init() {
	manager := &Manager{}
	Suite(manager)
}

func (m *Manager) TestRemoveDevice(c *C) {
	devs := make([]*device, 0)
	devs = append(devs, &device{Path: "path1", State: 0})
	devs = append(devs, &device{Path: "path2", State: 0})
	devs = m.doRemoveDevice(devs, "path1")
	c.Check(len(devs), Equals, 1)
	logger.Info(devs)
}
