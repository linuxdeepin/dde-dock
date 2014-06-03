package network

import (
	// liblogger "dlib/logger"
	. "launchpad.net/gocheck"
	"testing"
	"time"
)

func init() {
	manager := &Manager{}
	Suite(manager)
}

// func (m *Manager) TestRemoveDevice(c *C) {
// 	devs := make([]*device, 0)
// 	devs = append(devs, &device{Path: "path1", State: 0})
// 	devs = append(devs, &device{Path: "path2", State: 0})
// 	devs = m.doRemoveDevice(devs, "path1")
// 	c.Check(len(devs), Equals, 1)
// }

func TestModuleStartStop(t *testing.T) {
	// logger.SetLogLevel(liblogger.LEVEL_DEBUG)
	Stop()
	Stop()
	go func() {
		time.Sleep(3 * time.Second)
		Start()
		Stop()
		Stop()
	}()
	Start()
	Stop()
	Stop()
	time.Sleep(30 * time.Second)
}
