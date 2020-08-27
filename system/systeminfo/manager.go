package systeminfo

import (
	"encoding/xml"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName = "com.deepin.system.SystemInfo"
	dbusPath        = "/com/deepin/system/SystemInfo"
	dbusInterface   = dbusServiceName

	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
	TB = 1 << 40
	EB = 1 << 50
)

type Manager struct {
	service         *dbusutil.Service
	PropsMu         sync.RWMutex
	MemorySize      uint64
	MemorySizeHuman string
	CurrentSpeed    uint64
}

type lshwXmlList struct {
	Items []lshwXmlNode `xml:"node"`
}

type lshwXmlNode struct {
	Description string `xml:"description"`
	Size        uint64 `xml:"size"`
}

func formatFileSize(fileSize uint64) (size string) {
	if fileSize < KB {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < MB {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(KB))
	} else if fileSize < GB {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(MB))
	} else if fileSize < TB {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(GB))
	} else if fileSize < EB {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(TB))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(EB))
	}
}

func (m *Manager) GetInterfaceName() string {
	return dbusInterface
}

func NewManager(service *dbusutil.Service) *Manager {
	var m = &Manager{
		service: service,
	}
	return m
}

func runLshwMemory() (out []byte, err error) {
	cmd := exec.Command("lshw", "-c", "memory", "-sanitize", "-xml")
	out, err = cmd.Output()
	if err != nil {
		logger.Error(err)
		return out, err
	}
	return out, err
}

func parseXml(bytes []byte) (result lshwXmlNode, err error) {
	logger.Debug("ParseXml bytes: ", string(bytes))
	var list lshwXmlList
	err = xml.Unmarshal(bytes, &list)
	if err != nil {
		logger.Error(err)
		return result, err
	}
	len := len(list.Items)
	for i := 0; i < len; i++ {
		data := list.Items[i]
		logger.Debug("Description : ", data.Description, " , size : ", data.Size)
		if strings.ToLower(data.Description) == "system memory" {
			result = data
		}
	}
	return result, err
}

func (m *Manager) setMemorySize(value uint64) {
	m.MemorySize = value
	m.service.EmitPropertyChanged(m, "MemorySize", m.MemorySize)
}

func (m *Manager) setMemorySizeHuman(value string) {
	m.MemorySizeHuman = value
	m.service.EmitPropertyChanged(m, "MemorySizeHuman", m.MemorySizeHuman)
}

func (m *Manager) calculateMemoryViaLshw() error {
	cmdOutBuf, err := runLshwMemory()
	if err != nil {
		logger.Error(err)
		return err
	}
	ret, err1 := parseXml(cmdOutBuf)
	if err1 != nil {
		logger.Error(err1)
		return err1
	}
	memory := formatFileSize(ret.Size)
	m.PropsMu.Lock()
	//set property value
	m.setMemorySize(ret.Size)
	m.setMemorySizeHuman(memory)
	m.PropsMu.Unlock()
	logger.Debug("system memory : ", ret.Size)
	return nil
}

func GetCurrentSpeed(systemBit int) (uint64, error) {
	ret, err := getCurrentSpeed(systemBit)
	return ret, err
}

func getCurrentSpeed(systemBit int) (uint64, error) {
	var ret uint64 = 0
	cmdOutBuf, err := runDmidecode()
	if err != nil {
		return ret, err
	}
	ret, err = parseCurrentSpeed(cmdOutBuf, systemBit)
	if err != nil {
		logger.Error(err)
		return ret, err
	}
	logger.Debug("GetCurrentSpeed :", ret)
	return ret, err
}

func runDmidecode() (string, error) {
	cmd := exec.Command("dmidecode", "-t", "processor")
	out, err := cmd.Output()
	if err != nil {
		logger.Error(err)
	}
	return string(out), err
}

//From string parse "Current Speed"
func parseCurrentSpeed(bytes string, systemBit int) (result uint64, err error) {
	logger.Debug("parseCurrentSpeed data: ", bytes)
	lines := strings.Split(bytes, "\n")
	for _, line := range lines {
		if !strings.Contains(line, "Current Speed:") {
			continue
		}
		items := strings.Split(line, "Current Speed:")
		ret := ""
		if len(items) == 2 {
			//Current Speed: 3200 MHz
			ret = items[1]
			value, err := strconv.ParseUint(strings.TrimSpace(filterUnNumber(ret)), 10, systemBit)
			if err != nil {
				logger.Error(err)
				return result, err
			}
			result = value
		}
		break
	}
	return result, err
}

//仅保留字符串中的数字
func filterUnNumber(value string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		logger.Fatal(err)
	}
	return reg.ReplaceAllString(value, "")
}

func (m *Manager) systemBit() string {
	output, err := exec.Command("/usr/bin/getconf", "LONG_BIT").Output()
	if err != nil {
		return "64"
	}

	v := strings.TrimRight(string(output), "\n")
	return v
}

func (m *Manager) setPropCurrentSpeed(value uint64) (changed bool) {
	if m.CurrentSpeed != value {
		m.CurrentSpeed = value
		err := m.emitPropChangedsetPropCurrentSpeed(value)
		if err != nil {
			logger.Warning("emitPropChangedsetPropCurrentSpeed err : ", err)
			changed = false
		} else {
			changed = true
		}
	}
	return changed
}

func (m *Manager) emitPropChangedsetPropCurrentSpeed(value uint64) error {
	return m.service.EmitPropertyChanged(m, "CurrentSpeed", value)
}
