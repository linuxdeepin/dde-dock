/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package systeminfo

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

const (
	memKeyTotal   = "MemTotal"
	memKeyDelim   = ":"
	lscpuKeyDelim = ":"
)

func getMemoryFromFile(file string) (uint64, error) {
	ret, err := parseInfoFile(file, memKeyDelim)
	if err != nil {
		return 0, err
	}

	value, ok := ret[memKeyTotal]
	if !ok {
		return 0, fmt.Errorf("Can not find the key '%s'", memKeyTotal)
	}

	cap, err := strconv.ParseUint(strings.Split(value, " ")[0], 10, 64)
	if err != nil {
		return 0, err
	}

	return cap * 1024, nil
}

func systemBit() string {
	output, err := exec.Command("/usr/bin/getconf", "LONG_BIT").Output()
	if err != nil {
		return "64"
	}

	v := strings.TrimRight(string(output), "\n")
	return v
}

func runLscpu() (map[string]string, error) {
	cmd := exec.Command("lscpu")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	res := make(map[string]string, len(lines))
	for _, line := range lines {
		items := strings.Split(line, lscpuKeyDelim)
		if len(items) != 2 {
			continue
		}

		res[items[0]] = strings.TrimSpace(items[1])
	}

	return res, nil
}

func parseInfoFile(file, delim string) (map[string]string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var ret = make(map[string]string)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		array := strings.Split(line, delim)
		if len(array) != 2 {
			continue
		}

		ret[strings.TrimSpace(array[0])] = strings.TrimSpace(array[1])
	}

	return ret, nil
}
