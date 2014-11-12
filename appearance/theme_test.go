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

package appearance

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

func (*testWrapper) TestReadThemeConfig(c *C.C) {
	var tmp = Theme{
		Name:        "Custom",
		DisplayName: "自定义",
		GtkTheme:    "Deepin-Gray",
		IconTheme:   "Faenza",
		SoundTheme:  "LinuxDeepin",
		CursorTheme: "DMZ-White",
		FontName:    "Source Han Sans",
		FontMono:    "Source Code Pro",
		Background:  "file:///usr/share/personalization/themes/Deepin/wallpapers/time%201.jpg",
		FontSize:    10,
	}

	lang := os.Getenv("LANG")
	os.Setenv("LANG", "zh_CN.UTF-8")

	info, err := getThemeInfoFromFile("testdata/theme.ini")
	c.Check(err, C.Not(C.NotNil))
	c.Check(isThemeInfoSame(&info, &tmp), C.Equals, true)

	os.Setenv("LANG", lang)
}
