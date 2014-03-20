package main

import (
        "dbus/org/freedesktop/udisks2"
        "dlib/dbus"
        "dlib/logger"
        "io/ioutil"
        "os"
        "os/exec"
        "strconv"
        "strings"
)

type SystemInfo struct {
        Version    int32
        Processor  string
        MemoryCap  uint64
        SystemType int64
        DiskCap    uint64
}

const (
        _VERSION_ETC = "/etc/lsb-release"
        _VERSION_KEY = "DISTRIB_RELEASE"

        _PROC_CPU_INFO = "/proc/cpuinfo"
        _PROC_CPU_KEY  = "model name"

        _PROC_MEM_INFO = "/proc/meminfo"
        _PROC_MEM_KEY  = "MemTotal"
)

var (
        logObject = logger.NewLogger("daemon/system-info")
)

func IsFileNotExist(filename string) bool {
        _, err := os.Stat(filename)
        if os.IsNotExist(err) {
                return true
        }

        return false
}

func GetVersion() int32 {
        if IsFileNotExist(_VERSION_ETC) {
                return 0
        }
        contents, err := ioutil.ReadFile(_VERSION_ETC)
        if err != nil {
                logObject.Infof("Read File Failed In Get Version: %s",
                        err)
                return 0
        }

        lines := strings.Split(string(contents), "\n")
        version := int32(0)
        for _, line := range lines {
                if strings.Contains(line, _VERSION_KEY) {
                        vars := strings.Split(line, "=")
                        l := len(vars)
                        if l < 2 {
                                break
                        }
                        num, _ := strconv.ParseUint(vars[1], 10, 64)
                        version = int32(num)
                        break
                }
        }

        return version
}

func GetCpuInfo() string {
        if IsFileNotExist(_PROC_CPU_INFO) {
                return "Unknown"
        }
        contents, err := ioutil.ReadFile(_PROC_CPU_INFO)
        if err != nil {
                logObject.Infof("Read File Failed In Get CPU Info: %s",
                        err)
                return ""
        }

        info := ""
        cnt := 0
        lines := strings.Split(string(contents), "\n")
        for _, line := range lines {
                if strings.Contains(line, _PROC_CPU_KEY) {
                        vars := strings.Split(line, ":")
                        l := len(vars)
                        if l < 2 {
                                break
                        }
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
        if IsFileNotExist(_PROC_MEM_INFO) {
                return 0
        }
        contents, err := ioutil.ReadFile(_PROC_MEM_INFO)
        if err != nil {
                logObject.Infof("Read File Failed In Get Memory Cap: %s",
                        err)
                return 0
        }

        lines := strings.Split(string(contents), "\n")
        for _, line := range lines {
                if strings.Contains(line, _PROC_MEM_KEY) {
                        fields := strings.Fields(line)
                        l := len(fields)
                        if l < 2 {
                                break
                        }
                        memCap, _ = strconv.ParseUint(fields[1], 10, 64)
                        break
                }
        }

        return (memCap * 1024)
}

func GetSystemType() (sysType int64) {
        cmd := exec.Command("/bin/uname", "-m")
        out, err := cmd.Output()
        if err != nil {
                logObject.Infof("Exec 'uname -m' Failed In Get System Type: %s",
                        err)
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
        obj, err := udisks2.NewObjectManager("org.freedesktop.UDisks2", "/org/freedesktop/UDisks2")
        if err != nil {
                logObject.Infof("udisks2: New ObjectManager Failed:%v", err)
                return 0
        }
        managers, _ := obj.GetManagedObjects()

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

func NewSystemInfo() *SystemInfo {
        sys := &SystemInfo{}

        sys.Version = GetVersion()
        sys.Processor = GetCpuInfo()
        sys.MemoryCap = GetMemoryCap()
        sys.SystemType = GetSystemType()
        sys.DiskCap = GetDiskCap()

        return sys
}

func main() {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Fatal("recover error:", err)
                }
        }()
        logObject.SetRestartCommand("/usr/lib/deepin-daemon/system-info")

        sys := NewSystemInfo()
        err := dbus.InstallOnSystem(sys)
        if err != nil {
                panic(err)
        }
        dbus.DealWithUnhandledMessage()

        //select {}
        if err = dbus.Wait(); err != nil {
                logObject.Info("lost dbus session:", err)
                os.Exit(1)
        } else {
                os.Exit(0)
        }
}
