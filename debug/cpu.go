/*
 * Copyright (C) 2018 ~ 2018 Deepin Technology Co., Ltd.
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

package debug

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type cpuTimeInfo struct {
	utime  int64
	stime  int64
	cutime int64
	cstime int64
	nice   int64
	start  int64
	hertz  int64
}

func getCPUPercentage() (float64, error) {
	statFile := fmt.Sprintf("/proc/%d/stat", os.Getpid())
	contents, err := ioutil.ReadFile(statFile)
	if err != nil {
		logger.Warning("Failed to read contents:", err)
		return 0, err
	}

	list := strings.Split(strings.Split(string(contents), ") ")[1], " ")
	var info = &cpuTimeInfo{
		utime:  stoi64(list[11]),
		stime:  stoi64(list[12]),
		cutime: stoi64(list[13]),
		cstime: stoi64(list[14]),
		nice:   stoi64(list[16]),
		start:  stoi64(list[19]),
		hertz:  getHertz(),
	}

	return info.Percentage(), nil
}

func (info *cpuTimeInfo) Total() int64 {
	return info.utime + info.stime + info.nice +
		info.cutime + info.cstime
}

func (info *cpuTimeInfo) Percentage() float64 {
	uptime := getUptime()
	seconds := uptime - float64(info.start/info.hertz)
	return 100 * (float64(info.Total()/info.hertz) / seconds)
}

func getUptime() float64 {
	contents, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		logger.Warning("Failed to read uptime:", err)
		return 0
	}

	str := strings.Split(string(contents), " ")[0]
	v, _ := strconv.ParseFloat(str, 64)
	return v
}

var _hertz int64

func getHertz() int64 {
	if _hertz != 0 {
		return _hertz
	}

	outs, err := exec.Command("getconf", "CLK_TCK").CombinedOutput()
	if err != nil {
		logger.Warning("Failed to get hertz:", string(outs), err)
		return 100 // default? why?
	}
	_hertz = stoi64(strings.Split(string(outs), "\n")[0])
	return _hertz
}

func stoi64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}
