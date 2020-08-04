package systeminfo

import (
	"encoding/xml"
	"fmt"
	"os/exec"
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
	err := m.service.EmitPropertyChanged(m, "MemorySize", m.MemorySize)
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Manager) setMemorySizeHuman(value string) {
	m.MemorySizeHuman = value
	err := m.service.EmitPropertyChanged(m, "MemorySizeHuman", m.MemorySizeHuman)
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Manager) calculateMemoryViaLshw() {
	cmdOutBuf, err := runLshwMemory()
	if err != nil {
		logger.Error(err)
		return
	}
	ret, error := parseXml(cmdOutBuf)
	if error != nil {
		logger.Error(error)
		return
	}
	memory := formatFileSize(ret.Size)
	m.PropsMu.Lock()
	//set property value
	m.setMemorySize(ret.Size)
	m.setMemorySizeHuman(memory)
	m.PropsMu.Unlock()
	logger.Debug("System Memory : ", ret.Size)
}
