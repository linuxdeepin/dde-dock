package calendar

import (
	"encoding/json"
	"strconv"
	"time"

	"pkg.deepin.io/lib/calendar/util"
)

func parseInt(str string) (int, error) {
	return strconv.Atoi(str)
}

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

func parseDate(str string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", str, time.Local)
	return t, err
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func iterDays(start, end time.Time, f func(t time.Time) error) error {
	t := start
	for {
		err := f(t)
		if err != nil {
			return err
		}
		if t.Equal(end) {
			break
		}
		t = t.AddDate(0, 0, 1)
	}
	return nil
}

func getYearTimeRange(year int) TimeRange {
	start := newTimeYMDHM(year, 1, 1, 0, 0)
	end := newTimeYMDHM(year, 12, 31, 0, 0)
	return TimeRange{start, end}
}

const (
	modeDay   = 0
	modeWeek  = 1
	modeMonth = 2
	modeYear  = 3
)

func getMondayDate(t time.Time) time.Time {
	nDays := int(t.Weekday()) - 1
	if t.Weekday() == time.Sunday {
		nDays = 6
	}
	return t.AddDate(0, 0, -nDays)
}

func getMonthTimeRange(year int, month int) TimeRange {
	start := newTimeYMDHM(year, time.Month(month), 1, 0, 0)
	days := util.GetSolarMonthDays(year, month)
	end := start.AddDate(0, 0, days-1)
	return TimeRange{start, end}
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
