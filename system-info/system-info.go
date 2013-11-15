package main

import (
	"dlib/dbus"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type SystemInfo struct {
	Version    int32  `access:"read"`
	Processor  string `access:"read"`
	MemoryCap  uint64 `access:"read"`
	SystemType int64  `access:"read"`
	DiskCap    uint64 `access:"read"`
}

func (sys *SystemInfo) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.SystemInfo",
		"/com/deepin/daemon/SystemInfo",
		"com.deepin.daemon.SystemInfo",
	}
}

func IsFileNotExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return true
	}

	return false
}

func GetVersion() int32 {
	if IsFileNotExist("/etc/lsb-release") {
		return 0
	}
	contents, err := ioutil.ReadFile("/etc/lsb-release")
	if err != nil {
		return 0
	}

	lines := strings.Split(string(contents), "\n")
	version := int32(0)
	for _, line := range lines {
		vars := strings.Split(line, "=")
		l := len(vars)
		if l < 2 {
			break
		}
		if vars[0] == "DISTRIB_RELEASE" {
			num, _ := strconv.ParseUint(vars[1], 10, 64)
			version = int32(num)
			break
		}
	}

	return version
}

func GetCpuInfo() string {
	if IsFileNotExist("/proc/cpuinfo") {
		return "Unknown"
	}
	contents, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return ""
	}

	info := ""
	cnt := 0
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		vars := strings.Split(line, ":")
		l := len(vars)
		if l < 2 {
			break
		}
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
	if IsFileNotExist("/proc/meminfo") {
		return 0
	}
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		l := len(fields)
		if l < 2 {
			break
		}
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

	if strings.Contains(string(out), "i386") ||
		strings.Contains(string(out), "i586") ||
		strings.Contains(string(out), "i686") {
		sysType = 32
	} else if strings.Contains(string(out), "x86_64") {
		sysType = 64
	} else {
		sysType = 0
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
