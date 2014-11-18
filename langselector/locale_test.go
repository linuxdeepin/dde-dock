/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package langselector

import (
	C "launchpad.net/gocheck"
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

func (t *TestWrapper) TestConstructPamFile(c *C.C) {
	example := `LANG=en_US.UTF-8
LANGUAGE=en_US.UTF-8
LC_TIME="zh_CN.UTF-8"
`

	c.Check(constructPamFile("en_US.UTF-8",
		"testdata/pam_environment"), C.Equals, example)
	c.Check(constructPamFile("en_US.UTF-8", "xxxxxx"),
		C.Equals, generatePamContents("en_US.UTF-8"))
}

func (t *TestWrapper) TestGetLocale(c *C.C) {
	l, err := getLocaleFromFile("testdata/pam_environment")
	c.Check(err, C.Not(C.NotNil))
	c.Check(l, C.Equals, "zh_CN.UTF-8")

	l = getLocale()
	c.Check(len(l), C.Not(C.Equals), 0)
}

func (t *TestWrapper) TestWriteUserLocale(c *C.C) {
	c.Check(writeUserLocalePam("zh_CN.UTF-8", "testdata/pam"),
		C.Not(C.NotNil))
	os.RemoveAll("testdata/pam")
	c.Check(writeUserLocalePam("zh_CN.UTF-8", "/xxxxxxxxx"),
		C.NotNil)
}

func (t *TestWrapper) TestLocaleInfoList(c *C.C) {
	list, err := getLocaleInfoList("testdata/support_languages.json")
	c.Check(len(list), C.Not(C.Equals), 0)
	c.Check(err, C.Not(C.NotNil))

	list, err = getLocaleInfoList("testdata/zzxxxxxxx")
	c.Check(len(list), C.Equals, 0)
	c.Check(err, C.NotNil)
}

// TODO: panic in jenkins for the dbus interface
// func (t *TestWrapper) TestNetwork(c *C.C) {
// 	_, err := isNetworkEnable()
// 	c.Check(err, C.Not(C.NotNil))
// }

// func (t *TestWrapper) TestNotify(c *C.C) {
// 	err := sendNotify("", "", "Test")
// 	//c.Check(err, C.Not(C.NotNil))
// 	if err != nil {
// 		c.Skip(err.Error())
// 	}
// }
