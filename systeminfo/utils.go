package systeminfo

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

const (
	memKeyTotal = "MemTotal"
	memKeyDelim = ":"
)

func getMemoryFromFile(file string) (uint64, error) {
	ret, err := parseInfoFile(file, memKeyDelim)
	if err != nil {
		return 0, err
	}

	value, ok := ret[memKeyTotal]
	if !ok {
		return 0, fmt.Errorf("Can not find the key '%s'", memKeyTotal)
	}

	cap, err := strconv.ParseUint(strings.Split(value, " ")[0], 10, 64)
	if err != nil {
		return 0, err
	}

	return cap * 1024, nil
}

func getOSType() (int64, error) {
	arch, err := getOSArch()
	if err != nil {
		return 0, err
	}

	switch strings.ToLower(arch) {
	case "i386", "i586", "i686":
		return 32, nil
	case "x86_64", "alpha":
		return 64, nil
	}

	return 0, fmt.Errorf("Unknown architecture: %v", arch)
}

func getOSArch() (string, error) {
	out, err := exec.Command("/bin/sh", "-c", "uname -m").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func parseInfoFile(file, delim string) (map[string]string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var ret = make(map[string]string)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		array := strings.Split(line, delim)
		if len(array) != 2 {
			continue
		}

		ret[strings.TrimSpace(array[0])] = strings.TrimSpace(array[1])
	}

	return ret, nil
}
