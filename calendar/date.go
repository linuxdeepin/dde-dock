package calendar

import (
	"fmt"
	"time"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func (d Date) isWorkday() bool {
	w := d.weekday()
	return time.Monday <= w && w <= time.Friday
}

func (d Date) weekday() time.Weekday {
	t := newTimeYMDHM(d.Year, d.Month, d.Day, 0, 0)
	return t.Weekday()
}

const maxNanoSecs = 999999999

func (d Date) toTimeRange() TimeRange {
	return TimeRange{
		start: newTimeYMDHM(d.Year, d.Month, d.Day, 0, 0),
		end:   time.Date(d.Year, d.Month, d.Day, 23, 59, 59, maxNanoSecs, time.Local),
	}
}

func timeToDate(t time.Time) (d Date) {
	d.Year, d.Month, d.Day = t.Year(), t.Month(), t.Day()
	return
}

func sameDate(t1, t2 time.Time) bool {
	d1 := timeToDate(t1)
	d2 := timeToDate(t2)
	return d1 == d2
}

func (d Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", d.Year, d.Month, d.Day)
}
