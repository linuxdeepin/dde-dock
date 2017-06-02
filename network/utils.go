/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/utils"
	"regexp"
	"strings"
)

const iwBin = "iw"

func isStringInArray(s string, list []string) bool {
	for _, i := range list {
		if i == s {
			return true
		}
	}
	return false
}

func stringArrayBut(list []string, ignoreList ...string) (newList []string) {
	for _, s := range list {
		if !isStringInArray(s, ignoreList) {
			newList = append(newList, s)
		}
	}
	return
}

func appendStrArrayUnique(a1 []string, a2 ...string) (a []string) {
	a = a1
	for _, s := range a2 {
		if !isStringInArray(s, a) {
			a = append(a, s)
		}
	}
	return
}

func isDBusPathInArray(path dbus.ObjectPath, pathList []dbus.ObjectPath) bool {
	for _, i := range pathList {
		if i == path {
			return true
		}
	}
	return false
}

func isInterfaceNil(v interface{}) bool {
	return utils.IsInterfaceNil(v)
}

func isInterfaceEmpty(v interface{}) bool {
	if isInterfaceNil(v) {
		return true
	}
	switch v.(type) {
	case [][]interface{}: // ipv6Addresses
		if vd, ok := v.([][]interface{}); ok {
			if len(vd) == 0 {
				return true
			}
		}
	}
	return false
}

func marshalJSON(v interface{}) (jsonStr string, err error) {
	b, err := json.Marshal(v)
	if err != nil {
		logger.Error(err)
		return
	}
	jsonStr = string(b)
	return
}

func unmarshalJSON(jsonStr string, v interface{}) (err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	if err != nil {
		logger.Error(err)
	}
	return
}

func isUint32ArrayEmpty(a []uint32) (empty bool) {
	empty = true
	for _, v := range a {
		if v != 0 {
			empty = false
			break
		}
	}
	return
}

// convert local path to uri, etc "/the/path" -> "file:///the/path"
func toUriPath(path string) (uriPath string) {
	return utils.EncodeURI(path, utils.SCHEME_FILE)
}

// convert uri to local path, etc "file:///the/path" -> "/the/path"
func toLocalPath(path string) (localPath string) {
	return utils.DecodeURI(path)
}

// convert local path to uri, etc "/the/path" -> "file:///the/path"
func toUriPathFor8021x(path string) (uriPath string) {
	// the uri for 8021x cert files is specially, we just need append
	// suffix "file://" for it
	if !utils.IsURI(path) {
		uriPath = "file://" + path
	} else {
		uriPath = path
	}
	return
}

// convert uri to local path, etc "file:///the/path" -> "/the/path"
func toLocalPathFor8021x(path string) (uriPath string) {
	// the uri for 8021x cert files is specially, we just need remove
	// suffix "file://" from it
	if utils.IsURI(path) {
		uriPath = strings.TrimPrefix(path, "file://")
	} else {
		uriPath = path
	}
	return
}

// byte array should end with null byte
func strToByteArrayPath(path string) (bytePath []byte) {
	bytePath = []byte(path)
	bytePath = append(bytePath, 0)
	return
}
func byteArrayToStrPath(bytePath []byte) (path string) {
	if len(bytePath) < 1 {
		return
	}
	path = string(bytePath[:len(bytePath)-1])
	return
}

// strToUuid convert any given string to md5, and then to uuid, for
// example, a device address string "00:12:34:56:ab:cd" will be
// converted to "086e214c-1f20-bca4-9816-c0a11c8c0e02"
func strToUuid(str string) (uuid string) {
	md5, _ := utils.SumStrMd5(str)
	return doStrToUuid(md5)
}
func doStrToUuid(str string) (uuid string) {
	str = strings.ToLower(str)
	for i := 0; i < len(str); i++ {
		if (str[i] >= '0' && str[i] <= '9') ||
			(str[i] >= 'a' && str[i] <= 'f') {
			uuid = uuid + string(str[i])
		}
	}
	if len(uuid) < 32 {
		misslen := 32 - len(uuid)
		uuid = strings.Repeat("0", misslen) + uuid
	}
	uuid = fmt.Sprintf("%s-%s-%s-%s-%s", uuid[0:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:32])
	return
}

// execute program and read or write to it stdin/stdout pipe
func execWithIO(name string, arg ...string) (process *os.Process, stdin io.WriteCloser, stdout, stderr io.ReadCloser, err error) {
	cmd := exec.Command(name, arg...)
	stdin, _ = cmd.StdinPipe()
	stdout, _ = cmd.StdoutPipe()
	stderr, _ = cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		return
	}
	go cmd.Wait()

	process = cmd.Process
	return
}

// FIXME: temporary solution, please use libnl instead later
func isWirelessDeviceSuportHotspot(ifc string) (support bool) {
	var stdout, stderr string
	var err error
	var phynum, modes string
	var submatches []string

	// fixed 'iw' not in PATH
	pathEnv := os.Getenv("PATH")
	defer os.Setenv("PATH", pathEnv)
	os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")

	// get phy number for the device
	stdout, stderr, err = utils.ExecAndWait(5, iwBin, "dev", ifc, "info")
	if len(stderr) > 0 || err != nil {
		logger.Warningf("looks %s not exists, just let %s support hotspot mode: %s %s", iwBin, ifc, err, stderr)
		support = true
		return
	}
	regPhyNum := regexp.MustCompile("wiphy *([0-9]+)")
	submatches = regPhyNum.FindStringSubmatch(stdout)
	if len(submatches) >= 2 {
		phynum = submatches[1]
	}

	// get all supported modes
	stdout, stderr, err = utils.ExecAndWait(5, iwBin, "phy", "phy"+phynum, "info")
	if len(stderr) > 0 || err != nil {
		logger.Error(iwBin, "phy", "phy"+phynum, "info:", err, stderr)
		return
	}
	regModes := regexp.MustCompile("(?ims)Supported interface modes:(.*)software interface modes")
	submatches = regModes.FindStringSubmatch(stdout)
	if len(submatches) >= 2 {
		modes = submatches[1]
	}

	// check if support hotspot mode
	if strings.Contains(modes, "AP") {
		support = true
	}

	return
}
