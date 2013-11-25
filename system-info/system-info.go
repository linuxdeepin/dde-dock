package main

import (
	"dbus-gen/udisks2"
	"dlib/dbus"
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

func GetDiskCap() (diskCap uint64) {
	driList := []dbus.ObjectPath{}
	obj := udisks2.GetObjectManager("/org/freedesktop/UDisks2")
	managers := obj.GetManagedObjects()

	for _, value := range managers {
		if _, ok := value["org.freedesktop.UDisks2.Block"]; ok {
			v := value["org.freedesktop.UDisks2.Block"]["Drive"]
			path := v.Value().(dbus.ObjectPath)
			if path != dbus.ObjectPath("/") {
				flag := false
				l := len(driList)
				for i := 0; i < l; i++ {
					if driList[i] == path {
						flag = true
						break
					}
				}
				if !flag {
					driList = append(driList, path)
				}
			}
		}
	}

	for _, driver := range driList {
		_, driExist := managers[driver]
		rm, _ := managers[driver]["org.freedesktop.UDisks2.Drive"]["Removable"]
		if driExist && !(rm.Value().(bool)) {
			size := managers[driver]["org.freedesktop.UDisks2.Drive"]["Size"]
			diskCap += size.Value().(uint64)
		}
	}

	return diskCap
}

func main() {
	sys := SystemInfo{}

	sys.Version = GetVersion()
	sys.Processor = GetCpuInfo()
	sys.MemoryCap = GetMemoryCap()
	sys.SystemType = GetSystemType()
	sys.DiskCap = GetDiskCap()

	dbus.InstallOnSession(&sys)

	select {}
}
