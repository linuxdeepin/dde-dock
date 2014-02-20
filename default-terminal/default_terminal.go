/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
        "dlib/gio-2.0"
        "dlib/logger"
        "os/exec"
        "strings"
)

const (
        DEFAULT_TERMINAL_ID = "com.deepin.desktop.default-applications.terminal"
        SCHEMA_KEY_EXEC     = "exec"
        SCHEMA_KEY_ARG      = "exec-arg"
)

func getDefaultTerminal() (string, []string) {
        settings := gio.NewSettings(DEFAULT_TERMINAL_ID)
        defer settings.Unref()

        cmdStr := settings.GetString(SCHEMA_KEY_EXEC)
        argStr := settings.GetString(SCHEMA_KEY_ARG)

        if len(cmdStr) <= 0 {
                cmdStr = "/usr/bin/x-terminal-emulator"
                argStr = ""
        }

        return cmdStr, strings.Split(argStr, " ")
}

func execCommand(cmdStr string, args []string) {
        err := exec.Command(cmdStr, args...).Run()
        if err != nil {
                logger.Println("Exec", cmdStr, " ", args, " failed:", err)
                panic(err)
        }
}

func openDefaultTerminal() {
        defer func() {
                if err := recover(); err != nil {
                        logger.Println("Recover Error:", err)
                }
        }()

        cmd, args := getDefaultTerminal()
        execCommand(cmd, args)
}
