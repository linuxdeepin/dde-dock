/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package systeminfo

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

const (
	memKeyTotal = "MemTotal"
	memKeyDelim = ":"
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
