/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import "fmt"
import "gir/gio-2.0"
import "strconv"
import "strings"

const (
	proxyTypeHttp  = "http"
	proxyTypeHttps = "https"
	proxyTypeFtp   = "ftp"
	proxyTypeSocks = "socks"

	// The Deepin proxy gsettings schemas use the same path with
	// org.gnome.system.proxy which is /system/proxy. So in fact they
	// control the same values, and we don't need to synchronize them
	// at all.
	gsettingsIdProxy = "com.deepin.wrap.gnome.system.proxy"

	gkeyProxyMode   = "mode"
	proxyModeNone   = "none"
	proxyModeManual = "manual"
	proxyModeAuto   = "auto"

	gkeyProxyAuto        = "autoconfig-url"
	gkeyProxyIgnoreHosts = "ignore-hosts"
	gkeyProxyHost        = "host"
	gkeyProxyPort        = "port"

	gchildProxyHttp  = "http"
	gchildProxyHttps = "https"
	gchildProxyFtp   = "ftp"
	gchildProxySocks = "socks"
)

var (
	proxySettings           *gio.Settings
	proxyChildSettingsHttp  *gio.Settings
	proxyChildSettingsHttps *gio.Settings
	proxyChildSettingsFtp   *gio.Settings
	proxyChildSettingsSocks *gio.Settings
)

func initProxyGsettings() {
	proxySettings = gio.NewSettings(gsettingsIdProxy)
	proxyChildSettingsHttp = proxySettings.GetChild(gchildProxyHttp)
	proxyChildSettingsHttps = proxySettings.GetChild(gchildProxyHttps)
	proxyChildSettingsFtp = proxySettings.GetChild(gchildProxyFtp)
	proxyChildSettingsSocks = proxySettings.GetChild(gchildProxySocks)
}

func getProxyChildSettings(proxyType string) (childSettings *gio.Settings, err error) {
	switch proxyType {
	case proxyTypeHttp:
		childSettings = proxyChildSettingsHttp
	case proxyTypeHttps:
		childSettings = proxyChildSettingsHttps
	case proxyTypeFtp:
		childSettings = proxyChildSettingsFtp
	case proxyTypeSocks:
		childSettings = proxyChildSettingsSocks
	default:
		err = fmt.Errorf("not a valid proxy type: %s", proxyType)
		logger.Error(err)
	}
	return
}

// GetProxyMethod get current proxy method, it would be "none",
// "manual" or "auto".
func (m *Manager) GetProxyMethod() (proxyMode string, err error) {
	proxyMode = proxySettings.GetString(gkeyProxyMode)
	logger.Info("GetProxyMethod", proxyMode)
	return
}
func (m *Manager) SetProxyMethod(proxyMode string) (err error) {
	logger.Info("SetProxyMethod", proxyMode)
	err = checkProxyMethod(proxyMode)
	if err != nil {
		return
	}

	// ignore if proxyModeNone already set
	currentMethod, _ := m.GetProxyMethod()
	if proxyMode == proxyModeNone && currentMethod == proxyModeNone {
		return
	}

	ok := proxySettings.SetString(gkeyProxyMode, proxyMode)
	if !ok {
		err = fmt.Errorf("set proxy method through gsettings failed")
		return
	}
	switch proxyMode {
	case proxyModeNone:
		notifyProxyDisabled()
	default:
		notifyProxyEnabled()
	}
	return
}
func checkProxyMethod(proxyMode string) (err error) {
	switch proxyMode {
	case proxyModeNone, proxyModeManual, proxyModeAuto:
	default:
		err = fmt.Errorf("invalid proxy method %s", proxyMode)
		logger.Error(err)
	}
	return
}

// GetAutoProxy get proxy PAC file URL for "auto" proxy mode, the
// value will keep there even the proxy mode is not "auto".
func (m *Manager) GetAutoProxy() (proxyAuto string, err error) {
	proxyAuto = proxySettings.GetString(gkeyProxyAuto)
	return
}
func (m *Manager) SetAutoProxy(proxyAuto string) (err error) {
	logger.Debug("set autoconfig-url for proxy", proxyAuto)
	ok := proxySettings.SetString(gkeyProxyAuto, proxyAuto)
	if !ok {
		err = fmt.Errorf("set autoconfig-url proxy through gsettings failed %s", proxyAuto)
		logger.Error(err)
	}
	return
}

// GetProxyIgnoreHosts get the ignored hosts for proxy network which
// is a string separated by ",".
func (m *Manager) GetProxyIgnoreHosts() (ignoreHosts string, err error) {
	array := proxySettings.GetStrv(gkeyProxyIgnoreHosts)
	ignoreHosts = strings.Join(array, ", ")
	return
}
func (m *Manager) SetProxyIgnoreHosts(ignoreHosts string) (err error) {
	logger.Debug("set ignore-hosts for proxy", ignoreHosts)
	ignoreHostsFixed := strings.Replace(ignoreHosts, " ", "", -1)
	array := strings.Split(ignoreHostsFixed, ",")
	ok := proxySettings.SetStrv(gkeyProxyIgnoreHosts, array)
	if !ok {
		err = fmt.Errorf("set automatic proxy through gsettings failed %s", ignoreHosts)
		logger.Error(err)
	}
	return
}

// GetProxy get the host and port for target proxy type.
func (m *Manager) GetProxy(proxyType string) (host, port string, err error) {
	childSettings, err := getProxyChildSettings(proxyType)
	if err != nil {
		return
	}
	host = childSettings.GetString(gkeyProxyHost)
	port = strconv.Itoa(int(childSettings.GetInt(gkeyProxyPort)))
	return
}

// SetProxy set host and port for target proxy type.
func (m *Manager) SetProxy(proxyType, host, port string) (err error) {
	logger.Debugf("Manager.SetProxy proxyType: %q, host: %q, port: %q", proxyType, host, port)

	if port == "" {
		port = "0"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		logger.Error(err)
		return
	}

	childSettings, err := getProxyChildSettings(proxyType)
	if err != nil {
		return
	}

	var ok bool
	if ok = childSettings.SetString(gkeyProxyHost, host); ok {
		ok = childSettings.SetInt(gkeyProxyPort, int32(portInt))
	}
	if !ok {
		err = fmt.Errorf("set proxy value to gsettings failed: %s, %s:%s", proxyType, host, port)
		logger.Error(err)
	}

	return
}
