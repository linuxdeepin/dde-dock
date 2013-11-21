package main

import (
	"dlib/dbus"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"time"
)

const (
	_NTP_HOST           = "0.pool.ntp.org"
	_SET_DATE_TIME_DEST = "com.deepin.daemon.SetDateTime"
	_SET_DATE_TIME_PATH = "/com/deepin/daemon/SetDateTime"
	_SET_DATA_TIME_IFC  = "com.deepin.daemon.SetDateTime"
)

var (
	quitChan chan bool
)

type SetDateTime struct {
	NtpRunFlag bool
}

func (sdt *SetDateTime) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_SET_DATE_TIME_DEST,
		_SET_DATE_TIME_PATH,
		_SET_DATA_TIME_IFC,
	}
}

func (sdt *SetDateTime) SetCurrentDate(d string) {
	/* Date String Format: 2013-11-17 */
	if CountCharInString('-', d) != 2 {
		fmt.Println("date string format error")
		return
	}

	sysTime := time.Now()
	sysTmp := &sysTime
	_, tStr := GetDateTimeAny(sysTmp)
	cmd := exec.Command("date", "--set", d)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	sdt.SetCurrentTime(tStr)
}

func (sdt *SetDateTime) SetCurrentTime(t string) {
	/* Time String Format: 12:23:09 */
	if CountCharInString(':', t) != 2 {
		fmt.Println("time string format error")
		return
	}

	cmd := exec.Command("date", "+%T", "-s", t)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
}

func (sdt *SetDateTime) SyncNtpTime() {
	t, err := GetNtpNow()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, tStr := GetDateTimeAny(t)
	sdt.SetCurrentTime(tStr)
}

func (sdt *SetDateTime) SetNtpUsing(using bool) {
	if using {
		if sdt.NtpRunFlag {
			sdt.SyncNtpTime()
			return
		}

		sdt.NtpRunFlag = true
		go SetNtpThread(sdt)
	} else {
		sdt.NtpRunFlag = false
		quitChan <- true
	}
}

func SetNtpThread(sdt *SetDateTime) {
	for {
		sdt.SyncNtpTime()
		timer := time.NewTimer(time.Minute * 1)
		select {
		case <-timer.C:
		case <-quitChan:
			return
		}
	}
}

func NewSetDateTime() *SetDateTime {
	sdt := SetDateTime{}
	sdt.NtpRunFlag = false

	return &sdt
}

func CountCharInString(ch byte, str string) int {
	if l := len(str); l <= 0 {
		return 0
	}

	cnt := 0
	for i, _ := range str {
		if str[i] == ch {
			cnt++
		}
	}

	return cnt
}

func GetDateTimeAny(t *time.Time) (dStr, tStr string) {
	dStr += strconv.FormatInt(int64(t.Year()), 10) + "-" + strconv.FormatInt(int64(t.Month()), 10) + "-" + strconv.FormatInt(int64(t.Day()), 10)
	tStr += strconv.FormatInt(int64(t.Hour()), 10) + ":" + strconv.FormatInt(int64(t.Minute()), 10) + ":" + strconv.FormatInt(int64(t.Second()), 10)

	return dStr, tStr
}

func GetNtpNow() (*time.Time, error) {
	raddr, err := net.ResolveUDPAddr("udp", _NTP_HOST+":123")
	if err != nil {
		return nil, err
	}

	data := make([]byte, 48)
	data[0] = 3<<3 | 3

	con, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}

	defer con.Close()

	_, err = con.Write(data)
	if err != nil {
		return nil, err
	}

	con.SetDeadline(time.Now().Add(5 * time.Second))

	_, err = con.Read(data)
	if err != nil {
		return nil, err
	}

	var sec, frac uint64
	sec = uint64(data[43]) | uint64(data[42])<<8 | uint64(data[41])<<16 |
		uint64(data[40])<<24
	frac = uint64(data[47]) | uint64(data[46])<<8 | uint64(data[45])<<16 |
		uint64(data[44])<<24

	nsec := sec * 1e9
	nsec += (frac * 1e9) >> 32

	t := time.Date(1990, 1, 0, 0, 0, 0, 0, time.UTC).
		Add(time.Duration(nsec)).Local()

	return &t, nil
}

func main() {
	quitChan = make(chan bool)
	sdt := NewSetDateTime()
	err := dbus.InstallOnSystem(sdt)
	if err != nil {
		panic(err)
	}
	select {}
}
