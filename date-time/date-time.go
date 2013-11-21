package main

import (
	"dbus-gen/datetime"
	"dlib"
	"dlib/dbus"
	"dlib/dbus/property"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"time"
)

const (
	_NTP_HOST       = "0.pool.ntp.org"
	_DATE_TIME_DEST = "com.deepin.daemon.DateAndTime"
	_DATE_TIME_PATH = "/com/deepin/daemon/DateAndTime"
	_DATA_TIME_IFC  = "com.deepin.daemon.DateAndTime"

	_DATE_TIME_SCHEMA = "com.deepin.dde.datetime"
)

var (
	dtGSettings *dlib.Settings
	busConn     *dbus.Conn
	ntpRunFlag  bool
	quitChan    chan bool
)

type DateTime struct {
	AutoSetTime     bool `access:"read"`
	TimeShowFormat  dbus.Property
	CurrentTimeZone string `access:"read"`
}

func (date *DateTime) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_DATE_TIME_DEST, _DATE_TIME_PATH, _DATA_TIME_IFC}
}

func (date *DateTime) SetCurrentDate(d string) {
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
	fmt.Println(tStr)
	date.SetCurrentTime(tStr)
}

func (date *DateTime) SetCurrentTime(t string) {
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

func (date *DateTime) SetTimeZone(zone string) {
	gdate := datetime.GetDateTimeMechanism("/")
	gdate.SetTimezone(zone)
	date.CurrentTimeZone = zone
}

func (date *DateTime) SyncNtpTime() {
	t, err := GetNtpNow()
	if err != nil {
		/*date.AutoSetTime = false*/
		fmt.Println(err)
		return
	}

	_, tStr := GetDateTimeAny(t)
	date.SetCurrentTime(tStr)
}

func (date *DateTime) SetNtpUsing(using bool) {
	if using {
		if ntpRunFlag {
			date.SyncNtpTime()
			return
		}

		date.AutoSetTime = true
		ntpRunFlag = true
		dtGSettings.SetBoolean("is-auto-set", true)
		go SetNtpThread(date)
	} else {
		ntpRunFlag = false
		date.AutoSetTime = false
		quitChan <- true
		dtGSettings.SetBoolean("is-auto-set", false)
	}
	dbus.NotifyChange(busConn, date, "AutoSetTime")
}

func SetNtpThread(date *DateTime) {
	for {
		date.SyncNtpTime()
		timer := time.NewTimer(time.Minute * 1)
		select {
		case <-timer.C:
		case <-quitChan:
			return
		}
	}
}

func NewDateAndTime() *DateTime {
	dt := DateTime{}
	dtGSettings = dlib.NewSettings(_DATE_TIME_SCHEMA)

	dt.TimeShowFormat = property.NewGSettingsPropertyFull(dtGSettings,
		"is-24hour", true, busConn, _DATE_TIME_DEST, _DATA_TIME_IFC,
		"TimeShowFormat")
	d := datetime.GetDateTimeMechanism("/")
	dt.CurrentTimeZone = d.GetTimezone()
	dt.AutoSetTime = dtGSettings.GetBoolean("is-auto-set")

	return &dt
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

	fmt.Printf("date: %s\ntime: %s\n", dStr, tStr)
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
	var err error
	ntpRunFlag = false
	quitChan = make(chan bool)
	busConn, err = dbus.SystemBus()
	if err != nil {
		panic(err)
	}
	date := NewDateAndTime()
	err = dbus.InstallOnAny(busConn, date)
	if err != nil {
		panic(err)
	}

	dtGSettings.Connect("changed::is-auto-set", func(s *dlib.Settings, name string) {
		fmt.Println("is-auto-set changed:", s.GetBoolean("is-auto-set"))
		date.SetNtpUsing(s.GetBoolean("is-auto-set"))
	})

	if date.AutoSetTime {
		date.SetNtpUsing(true)
	}
	dlib.StartLoop()
}
