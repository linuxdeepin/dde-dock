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

package iw

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"pkg.deepin.io/lib/strv"
	"strings"
)

type DeviceInfo struct {
	Name       string
	MacAddress string
	Ciphers    []string
	IFCModes   []string
	Commands   []string
}
type DeviceInfos []*DeviceInfo

func ListDeviceInfo() (DeviceInfos, error) {
	var envPath = os.Getenv("PATH")
	os.Setenv("PATH", "/sbin:"+envPath)
	defer os.Setenv("PATH", envPath)
	outputs, err := exec.Command("/bin/sh", "-c",
		"exec iw list").CombinedOutput()
	if err != nil {
		return nil, err
	}
	return parseIwOutputs(outputs), nil
}

func (infos DeviceInfos) ListMiracastDevice() DeviceInfos {
	var ret DeviceInfos
	for _, info := range infos {
		if !info.SupportedMiracast() {
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func (infos DeviceInfos) ListHotspotDevice() DeviceInfos {
	var ret DeviceInfos
	for _, info := range infos {
		if !info.SupportedHotspot() {
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func (infos DeviceInfos) Get(macAddress string) *DeviceInfo {
	for _, info := range infos {
		if strings.ToLower(info.MacAddress) == strings.ToLower(macAddress) {
			return info
		}
	}
	return nil
}

func (info *DeviceInfo) SupportedHotspot() bool {
	return strv.Strv(info.IFCModes).Contains("AP")
}

func (info *DeviceInfo) SupportedMiracast() bool {
	list := strv.Strv(info.IFCModes)
	return list.Contains("P2P-client") &&
		list.Contains("P2P-GO")
	// list.Contains("P2P-device")
}

func debugDeviceInfos() {
	infos, err := ListDeviceInfo()
	if err != nil {
		fmt.Println("Failed to list wireless devices:", err)
		return
	}

	for _, info := range infos {
		fmt.Println(info.Name)
		fmt.Println("\tMac Address\t:", info.MacAddress)
		fmt.Println("\tCiphers\t:", info.Ciphers)
		fmt.Println("\tInterface Modes\t:", info.IFCModes)
		fmt.Println("\tCommands\t:", info.Commands)
	}
}

func parseIwOutputs(contents []byte) DeviceInfos {
	lines := strings.Split(string(contents), "\n")
	length := len(lines)
	var infos DeviceInfos
	for i := 0; i < length; {
		line := lines[i]
		if len(line) == 0 {
			i += 1
			continue
		}

		line = strings.TrimSpace(line)
		if strings.Contains(line, "Wiphy phy") {
			infos = append(infos, new(DeviceInfo))
			name := strings.Split(line, "Wiphy ")[1]
			infos[len(infos)-1].Name = name
			infos[len(infos)-1].MacAddress = getMacAddressByFile(macAddressFile(name))
			i += 1
			continue
		}

		if strings.Contains(line, "Supported Ciphers:") {
			i, infos[len(infos)-1].Ciphers = getValues(i+1, &lines)
			continue
		}

		if strings.Contains(line, "Supported interface modes:") {
			i, infos[len(infos)-1].IFCModes = getValues(i+1, &lines)
			continue
		}

		if strings.Contains(line, "Supported commands:") {
			i, infos[len(infos)-1].Commands = getValues(i+1, &lines)
			continue
		}

		i += 1
	}
	return infos
}

func getValues(idx int, lines *[]string) (int, []string) {
	var values []string
	length := len(*lines)
	for ; idx < length; idx++ {
		value := strings.TrimSpace((*lines)[idx])
		if value[0] != '*' {
			break
		}
		values = append(values, strings.Split(value, "* ")[1])
	}
	return idx, values
}

func macAddressFile(name string) string {
	return "/sys/class/ieee80211/" + name + "/macaddress"
}

func getMacAddressByFile(file string) string {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(contents))
}
