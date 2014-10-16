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

package language_info

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

type localeTest struct {
	locale string
	ret    bool
}

func (*testWrapper) TestGetLangList(c *C.C) {
	langList, err := GetLanguageInfoList("testdata/xxx.json")
	c.Check(err, C.NotNil)
	c.Check(len(langList), C.Equals, 0)

	langList, err = GetLanguageInfoList("testdata/support_languages.json")
	c.Check(err, C.Not(C.NotNil))
	c.Check(len(langList), C.Not(C.Equals), 0)
}

func (*testWrapper) TestLocaleValid(c *C.C) {
	var infos = []localeTest{
		localeTest{
			"ar_EG.UTF-8",
			true,
		},
		localeTest{
			"be_BY.UTF-8",
			true,
		},
		localeTest{
			"xxx_XX.UTF-8",
			false,
		},
	}

	for _, info := range infos {
		c.Check(IsLocaleValid(info.locale, "testdata/support_languages.json"),
			C.Equals, info.ret)
	}
}

type langCodeTest struct {
	locale  string
	lcode   string
	ccode   string
	variant string
	ok      bool
}

func (*testWrapper) TestGetLangCodeInfo(c *C.C) {
	infos := []langCodeTest{
		langCodeTest{"aa_DJ.UTF-8", "aa", "", "", false},
		langCodeTest{"be_BY.UTF-8", "be", "BY", "", false},
		langCodeTest{"be_BY@latin", "be", "", "latin", false},
		langCodeTest{"xxx_XX.UTF-8", "", "", "", true},
	}

	for _, info := range infos {
		tmp, err := GetCodeInfoByLocale(info.locale, "testdata/support_languages.json")
		if info.ok {
			c.Check(err, C.NotNil)
		} else {
			c.Check(err, C.Not(C.NotNil))
		}
		c.Check(tmp.LangCode, C.Equals, info.lcode)
		c.Check(tmp.CountryCode, C.Equals, info.ccode)
		c.Check(tmp.Variant, C.Equals, info.variant)
	}
}
