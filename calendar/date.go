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

func (d Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", d.Year, d.Month, d.Day)
}
