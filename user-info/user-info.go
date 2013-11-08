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
		"com.deepin.dss.userinfo",
		"/com/deepin/dss/userinfo",
		"com.deepin.dss.userinfo",
	}
}

func (info *UserInfo) SetAccountType (user, accountType string) bool {
	return true
}

func (info *UserInfo) SetAutoLogin (user string, autoLogin bool) bool {
	return true
}

func (info *UserInfo) SetFaceRecog (user string, aceRecog bool) bool {
	return true
}

func (info *UserInfo) SetUserPasswd (user, passwd string) bool {
	return true
}

func (info *UserInfo) GetUserInfo (user string) map[string]string {
	return nil
}

func main () {
	info := UserInfo {}
	dbus.InstallOnSession (&info)
	select {}
}
