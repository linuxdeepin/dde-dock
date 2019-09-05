package calendar

import "time"

type TimeRange struct {
	start time.Time
	end   time.Time
}

func getTimeRange(start, end time.Time) TimeRange {
	if end.Before(start) {
		start, end = end, start
	}
	return TimeRange{
		start: start,
		end:   end,
	}
}

func (r TimeRange) contains(r1 TimeRange) bool {
	// r.start <= r1.start
	// and r1.end <= r.end
	return !r.start.After(r1.start) && !r1.end.After(r.end)
}

func (r TimeRange) overlap(r1 TimeRange) bool {
	// 保证 r 在前，r1 在后
	if r.start.After(r1.start) {
		r, r1 = r1, r
	}

	// r1.start <= r.end
	return !r1.start.After(r.end)
}
