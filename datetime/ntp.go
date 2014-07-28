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

package datetime

import (
	"net"
	"strconv"
	"time"
)

const (
	_NTP_HOST = "0.pool.ntp.org"
)

func (obj *Manager) syncNtpTime() bool {
	for i := 0; i < 10; i++ {
		t, err := getNtpTime(obj.CurrentTimezone)
		if err == nil && t != nil {
			dStr, tStr := getDateTimeAny(t)
			logger.Infof("Date: %s, Time: %s", dStr, tStr)
			setDate.SetCurrentDate(dStr)
			setDate.SetCurrentTime(tStr)
			return true
		}
	}

	return false
}

func (obj *Manager) syncNtpThread() {
	for {
		obj.syncNtpTime()
		timer := time.NewTimer(time.Minute * 10)
		select {
		case <-timer.C:
		case <-obj.quitChan:
			obj.ntpRunning = false
			return
		}
	}
}

func (obj *Manager) enableNtp(enable bool) bool {
	if enable {
		if obj.ntpRunning {
			go obj.syncNtpTime()
			logger.Debug("Ntp is running")
			return true
		}

		obj.ntpRunning = true
		go obj.syncNtpThread()
	} else {
		if obj.ntpRunning {
			logger.Debug("Ntp will quit....")
			obj.quitChan <- true
		}

		obj.ntpRunning = false
	}

	return true
}

func getDateTimeAny(t *time.Time) (dStr, tStr string) {
	dStr += strconv.FormatInt(int64(t.Year()), 10) + "-" + strconv.FormatInt(int64(t.Month()), 10) + "-" + strconv.FormatInt(int64(t.Day()), 10)
	tStr += strconv.FormatInt(int64(t.Hour()), 10) + ":" + strconv.FormatInt(int64(t.Minute()), 10) + ":" + strconv.FormatInt(int64(t.Second()), 10)

	return dStr, tStr
}

func getNtpTime(locale string) (*time.Time, error) {
	if len(locale) < 1 {
		locale = "UTC"
	}

	if !timezoneIsValid(locale) {
		logger.Warningf("'%s': invalid locale", locale)
		locale = "UTC"
	}

	logger.Info("Locale:", locale)
	raddr, err := net.ResolveUDPAddr("udp", _NTP_HOST+":123")
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

	l := time.FixedZone(locale, 0)
	if l == nil {
		return nil, err
	}

	t := time.Date(1900, 1, 1, 0, 0, 0, 0, l).Add(time.Duration(nsec)).Local()

	return &t, nil
}
