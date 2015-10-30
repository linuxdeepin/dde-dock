package utils

import (
	C "launchpad.net/gocheck"
	"os"
)

type ConfigFileTestSuite struct {
}

var _ = C.Suite(&ConfigFileTestSuite{})

func (*ConfigFileTestSuite) TestConfigFilePath(c *C.C) {
	os.Setenv("XDG_CONFIG_HOME", "../testdata/.config")
	c.Assert(ConfigFilePath("launcher/test.ini"), C.Equals, "../testdata/.config/launcher/test.ini")
}
