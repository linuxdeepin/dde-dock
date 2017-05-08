/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package users

import (
	"fmt"
	"os/exec"
)

func isStrInArray(str string, array []string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}

	return false
}

func doAction(cmd string, args []string) error {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Printf("[doAction] exec '%s' failed: %s, %v\n", cmd, string(out), err)
	}
	return err
}
