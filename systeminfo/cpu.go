package systeminfo

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	cpuKeyDelim     = ":"
	cpuKeyProcessor = "processor"
	cpuKeyName      = "model name"
	cpuKeyCPU       = "cpu"
	cpuKeyMHz       = "CPU frequency [MHz]"
	cpuKeyActive    = "cpus active"
)

func getCPUInfo(file string) (string, error) {
	ret, err := parseInfoFile(file, cpuKeyDelim)
	if err != nil {
		return "", err
	}

	cpu := swCPUInfo(ret)
	if len(cpu) != 0 {
		return cpu, nil
	}

	name, err := getCPUName(cpuKeyName, ret)
	if err != nil {
		return "", err
	}

	number, _ := getCPUNumber(cpuKeyProcessor, ret)
	if number != 0 {
		name = fmt.Sprintf("%s x %v", name, number+1)
	}

	return name, nil
}

func swCPUInfo(ret map[string]string) string {
	cpu, err := getCPUName(cpuKeyCPU, ret)
	if err != nil {
		return ""
	}

	hz, err := getCPUHz(cpuKeyMHz, ret)
	if err == nil {
		cpu = fmt.Sprintf("%s %vGHz", cpu, hz)
	}

	number, _ := getCPUNumber(cpuKeyActive, ret)
	if number != 1 {
		cpu = fmt.Sprintf("%s x %v", cpu, number)
	}

	return cpu
}

func getCPUName(key string, ret map[string]string) (string, error) {
	value, ok := ret[key]
	if !ok {
		return "", fmt.Errorf("Can not find the key '%s'", key)
	}

	var name string
	array := strings.Split(value, " ")
	for i, v := range array {
		if len(v) == 0 {
			continue
		}
		name += v
		if i != len(array)-1 {
			name += " "
		}
	}

	return name, nil
}

func getCPUNumber(key string, ret map[string]string) (int, error) {
	value, ok := ret[key]
	if !ok {
		return 0, fmt.Errorf("Can not find the key '%s'", key)
	}

	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(number), nil
}

func getCPUHz(key string, ret map[string]string) (float64, error) {
	value, ok := ret[key]
	if !ok {
		return 0, fmt.Errorf("Can not find the key '%s'", key)
	}

	hz, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return hz / 1000, nil
}
