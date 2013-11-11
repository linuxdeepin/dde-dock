package main

import "dlib/dbus"

type UserInfo struct {
	AccountType string
	AutoLogin   bool
	FaceRecog   bool
}

func (info *UserInfo) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.UserInfo",
		"/com/deepin/daemon/UserInfo",
		"com.deepin.daemon.UserInfo",
	}
}

func (info *UserInfo) reset (propName string) {
}

func (info *UserInfo) SetUserPasswd(user, passwd string) bool {
	return true
}

func main() {
	info := UserInfo{}
	dbus.InstallOnSession(&info)
	select {}
}
