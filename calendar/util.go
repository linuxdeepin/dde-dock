package calendar

import (
	"time"

	"encoding/json"
	"pkg.deepin.io/lib/libc"
)

func newTimeYMDHM(y int, m time.Month, d int, h int, min int) time.Time {
	return time.Date(y, m, d, h, min, 0, 0, time.Local)
}

func newTimeYMDHMS(y int, m time.Month, d int, h int, min int, s int) time.Time {
	return time.Date(y, m, d, h, min, s, 0, time.Local)
}

func setClock(t1 time.Time, c Clock) time.Time {
	t := time.Date(t1.Year(), t1.Month(), t1.Day(),
		c.Hour, c.Minute, c.Second, t1.Nanosecond(), t1.Location())
	return t
}

type Clock struct {
	Hour   int
	Minute int
	Second int
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func fromJson(jStr string, v interface{}) error {
	return json.Unmarshal([]byte(jStr), v)
}

func toJson(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func cFormatTime(format string, t time.Time) string {
	tm := libc.NewTm(t)
	v := libc.Strftime(format, tm)
	return v
}
