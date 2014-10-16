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

package i18n_dependency

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

func (*testWrapper) TestInstallPkg(c *C.C) {
	err := InstallDependentPackages("en_US.UTF-8",
		"testdata/pkg_depends.json",
		"testdata/support_languages.json")
	//c.Check(err, C.Not(C.NotNil))
	if err != nil {
		c.Skip(err.Error())
	}
}

func (*testWrapper) TestGetPkgList(c *C.C) {
	list, err := getPkgDependList("testdata/pkg_depends.json")
	c.Check(err, C.Not(C.NotNil))
	c.Check(list, C.NotNil)

	trList := getDependentPkgListByKey("tr", "zh_CN.UTF-8",
		"testdata/pkg_depends.json", list)
	c.Check(len(trList), C.Not(C.Equals), 0)
	wwList := getDependentPkgListByKey("ww", "zh_CN.UTF-8",
		"testdata/pkg_depends.json", list)
	c.Check(len(wwList), C.Equals, 0)

	list, err = getPkgDependList("testdata/xxxxx.json")
	c.Check(err, C.NotNil)
	c.Check(list, C.Not(C.NotNil))
}

type pkgDependInfoTest struct {
	info   *dependentPkgInfo
	locale string
	list   []string
}

func (*testWrapper) TestParsePkgInfo(c *C.C) {
	//zh_CN.UTF-8: zh-hans, CN, ""
	infos := []pkgDependInfoTest{
		pkgDependInfoTest{
			&dependentPkgInfo{"", 1,
				"", "language-pack-"},
			"zh_CN.UTF-8",
			[]string{"language-pack-zh-hans"},
		},
		pkgDependInfoTest{
			&dependentPkgInfo{"", 2,
				"firefox", "firefox-locale-"},
			"zh_CN.UTF-8",
			[]string{"firefox-locale-zh-hans",
				"firefox-locale-zh-hans-cn"},
		},
		pkgDependInfoTest{
			&dependentPkgInfo{"", 3,
				"calligra-libs", "calligra-l10n-"},
			"zh_CN.UTF-8",
			[]string{"calligra-l10n-zh-hans",
				"calligra-l10n-zh-hanscn"},
		},
		pkgDependInfoTest{
			&dependentPkgInfo{"fi", 0,
				"epiphany", "xul-ext-mozvoikko"},
			"zh_CN.UTF-8",
			[]string{"xul-ext-mozvoikko"},
		},
		pkgDependInfoTest{
			&dependentPkgInfo{"", 5,
				"xxxx", "xxxxxxxxxxxxx"},
			"zh_CN.UTF-8",
			nil,
		},
		pkgDependInfoTest{
			nil,
			"zh_CN.UTF-8",
			nil,
		},
	}

	for _, info := range infos {
		tmp := parseFormatType(info.info, info.locale, "testdata/support_languages.json")
		c.Check(isStrListEqual(tmp, info.list),
			C.Equals, true)
	}
}

func isStrListEqual(l1, l2 []string) bool {
	len1 := len(l1)
	len2 := len(l2)

	if len1 != len2 {
		return false
	}

	for i := 0; i < len1; i++ {
		if l1[i] != l2[i] {
			return false
		}
	}

	return true
}

type langCodeTest struct {
	locale  string
	lcode   string
	ccode   string
	variant string
}

func (*testWrapper) TestGetLangCodeInfo(c *C.C) {
	infos := []langCodeTest{
		langCodeTest{"aa_DJ.UTF-8", "aa", "", ""},
		langCodeTest{"be_BY.UTF-8", "be", "BY", ""},
		langCodeTest{"be_BY@latin", "be", "", "latin"},
		langCodeTest{"xxx_XX.UTF-8", "", "", ""},
	}

	for _, info := range infos {
		lcode, ccode, variant := getLangCodeByLocale(
			info.locale,
			"testdata/support_languages.json")
		c.Check(lcode, C.Equals, info.lcode)
		c.Check(ccode, C.Equals, info.ccode)
		c.Check(variant, C.Equals, info.variant)
	}
}
