/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package appinfo

import (
	"os"

	C "gopkg.in/check.v1"
)

type ConfigFileTestSuite struct {
}

var _ = C.Suite(&ConfigFileTestSuite{})

func (*ConfigFileTestSuite) TestConfigFilePath(c *C.C) {
	os.Setenv("XDG_CONFIG_HOME", "../testdata/.config")
	c.Assert(ConfigFilePath("launcher/test.ini"), C.Equals, "../testdata/.config/launcher/test.ini")
}
