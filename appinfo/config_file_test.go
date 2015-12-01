package appinfo

import (
	"os"

	C "launchpad.net/gocheck"
)

type ConfigFileTestSuite struct {
}

var _ = C.Suite(&ConfigFileTestSuite{})

func (*ConfigFileTestSuite) TestConfigFilePath(c *C.C) {
	os.Setenv("XDG_CONFIG_HOME", "../testdata/.config")
	c.Assert(ConfigFilePath("launcher/test.ini"), C.Equals, "../testdata/.config/launcher/test.ini")
}
