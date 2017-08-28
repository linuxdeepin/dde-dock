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
