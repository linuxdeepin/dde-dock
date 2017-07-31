/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package zoneinfo

import (
	C "gopkg.in/check.v1"
	"testing"
)

type testWrapper struct{}

func init() {
	C.Suite(&testWrapper{})
}

func Test(t *testing.T) {
	C.TestingT(t)
}

func (*testWrapper) TestGetZoneList(c *C.C) {
	var ret = []string{
		"Europe/Andorra",
		"Asia/Dubai",
		"Asia/Kabul",
		"Europe/Tirane",
		"Asia/Yerevan",
	}

	list, err := getZoneListFromFile("testdata/zone1970.tab")
	c.Check(err, C.Equals, nil)
	for i, _ := range list {
		c.Check(list[i], C.Equals, ret[i])
	}
}

func (*testWrapper) TestZoneValid(c *C.C) {
	var infos = []struct {
		zone  string
		valid bool
	}{
		{
			zone:  "Asia/Shanghai",
			valid: true,
		},
		//{
		//zone:  "Asia/Beijing",
		//valid: true,
		//},
		{
			zone:  "Asia/xxxx",
			valid: false,
		},
	}

	for _, info := range infos {
		c.Check(IsZoneValid(info.zone), C.Equals, info.valid)
	}
}
