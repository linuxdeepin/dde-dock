// +build integration

package systeminfo

import (
	C "launchpad.net/gocheck"
	"testing"
)

type testWrapper struct{}

func Test(t *testing.T) {
	C.TestingT(t)
}

func init() {
	C.Suite(&testWrapper{})
}

func (*testWrapper) TestSystemVersion(c *C.C) {
	version, err := getVersionFromDeepin("/etc/deepin-version")
	c.Check(err, C.Not(C.NotNil))
	c.Check(len(version), C.Not(C.Equals), 0)
	version, err = getVersionFromLsb("/etc/lsb-release")
	c.Check(err, C.Not(C.NotNil))
	c.Check(len(version), C.Not(C.Equals), 0)
	version, err = getVersionFromDeepin("xxxxxxxxxx")
	c.Check(err, C.NotNil)
	c.Check(len(version), C.Equals, 0)
}

func (*testWrapper) TestSystemCPU(c *C.C) {
	cpu, err := getCPUInfoFromFile("/proc/cpuinfo")
	c.Check(err, C.Not(C.NotNil))
	c.Check(len(cpu), C.Not(C.Equals), 0)
	cpu, err = getCPUInfoFromFile("xxxxxxxxxxx")
	c.Check(err, C.NotNil)
	c.Check(cpu, C.Equals, "")
}

func (*testWrapper) TestSystemMemory(c *C.C) {
	mem, err := getMemoryCapFromFile("/proc/meminfo")
	c.Check(err, C.Not(C.NotNil))
	c.Check(mem, C.Not(C.Equals), uint64(0))
	mem, err = getMemoryCapFromFile("xxxxxxxxx")
	c.Check(err, C.NotNil)
	c.Check(mem, C.Equals, uint64(0))
}

func (*testWrapper) TestSystemType(c *C.C) {
	t, err := getSystemType()
	c.Check(err, C.Not(C.NotNil))
	c.Check(t, C.Not(C.Equals), 0)
}

func (*testWrapper) TestSystemDisk(c *C.C) {
	caps, err := getDiskCap()
	c.Check(err, C.Not(C.NotNil))
	c.Check(caps, C.Not(C.Equals), 0)
}
