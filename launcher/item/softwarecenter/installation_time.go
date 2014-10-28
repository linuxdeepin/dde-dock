package softwarecenter

import (
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

var timeRegExp = regexp.MustCompile(`(\d+)-(\d+)-(\d+) (\d+):(\d+):(\d+)`)

func extraTimeInstalled(installatiomTime string) int64 {
	subMatch := timeRegExp.FindStringSubmatch(installatiomTime)
	if len(subMatch) > 1 {
		year, _ := strconv.Atoi(subMatch[1])
		month, _ := strconv.Atoi(subMatch[2])
		day, _ := strconv.Atoi(subMatch[3])
		hour, _ := strconv.Atoi(subMatch[4])
		min, _ := strconv.Atoi(subMatch[5])
		sec, _ := strconv.Atoi(subMatch[6])
		date := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
		return date.Unix()
	}

	return 0
}

func GetTimeInstalled(pkgName string) int64 {
	// TODO: apt-history is a shell function.
	exec.Command("apt-history", "install", "|")
	return extraTimeInstalled("")
}
