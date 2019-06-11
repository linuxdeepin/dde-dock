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
	"strings"
	"sync"

	"pkg.deepin.io/dde/daemon/grub_common"
)

const (
	grubMkconfigCmd = "grub-mkconfig"
	adjustThemeCmd  = "/usr/lib/deepin-api/adjust-grub-theme"
)

func init() {
	// force use LANG=en_US.UTF-8 to make lsb-release/os-probe support
	// Unicode characters
	// FIXME: keep same with the current system language settings
	os.Setenv("LANG", "en_US.UTF-8")
	os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
}

type modifyManager struct {
	g           *Grub2
	ch          chan modifyTask
	modifyTasks []modifyTask

	running       bool
	stateChangeCb func(running bool)

	mu sync.Mutex
}

func newModifyManager() *modifyManager {
	m := &modifyManager{
		ch: make(chan modifyTask),
	}
	return m
}

func (m *modifyManager) notifyStateChange() {
	if m.stateChangeCb != nil {
		m.stateChangeCb(m.running)
	}
}

func (m *modifyManager) loop() {
	for {
		select {
		case t, ok := <-m.ch:
			if !ok {
				return
			}
			m.mu.Lock()

			if m.running {
				m.modifyTasks = append(m.modifyTasks, t)
			} else {
				m.start(t)
			}

			m.mu.Unlock()
		}
	}
}

func (m *modifyManager) start(tasks ...modifyTask) {
	logger.Infof("modifyManager start")
	defer logger.Infof("modifyManager start return")

	params, _ := grub_common.LoadGrubParams()

	logger.Debug("modifyManager.start len(tasks):", len(tasks))
	var adjustTheme bool
	var adjustThemeLang string
	for _, task := range tasks {
		f := task.paramsModifyFunc
		if f != nil {
			f(params)
		}
		if task.adjustTheme {
			adjustTheme = true
			adjustThemeLang = task.adjustThemeLang
		}
	}
	err := writeGrubParams(params)
	if err != nil {
		logger.Warning("failed to write grub params:", err)
		return
	}

	logStart()
	m.running = true
	m.notifyStateChange()
	go m.update(adjustTheme, adjustThemeLang)
}

func (m *modifyManager) update(adjustTheme bool, adjustThemeLang string) {
	if adjustTheme {
		logJobStart(logJobAdjustTheme)
		var args []string
		if adjustThemeLang != "" {
			args = append(args, "-lang", adjustThemeLang)
		}
		cmd := exec.Command(adjustThemeCmd, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		logger.Debugf("$ %s %s", adjustThemeCmd, strings.Join(args, " "))
		err := cmd.Run()
		if err != nil {
			logger.Warning("failed to adjust theme:", err)
		}
		logJobEnd(logJobAdjustTheme, err)
		m.g.theme.emitSignalBackgroundChanged()
	}

	logJobStart(logJobMkConfig)
	err := runGrubMkconfig()
	if err != nil {
		logger.Warning("failed to make config:", err)
	}
	logJobEnd(logJobMkConfig, err)
	m.updateEnd()
}

func runGrubMkconfig() error {
	cmd := exec.Command(grubMkconfigCmd, "-o", grubScriptFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logger.Debugf("$ %s -o %s", grubMkconfigCmd, grubScriptFile)
	return cmd.Run()
}

func (m *modifyManager) updateEnd() {
	m.mu.Lock()

	logEnd()
	logger.Info("modifyManager update end")

	if len(m.modifyTasks) > 0 {
		m.start(m.modifyTasks...)
		m.modifyTasks = nil
	} else {
		// loop end
		m.running = false
		m.notifyStateChange()
	}

	m.mu.Unlock()
}

func (m *modifyManager) IsRunning() bool {
	m.mu.Lock()
	running := m.running
	m.mu.Unlock()
	return running
}
