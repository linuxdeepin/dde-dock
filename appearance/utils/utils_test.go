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

package utils

import (
	C "launchpad.net/gocheck"
	"os"
	"testing"
)

type testWrapper struct{}

func init() {
	C.Suite(&testWrapper{})
}

func Test(t *testing.T) {
	C.TestingT(t)
}

func (*testWrapper) TestGetInfo(c *C.C) {
	dir := PathInfo{
		BaseName: "",
		FilePath: "testdata",
		FileFlag: FileFlagUserOwned,
	}

	themeConditions := []string{"theme.ini"}
	infoList := GetInfoListFromDirs([]PathInfo{dir}, themeConditions)
	c.Check(len(infoList), C.Equals, 2)
}

var infos = []PathInfo{
	{
		BaseName: "Deepin",
		FilePath: "testdata/Deepin",
		FileFlag: FileFlagUserOwned,
	},
	{
		BaseName: "Custom",
		FilePath: "testdata/Custom",
		FileFlag: FileFlagUserOwned,
	},
}

func (*testWrapper) TestNameAndList(c *C.C) {
	c.Check(IsNameInInfoList("Deepin", infos), C.Equals, true)
	c.Check(IsNameInInfoList("xxx", infos), C.Equals, false)

	_, err := GetInfoByName("Deepin", infos)
	c.Check(err, C.Not(C.NotNil))
	_, err = GetInfoByName("xxx", infos)
	c.Check(err, C.NotNil)

	list := GetBaseNameList(infos)
	c.Check(len(list), C.Equals, 2)

	c.Check(GetFileFlagByName("Deepin", infos),
		C.Equals, int32(1))
}

func (*testWrapper) TestSeedValid(c *C.C) {
	c.Check(isThumbSeedValid("--gtk"), C.Equals, true)
	c.Check(isThumbSeedValid("--icon"), C.Equals, true)
	c.Check(isThumbSeedValid("--cursor"), C.Equals, true)
	c.Check(isThumbSeedValid("--font"), C.Equals, false)
}

func (*testWrapper) TestGetConfig(c *C.C) {
	homeDir := os.Getenv("HOME")
	os.Setenv("HOME", "testdata")
	c.Check(GetUserGtk2Config(), C.Equals, "testdata/.gtkrc-2.0")
	os.Setenv("HOME", homeDir)
}

func (*testWrapper) TestWriteConfig(c *C.C) {
	c.Check(
		WriteUserGtk3Config("testdata/settings.ini",
			"gtk-theme-name", "Deepin"),
		C.Not(C.NotNil))
	c.Check(
		WriteUserGtk3Config("testdata/settings.ini",
			"", ""),
		C.NotNil)

	c.Check(
		WriteUserGtk2Config("testdata/gtkrc-2.0",
			"gtk-theme-name", "Deepin"),
		C.Not(C.NotNil))
	c.Check(
		WriteUserGtk3Config("testdata/gtkrc-2.0",
			"", ""),
		C.NotNil)
}
