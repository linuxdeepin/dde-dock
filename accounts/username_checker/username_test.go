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

package username_checker

import (
	C "launchpad.net/gocheck"
	"testing"
)

type testWrapper struct{}

func init() {
	initErrorInfo()
	C.Suite(&testWrapper{})
}

func Test(t *testing.T) {
	C.TestingT(t)
}

type errorInfo struct {
	name string
	code int32
	ret  bool
}

func (*testWrapper) TestUsernameValid(c *C.C) {
	var infos = []errorInfo{
		{
			name: "",
			code: ErrCodeEmpty,
			ret:  true,
		},
		{
			name: "Asdf",
			code: ErrCodeFirstNotLower,
			ret:  true,
		},
		{
			name: "root",
			code: ErrCodeSystemUsed,
			ret:  true,
		},
		{
			name: "xx12$%",
			code: ErrCodeInvalidChar,
			ret:  true,
		},
		{
			name: "0",
			code: ErrCodeFirstNotLower,
			ret:  true,
		},
		{
			name: "A",
			code: ErrCodeFirstNotLower,
			ret:  true,
		},
		{
			name: "-",
			code: ErrCodeFirstNotLower,
			ret:  true,
		},
		{
			name: "_",
			code: ErrCodeFirstNotLower,
			ret:  true,
		},
		{
			name: "i",
			code: -1,
			ret:  false,
		},
		{
			name: "iixxx",
			code: -1,
			ret:  false,
		},
	}

	for _, info := range infos {
		err := CheckUsernameValid(info.name)
		if info.ret {
			c.Check(err.Code, C.Equals, info.code)
		} else {
			c.Check(err, C.Not(C.NotNil))
		}
	}
}

func (*testWrapper) TestUsernameExist(c *C.C) {
	var infos = []errorInfo{
		{
			name: "",
			code: -1,
			ret:  false,
		},
		{
			name: "root",
			code: ErrCodeSystemUsed,
			ret:  true,
		},
		{
			name: "wen",
			code: ErrCodeExist,
			ret:  true,
		},
		{
			name: "iixxx",
			code: -1,
			ret:  false,
		},
	}

	for _, info := range infos {
		_, err := checkNameExist(info.name, "testdata/passwd")
		if info.ret {
			c.Check(err.Code, C.Equals, info.code)
		} else {
			c.Check(err, C.Not(C.NotNil))
		}
	}
}
