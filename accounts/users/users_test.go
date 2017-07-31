/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package users

import (
	C "gopkg.in/check.v1"
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

	list, err := getAdminUserList("testdata/group", "testdata/sudoers_deepin")
	c.Check(err, C.Equals, nil)

	for _, data := range datas {
		c.Check(isStrInArray(data.name, list), C.Equals, data.admin)
	}
}

func (*testWrapper) TestGetAutoLoginUser(c *C.C) {
	// lightdm
	name, err := getIniKeys("testdata/autologin/lightdm_autologin.conf",
		kfGroupLightdmSeat,
		[]string{kfKeyLightdmAutoLoginUser}, []string{""})
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "wen")
	name, err = getIniKeys("testdata/autologin/lightdm.conf",
		kfGroupLightdmSeat,
		[]string{kfKeyLightdmAutoLoginUser}, []string{""})
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "")
	_, err = getIniKeys("testdata/autologin/xxxxx.conf", "", nil, nil)
	c.Check(err, C.Not(C.Equals), nil)

	// gdm
	name, err = getIniKeys("testdata/autologin/custom_autologin.conf",
		kfGroupGDM3Daemon, []string{kfKeyGDM3AutomaticEnable,
			kfKeyGDM3AutomaticLogin}, []string{"True", ""})
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "wen")
	name, err = getIniKeys("testdata/autologin/custom.conf",
		kfGroupGDM3Daemon, []string{kfKeyGDM3AutomaticEnable,
			kfKeyGDM3AutomaticLogin}, []string{"True", ""})
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "")

	// kdm
	name, err = getIniKeys("testdata/autologin/kdmrc_autologin",
		kfGroupKDMXCore, []string{kfKeyKDMAutoLoginEnable,
			kfKeyKDMAutoLoginUser}, []string{"true", ""})
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "wen")
	name, err = getIniKeys("testdata/autologin/kdmrc",
		kfGroupKDMXCore, []string{kfKeyKDMAutoLoginEnable,
			kfKeyKDMAutoLoginUser}, []string{"true", ""})
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "")

	// sddm
	name, err = getIniKeys("testdata/autologin/sddm_autologin.conf",
		kfGroupSDDMAutologin,
		[]string{kfKeySDDMUser}, []string{""})
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "wen")
	name, err = getIniKeys("testdata/autologin/sddm.conf",
		kfGroupSDDMAutologin,
		[]string{kfKeySDDMUser}, []string{""})
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "")

	// lxdm
	name, err = getIniKeys("testdata/autologin/lxdm_autologin.conf",
		kfGroupLXDMBase,
		[]string{kfKeyLXDMAutologin}, []string{""})
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "wen")
	name, err = getIniKeys("testdata/autologin/lxdm.conf",
		kfGroupLXDMBase,
		[]string{kfKeyLXDMAutologin}, []string{""})
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "")

	// slim
	name, err = parseSlimConfig("testdata/autologin/slim_autologin.conf",
		"", false)
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "wen")
	name, err = parseSlimConfig("testdata/autologin/slim.conf", "", false)
	c.Check(err, C.Equals, nil)
	c.Check(name, C.Equals, "")
	// cp 'testdata/autologin/slim.conf' to '/tmp/slim_tmp.conf'
	// _, err = parseSlimConfig("/tmp/slim_tmp.conf", "wen", true)
	// c.Check(err, C.Equals, nil)
	// name, err = parseSlimConfig("/tmp/slim_tmp.conf", "", false)
	// c.Check(err, C.Equals, nil)
	// c.Check(name, C.Equals, "wen")

	m, err := getDefaultDM("testdata/autologin/default-display-manager")
	c.Check(err, C.Equals, nil)
	c.Check(m, C.Equals, "lightdm")
	_, err = getDefaultDM("testdata/autologin/xxxxx")
	c.Check(err, C.Not(C.Equals), nil)
}

func (*testWrapper) TestXSession(c *C.C) {
	session, _ := getIniKeys("testdata/autologin/lightdm.conf", kfGroupLightdmSeat,
		[]string{"user-session"}, []string{""})
	c.Check(session, C.Equals, "deepin")
	session, _ = getIniKeys("testdata/autologin/sddm.conf", kfGroupSDDMAutologin,
		[]string{kfKeySDDMSession}, []string{""})
	c.Check(session, C.Equals, "kde-plasma.desktop")
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

func (*testWrapper) TestGetAdmGroup(c *C.C) {
	groups, users, err := getAdmGroupAndUser("testdata/sudoers_deepin")
	c.Check(err, C.Equals, nil)
	c.Check(isStrInArray("sudo", groups), C.Equals, true)
	c.Check(isStrInArray("root", users), C.Equals, true)

	groups, users, err = getAdmGroupAndUser("testdata/sudoers_arch")
	c.Check(err, C.Equals, nil)
	c.Check(isStrInArray("sudo", groups), C.Equals, true)
	c.Check(isStrInArray("wheel", groups), C.Equals, true)
	c.Check(isStrInArray("root", users), C.Equals, true)
}

func (*testWrapper) TestDMFromService(c *C.C) {
	dm, _ := getDMFromSystemService("testdata/autologin/display-manager.service")
	c.Check(dm, C.Equals, "lightdm")
}
