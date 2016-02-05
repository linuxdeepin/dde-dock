/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package screenedge

import (
	"os/exec"
)

type edge struct {
	Command string
	Area    *areaRange
}

func (e *edge) ExecAction() {
	command := e.Command
	if len(command) == 0 {
		logger.Debug("command empty")
		return
	}
	logger.Debug("execute command :", command)
	go exec.Command("/bin/sh", "-c", command).Run()
}
