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

package dock

import (
	"dlib/gio-2.0"
	"os"
)

const (
	envAutoProxy  = "auto_proxy"
	envHttpProxy  = "http_proxy"
	envHttpsProxy = "https_proxy"
	envFtpProxy   = "ftp_proxy"
	envSocksProxy = "socks_proxy"

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
	proxySettings = gio.NewSettings(gsettingsIdProxy)
)

func startProxy() {
	updateProxyEnvs()
	listenProxyGsettings()
}

func updateProxyEnvs() {
	proxyMethod := proxySettings.GetString(gkeyProxyMethod)
	switch proxyMethod {
	case proxyMethodNone:
		unsetEnv(envAutoProxy)
		unsetEnv(envHttpProxy)
		unsetEnv(envHttpsProxy)
		unsetEnv(envFtpProxy)
		unsetEnv(envSocksProxy)
	case proxyMethodAuto:
		autoProxy := proxySettings.GetString(gkeyAutoProxy)
		os.Setenv(envAutoProxy, autoProxy)
		unsetEnv(envHttpProxy)
		unsetEnv(envHttpsProxy)
		unsetEnv(envFtpProxy)
		unsetEnv(envSocksProxy)
	case proxyMethodManual:
		unsetEnv(envAutoProxy)
		httpProxy := proxySettings.GetString(gkeyHttpProxy)
		os.Setenv(envHttpProxy, httpProxy)
		logger.Debug("update proxy envs, httpProxy:", httpProxy)

		httpsProxy := proxySettings.GetString(gkeyHttpsProxy)
		os.Setenv(envHttpsProxy, httpsProxy)
		logger.Debug("update proxy envs, httpsProxy:", httpsProxy)

		ftpProxy := proxySettings.GetString(gkeyFtpProxy)
		os.Setenv(envFtpProxy, ftpProxy)
		logger.Debug("update proxy envs, ftpProxy:", ftpProxy)

		socksProxy := proxySettings.GetString(gkeySocksProxy)
		os.Setenv(envSocksProxy, socksProxy)
		logger.Debug("update proxy envs, socksProxy:", socksProxy)
	}
}

func listenProxyGsettings() {
	proxySettings.Connect("changed", func(s *gio.Settings, key string) {
		logger.Debug("proxy value in gsettings changed:", key, s.GetString(key))
		updateProxyEnvs()
	})
}
