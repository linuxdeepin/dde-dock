package network

import "dlib/dbus"

var __conn *dbus.Conn = nil

func getBus() *dbus.Conn {
	if __conn == nil {
		var err error
		__conn, err = dbus.SystemBus()
		if err != nil {
			panic(err)
		}
	}
	return __conn
}
