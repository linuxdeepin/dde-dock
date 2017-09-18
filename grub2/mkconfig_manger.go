/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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

package grub2

import (
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
)

const (
	grubMkconfigCmd = "grub-mkconfig"
	grubParamsFile  = "/etc/default/grub"
)

func init() {
	// force use LANG=en_US.UTF-8 to make lsb-release/os-probe support
	// Unicode characters
	// FIXME: keep same with the current system language settings
	os.Setenv("LANG", "en_US.UTF-8")
	os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
}

type MkconfigManager struct {
	current     *Config
	future      *Config
	running     bool
	stateChange func(running bool)

	mu sync.Mutex
}

func newMkconfigManager(stateChange func(bool)) *MkconfigManager {
	m := &MkconfigManager{
		stateChange: stateChange,
	}
	return m
}

func (m *MkconfigManager) notifyStateChange() {
	if m.stateChange != nil {
		m.stateChange(m.running)
	}
}

func (m *MkconfigManager) Change(c *Config) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		m.future = c
	} else {
		m.start(c)
	}
}

func (m *MkconfigManager) start(c *Config) {
	logger.Infof("mkconfig start %s", c)
	m.running = true
	m.notifyStateChange()
	m.future = nil
	m.current = c

	writeGrubParams(c)
	cmd := exec.Command(grubMkconfigCmd, "-o", grubScriptFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logStart(c)
	err := cmd.Start()
	if err != nil {
		logger.Warning("start grubMkconfigCmd failed:", err)
		m.running = false
		m.notifyStateChange()
		return
	}
	go m.wait(cmd)
}

func (m *MkconfigManager) wait(cmd *exec.Cmd) {
	err := cmd.Wait()
	if err != nil {
		// exit status > 0
		logMkconfigFailed()
		logger.Warning(err)
	}
	logEnd()
	logger.Info("mkconfig end")

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.future != nil && !m.future.Equal(m.current) {
		m.start(m.future)
	} else {
		// loop end
		m.running = false
		m.notifyStateChange()
	}
}

func writeGrubParams(c *Config) error {
	content := c.GetGrubParamsContent()
	return ioutil.WriteFile(grubParamsFile, content, 0644)
}
