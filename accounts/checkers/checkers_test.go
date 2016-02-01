/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package checkers

import (
	C "launchpad.net/gocheck"
	"testing"
)

type testWrapper struct{}

func init() {
	C.Suite(&testWrapper{})
}

func Test(t *testing.T) {
	C.TestingT(t)
}

func (*testWrapper) TestCheckUsername(c *C.C) {
	type checkRet struct {
		name string
		code ErrorCode
	}

	var infos = []checkRet{
		{"", ErrCodeEmpty},
		{"a1111111111111111111111111111111", 0},
		{"a11111111111111111111111111111111", ErrCodeLenMoreThen},
		{"root", ErrCodeSystemUsed},
		{"123", ErrCodeFirstNotLower},
		{"a123*&", ErrCodeInvalidChar},
	}

	for _, v := range infos {
		tmp := CheckUsernameValid(v.name)
		if v.code == 0 {
			c.Check(tmp, C.Equals, (*ErrorInfo)(nil))
		} else {
			c.Check(tmp.Code, C.Equals, v.code)
		}
	}
}

func (*testWrapper) TestGetUsernames(c *C.C) {
	var datas = []struct {
		name string
		ret  bool
	}{
		{
			name: "test1",
			ret:  true,
		},
		{
			name: "test2",
			ret:  true,
		},
		{
			name: "test3",
			ret:  false,
		},
	}

	names, err := getAllUsername("testdata/passwd")
	c.Check(err, C.Equals, nil)
	c.Check(len(names), C.Equals, 2)

	for _, data := range datas {
		c.Check(isStrInArray(data.name, names), C.Equals, data.ret)
		c.Check(isStrInArray(data.name, names), C.Equals, data.ret)
		c.Check(isStrInArray(data.name, names), C.Equals, data.ret)
	}
}
