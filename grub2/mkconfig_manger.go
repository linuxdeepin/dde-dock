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

package grub2

import (
	"os"
	"os/exec"
	"sync"
)

const (
	grubMkconfigCmd = "grub-mkconfig"
)

func init() {
	// force use LANG=en_US.UTF-8 to make lsb-release/os-probe support
	// Unicode characters
	// FIXME: keep same with the current system language settings
	os.Setenv("LANG", "en_US.UTF-8")
	os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
}

type MkconfigManager struct {
	ch          chan ModifyFunc
	modifyFuncs []ModifyFunc

	running     bool
	stateChange func(running bool)

	mu sync.Mutex
}

func newMkconfigManager(ch chan ModifyFunc, stateChange func(bool)) *MkconfigManager {
	m := &MkconfigManager{
		ch:          ch,
		stateChange: stateChange,
	}
	return m
}

func (m *MkconfigManager) notifyStateChange() {
	if m.stateChange != nil {
		m.stateChange(m.running)
	}
}

func (m *MkconfigManager) loop() {
	for {
		select {
		case f, ok := <-m.ch:
			if !ok {
				return
			}
			logger.Debug("mkconfigManager.loop receive f")
			m.mu.Lock()

			if m.running {
				m.modifyFuncs = append(m.modifyFuncs, f)
			} else {
				m.start(f)
			}

			m.mu.Unlock()
		}
	}
}

func (m *MkconfigManager) start(funcs ...ModifyFunc) {
	logger.Infof("mkconfig start")
	defer logger.Infof("mkconfig start return")

	params, _ := loadGrubParams()

	logger.Debug("mkconfigManager.start len(funcs):", len(funcs))
	for _, fn := range funcs {
		fn(params)
	}
	paramsHash, err := writeGrubParams(params)
	if err != nil {
		logger.Warning("failed to write grub params:", err)
		return
	}

	cmd := exec.Command(grubMkconfigCmd, "-o", grubScriptFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logStart(paramsHash)
	err = cmd.Start()
	if err != nil {
		logger.Warning("start grubMkconfigCmd failed:", err)
		return
	}
	go m.wait(cmd)
	m.running = true
	m.notifyStateChange()
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

	if len(m.modifyFuncs) > 0 {
		m.start(m.modifyFuncs...)
		m.modifyFuncs = nil
	} else {
		// loop end
		m.running = false
		m.notifyStateChange()
	}

	m.mu.Unlock()
}

func (m *MkconfigManager) IsRunning() bool {
	m.mu.Lock()
	running := m.running
	m.mu.Unlock()
	return running
}
