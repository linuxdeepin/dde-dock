/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"fmt"
	dutils "pkg.deepin.io/lib/utils"
)

type secretProxyType []uint32

func (l *secretProxyType) Add(pid uint32) {
	if l.Last() == pid {
		return
	}

	l.delete(pid)
	*l = append(*l, pid)
}

func (l *secretProxyType) Last() uint32 {
	len := len(*l)
	if len == 0 {
		return 0
	}

	pid := (*l)[len-1]
	file := fmt.Sprintf("/proc/%v", pid)
	if !dutils.IsFileExist(file) {
		l.delete(pid)
		return l.Last()
	}
	return pid
}

func (l *secretProxyType) delete(pid uint32) {
	var ret secretProxyType
	for _, v := range *l {
		if v == pid {
			continue
		}
		ret = append(ret, v)
	}
	*l = ret
}
