package network

import "strings"
import "fmt"
import "dlib/gio-2.0"

// example of /etc/environment
// http_proxy="http://127.0.0.1:0/"
// https_proxy="https://127.0.0.1:0/"
// ftp_proxy="ftp://127.0.0.1:0/"
// socks_proxy="socks://127.0.0.1:0/"

const (
	proxyHttp  = "http"
	proxyHttps = "https"
	proxyFtp   = "ftp"
	proxySocks = "socks"

	gsettingsIdProxy = "com.deepin.dde.proxy"
	gkeyHttpProxy    = "http-proxy"
	gkeyHttpsProxy   = "https-proxy"
	gkeyFtpProxy     = "ftp-proxy"
	gkeySocksProxy   = "socks-proxy"
)

var (
	proxySettings = gio.NewSettings(gsettingsIdProxy)
)

func (m *Manager) GetProxy(proxyType string) (addr, port string, err error) {
	proxy, err := doGetProxy(proxyType)
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
	var proxy string
	if len(addr) > 0 {
		if len(port) == 0 {
			err = fmt.Errorf("proxy port is empty")
			return
		}
		if !strings.HasPrefix(addr, proxyType+"://") {
			addr = proxyType + "://" + addr
		}
		proxy = addr + ":" + port + "/"
	}
	err = doSetProxy(proxyType, proxy)
	return
}

func checkProxyType(proxyType string) (err error) {
	switch proxyType {
	case proxyHttp, proxyHttps, proxyFtp, proxySocks:
	default:
		err = fmt.Errorf("not a valid proxy type: %s", proxyType)
		logger.Error(err)
	}
	return
}

func doGetProxy(proxyType string) (proxy string, err error) {
	switch proxyType {
	case proxyHttp:
		proxy = proxySettings.GetString(gkeyHttpProxy)
	case proxyHttps:
		proxy = proxySettings.GetString(gkeyHttpsProxy)
	case proxyFtp:
		proxy = proxySettings.GetString(gkeyFtpProxy)
	case proxySocks:
		proxy = proxySettings.GetString(gkeySocksProxy)
	default:
		err = fmt.Errorf("not a valid proxy type: %s", proxyType)
		logger.Error(err)
	}
	return
}

func doSetProxy(proxyType, proxy string) (err error) {
	var ok bool
	switch proxyType {
	case proxyHttp:
		ok = proxySettings.SetString(gkeyHttpProxy, proxy)
	case proxyHttps:
		ok = proxySettings.SetString(gkeyHttpsProxy, proxy)
	case proxyFtp:
		ok = proxySettings.SetString(gkeyFtpProxy, proxy)
	case proxySocks:
		ok = proxySettings.SetString(gkeySocksProxy, proxy)
	default:
		err = fmt.Errorf("not a valid proxy type: %s", proxyType)
		logger.Error(err)
	}
	if !ok {
		err = fmt.Errorf("set proxy value to gsettings failed: %s, %s", proxyType, proxy)
		logger.Error(err)
	}
	return
}
