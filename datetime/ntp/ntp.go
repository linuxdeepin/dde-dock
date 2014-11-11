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

package ntp

import (
	"fmt"
	"net"
	. "pkg.linuxdeepin.com/dde-daemon/datetime/utils"
	"sync"
	"time"
)

const (
	ntpHost = "0.pool.ntp.org"
)

const (
	NTPStateDisabled int32 = 0
	NTPStateEnabled  int32 = 1
)

// update time when timezone changed
var Timezone string

var (
	stateLock sync.Mutex
	ntpState  int32
)

func InitNtpModule() error {
	ntpState = NTPStateDisabled
	return InitSetDateTime()
}

func FiniNtpModule() {
	DestroySetDateTime()
}

func Enabled(enable bool, zone string) {
	Timezone = zone
	if enable {
		if ntpState == NTPStateEnabled {
			go SyncNetworkTime()
			return
		}

		stateLock.Lock()
		ntpState = NTPStateEnabled
		stateLock.Unlock()
		go syncThread()
	} else {
		stateLock.Lock()
		ntpState = NTPStateDisabled
		stateLock.Unlock()
	}

	return
}

func SyncNetworkTime() bool {
	for i := 0; i < 10; i++ {
		dStr, tStr, err := getDateTime()
		if err == nil {
			SetDate(dStr)
			SetTime(tStr)
			return true
		}
	}

	return false
}

func syncThread() {
	for {
		SyncNetworkTime()
		timer := time.NewTimer(time.Minute * 10)
		select {
		case <-timer.C:
			if ntpState == NTPStateDisabled {
				return
			}
		}
	}
}

func getDateTime() (string, string, error) {
	t, err := getNetworkTime()
	if err != nil {
		return "", "", err
	}

	dStr := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
	tStr := fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second())

	return dStr, tStr, nil
}

func getNetworkTime() (*time.Time, error) {
	loc, err := time.LoadLocation(Timezone)
	if err != nil {
		return nil, err
	}
	time.Local = loc

	raddr, err := net.ResolveUDPAddr("udp", ntpHost+":123")
	if err != nil {
		return nil, err
	}

	data := make([]byte, 48)
	data[0] = 3<<3 | 3

	con, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}

	defer con.Close()

	_, err = con.Write(data)
	if err != nil {
		return nil, err
	}

	con.SetDeadline(time.Now().Add(5 * time.Second))

	_, err = con.Read(data)
	if err != nil {
		return nil, err
	}

	var sec, frac uint64
	sec = uint64(data[43]) | uint64(data[42])<<8 | uint64(data[41])<<16 |
		uint64(data[40])<<24
	frac = uint64(data[47]) | uint64(data[46])<<8 | uint64(data[45])<<16 |
		uint64(data[44])<<24

	nsec := sec * 1e9
	nsec += (frac * 1e9) >> 32

	t := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(nsec)).Local()

	return &t, nil
}
