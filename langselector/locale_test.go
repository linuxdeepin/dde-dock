/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package langselector

import (
	C "gopkg.in/check.v1"
	"os"
	"testing"
)

func Test(t *testing.T) {
	C.TestingT(t)
}

type TestWrapper struct{}

func init() {
	C.Suite(&TestWrapper{})
}

type localeDescTest struct {
	locale string
	ret    bool
}

func (t *TestWrapper) TestGenerateLocaleEnvFile(c *C.C) {
	example := `LANG=en_US.UTF-8
LANGUAGE=en_US
LC_TIME="zh_CN.UTF-8"`

	c.Check(generateLocaleEnvFile("en_US.UTF-8",
		"testdata/pam_environment"), C.Equals, example)
}

func (t *TestWrapper) TestGetLocale(c *C.C) {
	l, err := getLocaleFromFile("testdata/pam_environment")
	c.Check(err, C.Not(C.NotNil))
	c.Check(l, C.Equals, "zh_CN.UTF-8")

	l = getCurrentUserLocale()
	c.Check(len(l), C.Not(C.Equals), 0)
}

func (t *TestWrapper) TestWriteUserLocale(c *C.C) {
	c.Check(writeLocaleEnvFile("zh_CN.UTF-8", "testdata/pam"),
		C.Not(C.NotNil))
	os.RemoveAll("testdata/pam")
	c.Check(writeLocaleEnvFile("zh_CN.UTF-8", "/xxxxxxxxx"),
		C.NotNil)
}
