/*
 * Copyright (C) 2016 ~ 2020 Deepin Technology Co., Ltd.
 *
 * Author:     hubenchang <hubenchang@uniontech.com>
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

package power

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	globalCpuDirPath                = "/sys/devices/system/cpu"
	globalCpuFreqDirName            = "cpufreq"
	globalGovernorFileName          = "scaling_governor"
	globalAvailableGovernorFileName = "scaling_available_governors"
	globalBoostFilePath             = "/sys/devices/system/cpu/cpufreq/boost"
)

type CpuHandler struct {
	path               string
	availableGovernors map[string]bool
	governor           string
}

type CpuHandlers []CpuHandler

func NewCpuHandlers() *CpuHandlers {
	cpus := make(CpuHandlers, 0)
	dirs, err := ioutil.ReadDir(globalCpuDirPath)
	if err != nil {
		logger.Warning(err)
		return &cpus
	}

	pattern, _ := regexp.Compile(`cpu[0-9]+`)
	for _, dir := range dirs {
		dirNane := dir.Name()
		cpuPath := filepath.Join(globalCpuDirPath, dirNane)
		isMatch := pattern.MatchString(dirNane)
		if isMatch {
			logger.Debugf("append %s", cpuPath)
			freqPath := filepath.Join(cpuPath, globalCpuFreqDirName)
			cpu := CpuHandler{
				path: freqPath,
			}
			_, err = cpu.GetAvailableGovernors(true)
			if err != nil {
				logger.Warning(err)
			}
			_, err = cpu.GetGovernor(true)
			if err != nil {
				logger.Warning(err)
			}
			cpus = append(cpus, cpu)
		} else {
			logger.Debugf("skip %s", cpuPath)
		}
	}

	logger.Debugf("total %d cpus", len(cpus))
	return &cpus
}

func (cpu *CpuHandler) GetAvailableGovernors(force bool) (map[string]bool, error) {
	if force {
		governors := make(map[string]bool)
		data, err := ioutil.ReadFile(filepath.Join(cpu.path, globalAvailableGovernorFileName))
		if err != nil {
			logger.Warning(err)
			return governors, err
		}

		for _, g := range strings.Split(string(data), " ") {
			g = strings.TrimSpace(g)
			governors[g] = true
		}
		cpu.availableGovernors = governors
	}

	return cpu.availableGovernors, nil
}

func (cpu *CpuHandler) GetGovernor(force bool) (string, error) {
	if force {
		data, err := ioutil.ReadFile(filepath.Join(cpu.path, globalGovernorFileName))
		if err != nil {
			logger.Warning(err)
			return "", err
		}
		cpu.governor = strings.TrimSpace(string(data))
	}

	return cpu.governor, nil
}

func (cpu *CpuHandler) SetGovernor(governor string) error {
	_, ok := cpu.availableGovernors[governor]
	if ok {
		err := ioutil.WriteFile(filepath.Join(cpu.path, globalGovernorFileName), []byte(governor), 0644)
		if err != nil {
			logger.Warning(err)
			return err
		}
		return nil
	} else {
		logger.Warningf("governor %q is unavailable.", governor)
		return fmt.Errorf("governor %q is unavailable.", governor)
	}
}

func (cpus *CpuHandlers) GetAvailableGovernors() (map[string]bool, error) {
	if len(*cpus) < 1 {
		return nil, fmt.Errorf("cannot find cpu files")
	}

	// 理论上应该都是一样的，但是这里求交集
	availableGovernors, _ := (*cpus)[0].GetAvailableGovernors(false)
	for i := 1; i < len(*cpus); i++ {
		buff := make(map[string]bool)
		available, _ := (*cpus)[i].GetAvailableGovernors(false)
		for key := range availableGovernors {
			_, ok := available[key]
			if ok {
				buff[key] = true
			}
		}
		availableGovernors = buff
	}

	return availableGovernors, nil
}

func (cpus *CpuHandlers) GetGovernor() (string, error) {
	if len(*cpus) < 1 {
		return "", fmt.Errorf("cannot find cpu files")
	}

	// 理论上应该都是一样的，但是这里确认一遍
	governor, _ := (*cpus)[0].GetGovernor(true)
	for _, cpu := range *cpus {
		temp, _ := cpu.GetGovernor(true)
		if governor != temp {
			logger.Warning("governors are not same")
			return "", fmt.Errorf("governors are not same")
		}
	}
	return governor, nil
}

func (cpus *CpuHandlers) SetGovernor(governor string) error {
	for _, cpu := range *cpus {
		err := cpu.SetGovernor(governor)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cpus *CpuHandlers) IsBoostFileExist() bool {
	_, err := os.Lstat(globalBoostFilePath)
	return err == nil
}

func (cpus *CpuHandlers) SetBoostEnabled(enabled bool) error {
	var err error
	if enabled {
		err = ioutil.WriteFile(globalBoostFilePath, []byte("0"), 0644)
	} else {
		err = ioutil.WriteFile(globalBoostFilePath, []byte("1"), 0644)
	}
	return err
}

func (cpus *CpuHandlers) GetBoostEnabled() (bool, error) {
	data, err := ioutil.ReadFile(globalBoostFilePath)
	return strings.TrimSpace(string(data)) == "0", err
}
