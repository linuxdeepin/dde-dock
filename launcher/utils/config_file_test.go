package utils

import (
	C "launchpad.net/gocheck"
	"os"
)

type ConfigFileTestSuite struct {
}

var _ = C.Suite(&ConfigFileTestSuite{})

func (self *ConfigFileTestSuite) TestConfigFilePath(c *C.C) {
	old := os.Getenv("HOME")
	os.Setenv("HOME", "../testdata/")
	c.Assert(ConfigFilePath("launcher/test.ini"), C.Equals, "../testdata/.config/launcher/test.ini")
	os.Setenv("HOME", old)
}
