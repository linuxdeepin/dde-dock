package main

import "dlib/dbus"

type UserInfo struct {
	AccountType	string
	AutoLogin	bool
	FaceRecog	bool

	AccoutTypeChanged	func (user, accountType string)
	AutoLoginChanged	func (user string, autoLogin bool)
	FaceRecogChanged	func (user string, aceRecog bool)
}

func (info *UserInfo) GetDBusInfo () dbus.DBusInfo {
	return dbus.DBusInfo {
		"com.deepin.daemon.UserInfo",
		"/com/deepin/daemon/UserInfo",
		"com.deepin.daemon.UserInfo",
	}
}

func (info *UserInfo) SetUserPasswd (user, passwd string) bool {
	return true
}

func main () {
	info := UserInfo {}
	dbus.InstallOnSession (&info)
	select {}
}
