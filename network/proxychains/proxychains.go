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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/xdg/basedir"
)

var logger *log.Logger

func SetLogger(l *log.Logger) {
	logger = l
}

type Manager struct {
	Type     string
	IP       string
	Port     uint32
	User     string
	Password string

	jsonFile string
	confFile string
	mu       sync.Mutex
}

func NewManager() *Manager {
	cfgDir := basedir.GetUserConfigDir()
	jsonFile := filepath.Join(cfgDir, "deepin", "proxychains.json")
	confFile := filepath.Join(cfgDir, "deepin", "proxychains.conf")
	m := &Manager{
		jsonFile: jsonFile,
		confFile: confFile,
	}
	go m.init()
	return m
}

const (
	dbusObjPath   = "/com/deepin/daemon/Network/ProxyChains"
	dbusInterface = "com.deepin.daemon.Network.ProxyChains"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Network",
		ObjectPath: dbusObjPath,
		Interface:  dbusInterface,
	}
}

const defaultType = "http"

func (m *Manager) init() {
	cfg, err := loadConfig(m.jsonFile)
	logger.Debug("load proxychains config file:", m.jsonFile)
	if err != nil {
		logger.Warning("load proxychains config failed:", err)
		m.Type = defaultType
	} else {
		m.Type = cfg.Type
		m.IP = cfg.IP
		m.Port = cfg.Port
		m.User = cfg.User
		m.Password = cfg.Password
	}

	changed := m.fixConfig()
	logger.Debug("fixConfig changed:", changed)
	if changed {
		if err := m.saveConfig(); err != nil {
			logger.Warning("save config failed", err)
		}
	}

	if !m.checkConfig() {
		// config is invalid
		logger.Warning("config is invalid")
		if err := m.removeConf(); err != nil {
			logger.Warning("remove conf file failed:", err)
		}
	}
}

func (m *Manager) saveConfig() error {
	cfg := &Config{
		Type:     m.Type,
		IP:       m.IP,
		Port:     m.Port,
		User:     m.User,
		Password: m.Password,
	}
	return cfg.save(m.jsonFile)
}

func (m *Manager) notifyChange(prop string) {
	dbus.NotifyChange(m, prop)
}

func (m *Manager) fixConfig() bool {
	var changed bool
	if !validType(m.Type) {
		m.Type = defaultType
		changed = true
	}

	if m.IP != "" && !validIPv4(m.IP) {
		m.IP = ""
		changed = true
	}

	if !validUser(m.User) {
		m.User = ""
		changed = true
	}

	if !validPassword(m.Password) {
		m.Password = ""
		changed = true
	}

	return changed
}

func (m *Manager) checkConfig() bool {
	if !validType(m.Type) {
		return false
	}

	if !validIPv4(m.IP) {
		return false
	}

	if !validUser(m.User) {
		return false
	}

	if !validPassword(m.Password) {
		return false
	}
	return true
}

type InvalidParamError struct {
	Param string
}

func (err InvalidParamError) Error() string {
	return fmt.Sprintf("invalid param %s", err.Param)
}

func (m *Manager) Set(type0, ip string, port uint32, user, password string) error {
	// allow type0 is empty
	if type0 == "" {
		type0 = defaultType
	}
	if !validType(type0) {
		return InvalidParamError{"Type"}
	}

	var disable bool
	if ip == "" && port == 0 {
		disable = true
	}

	if !disable && !validIPv4(ip) {
		return InvalidParamError{"IP"}
	}

	if !validUser(user) {
		return InvalidParamError{"User"}
	}

	if !validPassword(password) {
		return InvalidParamError{"Password"}
	}

	if (user == "" && password != "") || (user != "" && password == "") {
		return errors.New("user and password are not provided at the same time")
	}

	// all params are ok
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Type != type0 {
		m.Type = type0
		m.notifyChange("Type")
	}

	if m.IP != ip {
		m.IP = ip
		m.notifyChange("IP")
	}

	if m.Port != port {
		m.Port = port
		m.notifyChange("Port")
	}

	if m.User != user {
		m.User = user
		m.notifyChange("User")
	}

	if m.Password != password {
		m.Password = password
		m.notifyChange("Password")
	}

	err := m.saveConfig()
	if err != nil {
		return err
	}

	if disable {
		return m.removeConf()
	}

	// enable
	return m.writeConf()
}

func (m *Manager) writeConf() error {
	const head = `# Written by ` + dbusInterface + `
strict_chain
quiet_mode
proxy_dns
remote_dns_subnet 224
tcp_read_time_out 15000
tcp_connect_time_out 8000
localnet 127.0.0.0/255.0.0.0

[ProxyList]
`
	fh, err := os.Create(m.confFile)
	if err != nil {
		return err
	}
	_, err = fh.WriteString(head)
	if err != nil {
		return err
	}

	proxy := fmt.Sprintf("%s\t%s\t%v", m.Type, m.IP, m.Port)
	if m.User != "" && m.Password != "" {
		proxy += fmt.Sprintf("\t%s\t%s", m.User, m.Password)
	}
	_, err = fh.WriteString(proxy + "\n")
	if err != nil {
		return err
	}

	err = fh.Sync()
	if err != nil {
		return err
	}

	return fh.Close()
}

func (m *Manager) removeConf() error {
	err := os.Remove(m.confFile)
	if os.IsNotExist(err) {
		// ignore file not exist error
		return nil
	}
	return err
}
