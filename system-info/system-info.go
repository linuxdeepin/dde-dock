package main

import (
	"dlib/dbus"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

type SystemInfo struct {
	Version    int32	`access:"read"`
	Processor  string	`access:"read"`
	MemoryCap  uint64	`access:"read"`
	SystemType int64	`access:"read"`
	DiskCap    uint64	`access:"read"`
}

func (sys *SystemInfo) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.SystemInfo",
		"/com/deepin/daemon/SystemInfo",
		"com.deepin.daemon.SystemInfo",
	}
}

func GetVersion() int32 {
	contents, err := ioutil.ReadFile("/etc/lsb-release")
	if err != nil {
		return 0
	}

	lines := strings.Split(string(contents), "\n")
	version := int32(0)
	for _, line := range lines {
		vars := strings.Split(line, "=")
		if vars[0] == "DISTRIB_RELEASE" {
			num, _ := strconv.ParseUint(vars[1], 10, 64)
			version = int32(num)
			break
		}
	}

	return version
}

func GetCpuInfo() string {
	contents, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return ""
	}

	info := ""
	cnt := 0
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		vars := strings.Split(line, ":")
		if strings.Contains(vars[0], "model name") {
			cnt++
			if info == "" {
				info += vars[1]
			}
		}
	}
	info += " x "
	info += strconv.FormatInt(int64(cnt), 10)

	return strings.TrimSpace(info)
}

func GetMemoryCap() (memCap uint64) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "MemTotal:" {
			memCap, _ = strconv.ParseUint(fields[1], 10, 64)
			break
		}
	}

	return memCap
}

func GetSystemType() (sysType int64) {
	cmd := exec.Command("uname", "-m")
	out, err := cmd.Output()
	if err != nil {
		return int64(0)
	}

	t := strings.TrimSpace(string(out))
	switch t {
	case "i386", "i586", "i686":
		sysType = 32
	case "x86_64":
		sysType = 64
	}

	return sysType
}

func GetDiskCap() (disCcap uint64) {
	return uint64(512000)
}

func main() {
	sys := SystemInfo{}

	sys.Version = GetVersion()
	sys.Processor = GetCpuInfo()
	sys.MemoryCap = GetMemoryCap()
	sys.SystemType = GetSystemType()
	sys.DiskCap = GetDiskCap()

	dbus.InstallOnSession(&sys)

	fmt.Println("Version:", sys.Version)
	fmt.Println("CPU:", sys.Processor)
	fmt.Println("Memory:", sys.MemoryCap)
	fmt.Println("System Type:", sys.SystemType)
	fmt.Println("Disk:", sys.DiskCap)
	select {}
}
