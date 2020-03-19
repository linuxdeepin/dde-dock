package systeminfo

import (
	"fmt"
	"crypto/sha256"
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

const (
	uuidDelim = "+"
)

// Disk store disk info
type Disk struct {
	Name   string
	Model  string
	Serial string // if empty, use children uuid's sha256 replace
	Vendor string

	Size int64 // byte

	RootMounted bool
}

// DiskList multi disk
type DiskList []*Disk

type lsblkDevice struct {
	Name       string `json:"name"`
	Serial     string `json:"serial"`
	Type       string `json:"type"`
	Vendor     string `json:"vendor"`
	Model      string `json:"model"`
	UUID       string `json:"uuid"`
	MountPoint string `json:"mountpoint"`

	Size interface{} `json:"size"`

	Children lsblkDeviceList `json:"children"`
}
type lsblkDeviceList []*lsblkDevice

type lsblkOutput struct {
	Blockdevices lsblkDeviceList `json:"blockdevices"`
}

func GetDiskList() (DiskList, error) {
	out, err := execLsblk()
	if err != nil {
		return nil, err
	}
	return newDiskListFromOutput(out)
}

func (list DiskList) GetRoot() *Disk {
	for _, d := range list {
		if d.RootMounted {
			return d
		}
	}
	return nil
}

func newDiskListFromOutput(out []byte) (DiskList, error) {
	lsblk, err := parseLsblkOutput(out)
	if err != nil {
		return nil, err
	}

	var disks DiskList
	for _, info := range lsblk.Blockdevices {
		disks = append(disks, newDiskFromDevice(info))
	}
	return disks, nil
}

func newDiskFromDevice(dev *lsblkDevice) *Disk {
	var info = Disk{
		Name:        dev.Name,
		Model:       dev.Model,
		Serial:      dev.Serial,
		Vendor:      dev.Vendor,
		RootMounted: dev.RootMounted(),
	}

	if v, ok := dev.Size.(string); ok {
		info.Size, _ = strconv.ParseInt(v, 10, 64)
	} else if v, ok := dev.Size.(float64); ok {
		info.Size = int64(v)
	}

	if len(info.Serial) == 0 {
		// using children uuid list's sha256 as serial
		info.Serial = genSerialByUUIDList(dev.GetUUIDList())
	}
	return &info
}

func genSerialByUUIDList(list []string) string {
	if len(list) == 0 {
		return ""
	}

	str := strings.Join(list, uuidDelim)
	return sha256Sum([]byte(str))
}

func (dev *lsblkDevice) RootMounted() bool {
	for _, child := range dev.Children {
		if child.MountPoint == "/" {
			return true
		}
	}
	return false
}

func (dev *lsblkDevice) GetUUIDList() []string {
	var list []string
	for _, child := range dev.Children {
		if len(child.UUID) == 0 {
			continue
		}
		list = append(list, child.UUID)
	}
	return list
}

// SHA256Sum sum data by sha256
func sha256Sum(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func parseLsblkOutput(out []byte) (*lsblkOutput, error) {
	var info lsblkOutput
	err := json.Unmarshal(out, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func execLsblk() ([]byte, error) {
	lsblk := "lsblk -J -bno NAME,SERIAL,TYPE,SIZE,VENDOR,MODEL,MOUNTPOINT,UUID"
	return exec.Command("/bin/sh", "-c", lsblk).CombinedOutput()
}
