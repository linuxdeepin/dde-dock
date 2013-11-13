package main

import (
	"bitbucket.org/jpoirier/cpu"
	"dlib/dbus"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

type SystemInfo struct {
	Version    string
	Processor  string
	MemoryCap  uint64
	SystemType string
	DiskCap    string
}

func (sys *SystemInfo) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.SystemInfo",
		"/com/deepin/daemon/SystemInfo",
		"com.deepin.daemon.SystemInfo",
	}
}

func GetCpuInfo() string {
	cpuInfo := ""

	cpuInfo += string(cpu.ProcessorFamily)
	cpuInfo += " x "
	cpuInfo += strconv.FormatInt(int64(cpu.MaxProcs), 10)

	info := strings.TrimLeft(cpuInfo, " ")
	fmt.Println(info)
	return info
}

func GetMemoryCap() uint64 {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}

	lines := strings.Split(string(contents), "\n")
	memCap := uint64(0)
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "MemTotal:" {
			size, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return 0
			}
			memCap = size
			break
		}
	}

	return memCap
}

func GetSystemType() string {
	cmd := exec.Command("arch")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	ts := strings.Split(string(out), "_")
	return string(ts[1])
}

func main() {
	sys := SystemInfo{}

	sys.Version = "2013"
	sys.Processor = GetCpuInfo()
	sys.MemoryCap = GetMemoryCap()
	sys.SystemType = GetSystemType()
	sys.DiskCap = "500G"

	dbus.InstallOnSession(&sys)
	fmt.Println(sys.Processor)
	select {}
}
