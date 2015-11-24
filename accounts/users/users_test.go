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

package users

import (
	C "launchpad.net/gocheck"
	dutils "pkg.deepin.io/lib/utils"
	"testing"
)

type testWrapper struct{}

func init() {
	C.Suite(&testWrapper{})
}

func Test(t *testing.T) {
	C.TestingT(t)
}

func (*testWrapper) TestGetUserInfos(c *C.C) {
	var names = []string{"test1", "test2", "vbox"}
	infos, err := getUserInfosFromFile("testdata/passwd")
	c.Check(err, C.Equals, nil)
	c.Check(len(infos), C.Equals, 3)
	for i, info := range infos {
		c.Check(info.Name, C.Equals, names[i])
	}
}

func (*testWrapper) TestUserInfoValid(c *C.C) {
	var infos = []struct {
		name  UserInfo
		valid bool
	}{
		{
			UserInfo{Name: "root", Uid: "0", Gid: "0"},
			false,
		},
		{
			UserInfo{Name: "test1", Shell: "/bin/bash", Uid: "1000", Gid: "1000"},
			true,
		},
		{
			UserInfo{Name: "test1", Shell: "/bin/false", Uid: "1000", Gid: "1000"},
			false,
		},
		{
			UserInfo{Name: "test1", Shell: "/bin/bash", Uid: "60000", Gid: "60000"},
			true,
		},
		{
			UserInfo{Name: "test1", Shell: "/bin/bash", Uid: "999", Gid: "999"},
			false,
		},
		{
			UserInfo{Name: "test1", Shell: "/bin/bash", Uid: "60001", Gid: "60001"},
			false,
		},
		{
			UserInfo{Name: "test1", Shell: "/bin/nologin", Uid: "1000", Gid: "1000"},
			false,
		},
		{
			UserInfo{Name: "test3", Shell: "/bin/bash", Uid: "1000", Gid: "1000"},
			false,
		},
		{
			UserInfo{Name: "test4", Shell: "/bin/bash", Uid: "1000", Gid: "1000"},
			false,
		},
	}

	for _, v := range infos {
		c.Check(v.name.isHumanUser("testdata/shadow", "testdata/login.defs"), C.Equals, v.valid)
	}
}

func (*testWrapper) TestFoundUserInfo(c *C.C) {
	info, err := getUserInfo(UserInfo{Name: "test1"}, "testdata/passwd")
	c.Check(err, C.Equals, nil)
	c.Check(info.Name, C.Equals, "test1")

	info, err = getUserInfo(UserInfo{Uid: "1001"}, "testdata/passwd")
	c.Check(err, C.Equals, nil)
	c.Check(info.Name, C.Equals, "test1")

	info, err = getUserInfo(UserInfo{Name: "1006"}, "testdata/passwd")
	c.Check(err, C.NotNil)

	info, err = getUserInfo(UserInfo{Uid: "1006"}, "testdata/passwd")
	c.Check(err, C.NotNil)

	info, err = getUserInfo(UserInfo{Uid: "1006"}, "testdata/xxxxx")
	c.Check(err, C.NotNil)
}

func (*testWrapper) TestAdminUser(c *C.C) {
	var datas = []struct {
		name  string
		admin bool
	}{
		{
			name:  "wen",
			admin: true,
		},
		{
			name:  "test1",
			admin: true,
		},
		{
			name:  "test2",
			admin: false,
		},
	}

	list, err := getAdminUserList("testdata/group")
	c.Check(err, C.Equals, nil)

	for _, data := range datas {
		c.Check(isStrInArray(data.name, list), C.Equals, data.admin)
	}
}

func (*testWrapper) TestGetAutoLoginUser(c *C.C) {
	name, err := getLightdmAutoLoginUser("testdata/autologin/lightdm_autologin.conf")
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "wen")
	name, err = getLightdmAutoLoginUser("testdata/autologin/lightdm.conf")
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "")
	_, err = getLightdmAutoLoginUser("testdata/autologin/xxxxx.conf")
	c.Check(err, C.Not(C.Equals), nil)

	name, err = getGDMAutoLoginUser("testdata/autologin/custom_autologin.conf")
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "wen")
	name, err = getGDMAutoLoginUser("testdata/autologin/custom.conf")
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "")
	_, err = getGDMAutoLoginUser("testdata/autologin/xxxx.conf")
	c.Check(err, C.Not(C.Equals), nil)

	name, err = getKDMAutoLoginUser("testdata/autologin/kdmrc_autologin")
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "wen")
	name, err = getKDMAutoLoginUser("testdata/autologin/kdmrc")
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "")
	_, err = getKDMAutoLoginUser("testdata/autologin/xxxxx")
	c.Check(err, C.Not(C.Equals), nil)

	m, err := getDefaultDM("testdata/autologin/default-display-manager")
	c.Check(err, C.Equals, nil)
	c.Check(m, C.Equals, "lightdm")
	_, err = getDefaultDM("testdata/autologin/xxxxx")
	c.Check(err, C.Not(C.Equals), nil)
}

func (*testWrapper) TestWriteStrvData(c *C.C) {
	var (
		datas = []string{"123", "abc", "xyz"}
		file  = "/tmp/write_strv"
	)
	err := writeStrvToFile(datas, file, 0644)
	c.Check(err, C.Equals, nil)

	md5, _ := dutils.SumFileMd5(file)
	c.Check(md5, C.Equals, "0b188e42e5f8d5bc5a6560ce68d5fbc6")
}

func (*testWrapper) TestGetDefaultShell(c *C.C) {
	shell, err := getDefaultShell("testdata/adduser.conf")
	c.Check(err, C.Equals, nil)
	c.Check(shell, C.Equals, "/bin/zsh")

	shell, err = getDefaultShell("testdata/adduser1.conf")
	c.Check(err, C.Equals, nil)
	c.Check(shell, C.Equals, "")

	_, err = getDefaultShell("testdata/xxxxx.conf")
	c.Check(err, C.Not(C.Equals), nil)
}

func (*testWrapper) TestStrInArray(c *C.C) {
	var array = []string{"abc", "123", "xyz"}

	var datas = []struct {
		value string
		ret   bool
	}{
		{
			value: "abc",
			ret:   true,
		},
		{
			value: "xyz",
			ret:   true,
		},
		{
			value: "abcd",
			ret:   false,
		},
	}

	for _, data := range datas {
		c.Check(isStrInArray(data.value, array), C.Equals, data.ret)
	}
}
