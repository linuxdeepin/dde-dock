package network

import "strings"
import "fmt"

// TODO implement dde-api/proxy

// example of /etc/environment
// http_proxy="http://127.0.0.1:0/"
// https_proxy="https://127.0.0.1:0/"
// ftp_proxy="ftp://127.0.0.1:0/"
// socks_proxy="socks://127.0.0.1:0/"

func (m *Manager) GetProxy(proxyType string) (addr, port string, err error) {
	err = checkProxyType(proxyType)
	if err != nil {
		return
	}
	proxy, err := ubuntuGetProxy(proxyType)
	if err != nil {
		return
	}
	proxy = strings.TrimPrefix(proxy, proxyType+"://")
	proxy = strings.TrimSuffix(proxy, "/")
	a := strings.Split(proxy, ":")
	if len(a) == 2 {
		addr = a[0]
		port = a[1]
	}
	return
}

// if address is empty, means to remove proxy setting
func (m *Manager) SetProxy(proxyType, addr, port string) (err error) {
	err = checkProxyType(proxyType)
	if err != nil {
		return
	}
	var proxy string
	if len(addr) > 0 {
		proxy = proxyType + "://" + addr + port + "/"
	}
	err = ubuntuSetProxy(proxyType, proxy)
	return
}

func checkProxyType(proxyType string) (err error) {
	switch proxyType {
	case "http", "https", "ftp", "socks":
	default:
		err = fmt.Errorf("not a valid proxy type: %s", proxyType)
		logger.Error(err)
	}
	return
}

func newUbuntuSystemService() (s *SystemService, err error) {
	// s, err = NewSystemService("com.ubuntu.SystemService", "/com/ubuntu/SystemService")
	s, err = NewSystemService("com.ubuntu.SystemService", "/")
	if err != nil {
		logger.Error(err)
	}
	return
}

func ubuntuGetProxy(proxyType string) (proxy string, err error) {
	var s *SystemService
	s, err = newUbuntuSystemService()
	if err != nil {
		logger.Error(err)
		return
	}
	proxy, err = s.GetProxy(proxyType)
	if err != nil {
		// do not print error here, for if a proxy is not defined, it
		// will get error too
		return
	}
	return
}

func ubuntuSetProxy(proxyType, proxy string) (err error) {
	var s *SystemService
	s, err = newUbuntuSystemService()
	if err != nil {
		logger.Error(err)
		return
	}
	ok, err := s.SetProxy(proxyType, proxy)
	if !ok {
		err = fmt.Errorf("set proxy failed: %s, %s", proxyType, proxy)
		logger.Error(err)
		return
	}
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func ubuntuSetNoProxy(proxyType, proxy string) (err error) {
	var s *SystemService
	s, err = newUbuntuSystemService()
	if err != nil {
		logger.Error(err)
		return
	}
	ok, err := s.SetProxy(proxyType, proxy)
	if !ok {
		err = fmt.Errorf("set proxy failed: %s, %s", proxyType, proxy)
		logger.Error(err)
		return
	}
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
