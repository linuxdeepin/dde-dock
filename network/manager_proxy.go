/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
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

package network

import "strings"
import "fmt"
import "dlib/gio-2.0"
import "regexp"

// example of /etc/environment
// http_proxy="http://127.0.0.1:0/"
// https_proxy="https://127.0.0.1:0/"
// ftp_proxy="ftp://127.0.0.1:0/"
// socks_proxy="socks://127.0.0.1:0/"

const (
	proxyAuto  = "auto"
	proxyHttp  = "http"
	proxyHttps = "https"
	proxyFtp   = "ftp"
	proxySocks = "socks"

	gsettingsIdProxy = "com.deepin.dde.proxy"
	gkeyProxyMethod  = "proxy-method"

	proxyMethodNone   = "none"
	proxyMethodManual = "manual"
	proxyMethodAuto   = "auto"

	gkeyAutoProxy = "auto-proxy"

	gkeyHttpProxy  = "http-proxy"
	gkeyHttpsProxy = "https-proxy"
	gkeyFtpProxy   = "ftp-proxy"
	gkeySocksProxy = "socks-proxy"
)

var (
	proxySettings  = gio.NewSettings(gsettingsIdProxy)
	proxyPrefixReg = regexp.MustCompile(`^.*?://(.*)$`)
)

func (m *Manager) GetProxyMethod() (proxyMethod string, err error) {
	proxyMethod = proxySettings.GetString(gkeyProxyMethod)
	logger.Info("GetProxyMethod", proxyMethod)
	return
}
func (m *Manager) SetProxyMethod(proxyMethod string) (err error) {
	logger.Info("SetProxyMethod", proxyMethod)
	err = checkProxyMethod(proxyMethod)
	if err != nil {
		return
	}
	ok := proxySettings.SetString(gkeyProxyMethod, proxyMethod)
	if !ok {
		err = fmt.Errorf("set proxy method through gsettings failed")
	}
	return
}
func checkProxyMethod(proxyMethod string) (err error) {
	switch proxyMethod {
	case proxyMethodNone, proxyMethodManual, proxyMethodAuto:
	default:
		err = fmt.Errorf("invalid proxy method", proxyMethod)
		logger.Error(err)
	}
	return
}

func (m *Manager) GetAutoProxy() (proxyAuto string, err error) {
	proxyAuto = proxySettings.GetString(gkeyAutoProxy)
	logger.Info("GetAutoProxy", proxyAuto)
	return
}
func (m *Manager) SetAutoProxy(proxyAuto string) (err error) {
	logger.Info("SetAutoProxy", proxyAuto)
	ok := proxySettings.SetString(gkeyAutoProxy, proxyAuto)
	if !ok {
		err = fmt.Errorf("set automatic proxy through gsettings failed", proxyAuto)
	}
	return
}

func (m *Manager) GetProxy(proxyType string) (addr, port string, err error) {
	proxy, err := doGetProxy(proxyType)
	if err != nil {
		return
	}
	if proxyPrefixReg.MatchString(proxy) {
		proxy = proxyPrefixReg.FindStringSubmatch(proxy)[1]
	}
	proxy = strings.TrimSuffix(proxy, "/")
	a := strings.Split(proxy, ":")
	if len(a) == 2 {
		addr = a[0]
		port = a[1]
	} else if len(a) == 3 {
		// with user and password
		addr = a[0] + ":" + a[1]
		port = a[2]
	}
	logger.Info("GetProxy:", proxyType, addr, port)
	return
}

// if address is empty, means to remove proxy setting
func (m *Manager) SetProxy(proxyType, addr, port string) (err error) {
	logger.Info("SetProxy:", proxyType, addr, port)
	var proxy string
	if len(addr) > 0 {
		if len(port) == 0 {
			err = fmt.Errorf("proxy port is empty")
			return
		}
		if !proxyPrefixReg.MatchString(proxy) {
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
	logger.Info("doGetProxy:", proxyType, proxy)
	return
}

func doSetProxy(proxyType, proxy string) (err error) {
	logger.Info("doSetProxy:", proxyType, proxy)
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
