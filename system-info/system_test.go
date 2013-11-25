package main

import (
	"testing"
	"fmt"
)

func TestSystem (t *testing.T) {
	sys := SystemInfo{}

	sys.Version = GetVersion()
	if sys.Version == int32(0) {
		t.Error("get version failed")
	}

	sys.Processor = GetCpuInfo()
	if sys.Processor == "" {
		t.Error("get cpu info failed")
	}

	sys.MemoryCap = GetMemoryCap()
	if sys.MemoryCap == uint64(0) {
		t.Error("get memory info failed")
	}

	sys.SystemType = GetSystemType()
	if sys.SystemType == int64(0) {
		t.Error("get system type failed")
	}

	sys.DiskCap = GetDiskCap()
	if sys.MemoryCap == uint64(0) {
		t.Error("get disk info failed")
	}

	fmt.Println("Version:", sys.Version)
	fmt.Println("CPU:", sys.Processor)
	fmt.Println("Memory:", sys.MemoryCap)
	fmt.Println("System Type:", sys.SystemType)
	fmt.Println("Disk:", sys.DiskCap)
}
