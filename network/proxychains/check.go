/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package proxychains

import (
	"net"
	"regexp"
	"strings"
)

var ipReg = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)

func checkType(type0 string) bool {
	switch type0 {
	case "http", "socks4", "socks5":
		return true
	default:
		return false
	}
}

func checkIP(ipstr string) bool {
	if !ipReg.MatchString(ipstr) {
		return false
	}

	ip := net.ParseIP(ipstr)
	if ip == nil {
		return false
	}
	return true
}

func checkUser(user string) bool {
	if strings.ContainsAny(user, "\t ") {
		return false
	}

	return true
}

func checkPassword(password string) bool {
	return checkUser(password)
}
