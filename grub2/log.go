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
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"pkg.deepin.io/lib/encoding/kv"
)

const (
	grubScriptFile = "/boot/grub/grub.cfg"
	dataDir        = "/var/cache/deepin"
	logFile        = dataDir + "/grub2.log"
	logFileMode    = 0644
)

func getGrubScriptMD5sum() (string, error) {
	return getFileMD5sum(grubScriptFile)
}

// write text:
// start= now
// paramsMD5Sum=
func logStart(paramHash string) {
	content := fmt.Sprintf("%s=%s\n%s=%s\n", logKeyStart, time.Now(),
		logKeyParamsMD5Sum, paramHash)
	err := ioutil.WriteFile(logFile, []byte(content), logFileMode)
	if err != nil {
		logger.Warning("logStart write failed:", err)
	}
}

// append text:
// mkconfigFailed=1
func logMkconfigFailed() {
	logAppendText(logKeyMkconfigFailed + "=1\n")
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

// append text:
// scriptMD5Sum=
// end= now
func logEnd() {
	sum, err := getGrubScriptMD5sum()
	if err != nil {
		logger.Warning("logEnd: getGrubScriptMD5sum failed:", err)
		return
	}

	logAppendText(fmt.Sprintf("%s=%s\n%s=%s\n", logKeyScriptMD5Sum, sum,
		logKeyEnd, time.Now()))
}

type Log struct {
	hasStart       bool
	hasEnd         bool
	paramsMD5Sum   string
	scriptMD5Sum   string
	mkconfigFailed bool
}

const (
	logKeyStart          = "start"
	logKeyEnd            = "end"
	logKeyParamsMD5Sum   = "paramsMD5Sum"
	logKeyScriptMD5Sum   = "scriptMD5Sum"
	logKeyMkconfigFailed = "mkconfigFailed"
)

func loadLog() (*Log, error) {
	f, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dict := make(map[string]string)
	reader := kv.NewReader(f)

	for {
		pair, err := reader.Read()
		if err != nil {
			break
		}
		dict[pair.Key] = pair.Value
	}

	l := &Log{}

	if dict[logKeyStart] != "" {
		l.hasStart = true
	}
	if dict[logKeyEnd] != "" {
		l.hasEnd = true
	}

	l.paramsMD5Sum = dict[logKeyParamsMD5Sum]
	l.scriptMD5Sum = dict[logKeyScriptMD5Sum]

	if dict[logKeyMkconfigFailed] == "1" {
		l.mkconfigFailed = true
	}
	return l, nil
}

func (l *Log) Verify(paramsMD5Sum string) (ok bool, err error) {
	// check start and end
	if !l.hasStart || !l.hasEnd {
		return false, nil
	}

	if l.mkconfigFailed {
		return false, nil
	}

	// check configHash
	if paramsMD5Sum != l.paramsMD5Sum {
		return false, nil
	}

	// check scriptMD5sum
	scriptMD5sum, err := getGrubScriptMD5sum()
	if err != nil {
		return false, err
	}

	return scriptMD5sum == l.scriptMD5Sum, nil
}
