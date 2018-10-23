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
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"pkg.deepin.io/lib/encoding/kv"
)

const (
	dataDir     = "/var/cache/deepin"
	logFile     = dataDir + "/grub2.log"
	logFileMode = 0644

	logJobMkConfig    = "mkconfig"
	logJobAdjustTheme = "adjustTheme"
)

func logStart() {
	content := fmt.Sprintf("start=%s\n", time.Now())
	err := ioutil.WriteFile(logFile, []byte(content), logFileMode)
	if err != nil {
		logger.Warning("logStart write failed:", err)
	}
}

func logAppendText(text string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, logFileMode)
	if err != nil {
		logger.Warning("logAppendText open failed:", err)
		return
	}
	defer f.Close()
	_, err = f.WriteString(text)
	if err != nil {
		logger.Warning("logAppendText write failed:", err)
	}
}

func logEnd() {
	logAppendText(fmt.Sprintf("end=%s\n", time.Now()))
}

func logJobStart(jobName string) {
	logAppendText(fmt.Sprintf("%sStart=%s\n", jobName, time.Now()))
}

func logJobEnd(jobName string, err error) {
	text := fmt.Sprintf("%sEnd=%s\n", jobName, time.Now())
	if err != nil {
		text += jobName + "Failed=1\n"
	}
	logAppendText(text)
}

type Log map[string]string

func loadLog() (Log, error) {
	f, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	l := make(Log)
	reader := kv.NewReader(f)

	for {
		pair, err := reader.Read()
		if err != nil {
			break
		}
		l[pair.Key] = pair.Value
	}

	return l, nil
}

func (l Log) hasJob(jobName string) bool {
	_, ok := l[jobName+"Start"]
	return ok
}

func (l Log) isJobDone(jobName string) bool {
	_, ok := l[jobName+"End"]
	return ok
}
