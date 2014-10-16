// +build integration

package systeminfo

import (
	C "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) {
	C.TestingT(t)
}

func init() {
	C.Suite(NewSystemInfo())
}

func (si *SystemInfo) TestSystem(c *C.C) {
	if ret := GetVersion(); ret == "" {
		c.Error("GetVersion failed")
		return
	}

	if ret := GetCpuInfo(); len(ret) < 1 {
		c.Error("GetCpuInfo failed")
		return
	}

	if ret := GetMemoryCap(); ret == 0 {
		c.Error("GetMemoryCap failed")
		return
	}

	if ret := GetSystemType(); ret == 0 {
		c.Error("GetSystemType failed")
		return
	}

	if ret := GetDiskCap(); ret == 0 {
		c.Error("GetDiskCap failed")
		return
	}
}
