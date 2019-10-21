package calendar

import (
	"testing"
	"time"

	libdate "github.com/rickb777/date"
	"github.com/stretchr/testify/assert"
)

func TestParseRemind(t *testing.T) {
	tests := []struct {
		t      time.Time
		remind string
		result time.Time
		hasErr bool
	}{
		{
			t:      newTimeYMDHM(2019, 8, 29, 0, 0),
			remind: "",
			hasErr: false,
			result: time.Time{},
		},

		{
			t:      newTimeYMDHM(2019, 8, 29, 0, 0),
			remind: "0",
			hasErr: false,
			result: newTimeYMDHM(2019, 8, 29, 0, 0),
		},

		{
			t:      newTimeYMDHM(2019, 8, 29, 2, 0),
			remind: "60", // 提前 60 min
			hasErr: false,
			result: newTimeYMDHM(2019, 8, 29, 1, 0),
		},

		{
			t:      newTimeYMDHM(2019, 8, 29, 2, 0),
			remind: "abc", // 错误的值
			hasErr: true,
			result: time.Time{},
		},

		{
			t:      newTimeYMDHM(2019, 8, 29, 0, 0),
			remind: "0;08:00", // 当天的08:00
			hasErr: false,
			result: newTimeYMDHM(2019, 8, 29, 8, 0),
		},

		{
			t:      newTimeYMDHM(2019, 8, 29, 0, 0),
			remind: "1;08:00", // 提前一天
			hasErr: false,
			result: newTimeYMDHM(2019, 8, 28, 8, 0),
		},

		{
			t:      newTimeYMDHM(2019, 8, 29, 0, 0),
			remind: "0;24:00", // 错误的时间
			hasErr: true,
			result: time.Time{},
		},

		{
			t:      newTimeYMDHM(2019, 8, 29, 0, 0),
			remind: "0;00:60", // 错误的时间
			hasErr: true,
			result: time.Time{},
		},
	}
	for idx, test := range tests {
		rt, err := parseRemind(test.t, test.remind)
		assert.Equal(t, test.result, rt, "test idx: %d", idx)
		if test.hasErr {
			assert.NotNil(t, err, "test idx: %d", idx)
		} else {
			assert.Nil(t, err, "test idx: %d", idx)
		}
	}
}

func TestGetRemindAdvanceDays(t *testing.T) {
	n, err := getRemindAdvanceDays("0;09:00")
	assert.Nil(t, err)
	assert.Equal(t, 0, n)
	n, err = getRemindAdvanceDays("1;09:00")
	assert.Nil(t, err)
	assert.Equal(t, 1, n)

	n, err = getRemindAdvanceDays("0")
	assert.Nil(t, err)
	assert.Equal(t, 0, n)

	n, err = getRemindAdvanceDays("2880")
	assert.Nil(t, err)
	assert.Equal(t, 2, n)
}

func TestTimeRangeContains(t *testing.T) {
	r := getTimeRange(newTimeYMDHM(2019, 1, 1, 0, 0),
		newTimeYMDHM(2019, 1, 1, 2, 0))
	r1 := getTimeRange(newTimeYMDHM(2019, 1, 1, 0, 0),
		newTimeYMDHM(2019, 1, 1, 1, 0))
	assert.True(t, r.contains(r1))
	assert.False(t, r1.contains(r))

	r = getTimeRange(newTimeYMDHM(2019, 1, 1, 1, 0),
		newTimeYMDHM(2019, 1, 1, 2, 0))
	r1 = getTimeRange(newTimeYMDHM(2019, 1, 1, 0, 0),
		newTimeYMDHM(2019, 1, 1, 1, 0))
	assert.False(t, r.contains(r1))
	assert.False(t, r1.contains(r))

	r = getTimeRange(newTimeYMDHM(2019, 1, 1, 1, 0),
		newTimeYMDHM(2019, 1, 1, 2, 0))
	r1 = getTimeRange(newTimeYMDHM(2019, 1, 1, 0, 0),
		newTimeYMDHM(2019, 1, 1, 3, 0))
	assert.False(t, r.contains(r1))
	assert.True(t, r1.contains(r))
}

func TestTimeRangeOverlap(t *testing.T) {
	r := getTimeRange(newTimeYMDHM(2019, 1, 1, 0, 0),
		newTimeYMDHM(2019, 1, 1, 2, 0))
	r1 := getTimeRange(newTimeYMDHM(2019, 1, 1, 0, 0),
		newTimeYMDHM(2019, 1, 1, 1, 0))
	assert.True(t, r.overlap(r1))
	assert.True(t, r1.overlap(r))

	r = getTimeRange(newTimeYMDHM(2019, 1, 1, 0, 0),
		newTimeYMDHM(2019, 1, 1, 1, 0))
	r1 = getTimeRange(newTimeYMDHM(2019, 1, 1, 2, 0),
		newTimeYMDHM(2019, 1, 1, 3, 0))
	assert.False(t, r.overlap(r1))
	assert.False(t, r1.overlap(r))

	r = getTimeRange(newTimeYMDHM(2019, 1, 1, 0, 0),
		newTimeYMDHMS(2019, 1, 1, 23, 59, 59))
	r1 = getTimeRange(newTimeYMDHM(2019, 1, 2, 0, 0),
		newTimeYMDHMS(2019, 1, 2, 23, 59, 59))
	assert.False(t, r.overlap(r1))
	assert.False(t, r1.overlap(r))
}

func TestBetween(t *testing.T) {
	// 单日任务，无重复
	job := &Job{
		Start: newTimeYMDHM(2019, 9, 1, 9, 0),
		End:   newTimeYMDHM(2019, 9, 1, 10, 0),
	}
	startDate := libdate.New(2019, 9, 1)
	endDate := libdate.New(2019, 9, 10)
	jobTimes, err := job.between(startDate, endDate)
	assert.Nil(t, err)
	assert.Len(t, jobTimes, 1)
	assert.Equal(t, jobTime{start: newTimeYMDHM(2019, 9, 1, 9, 0)}, jobTimes[0])

	startDate = libdate.New(2019, 8, 1)
	endDate = libdate.New(2019, 8, 31)
	jobTimes, err = job.between(startDate, endDate)
	assert.Nil(t, err)
	assert.Len(t, jobTimes, 0)

	startDate = libdate.New(2019, 9, 2)
	endDate = libdate.New(2019, 9, 31)
	jobTimes, err = job.between(startDate, endDate)
	assert.Nil(t, err)
	assert.Len(t, jobTimes, 0)

	// 单日任务，重复：每天
	job = &Job{
		Start: newTimeYMDHM(2019, 9, 1, 9, 0),
		End:   newTimeYMDHM(2019, 9, 1, 10, 0),
		RRule: "FREQ=DAILY",
	}
	startDate = libdate.New(2019, 9, 1)
	endDate = libdate.New(2019, 9, 10)
	jobTimes, err = job.between(startDate, endDate)
	assert.Nil(t, err)
	assert.Equal(t, len(jobTimes), 10)
	assert.Equal(t, jobTimes[0], jobTime{start: newTimeYMDHM(2019, 9, 1, 9, 0)})
	assert.Equal(t, jobTimes[1], jobTime{start: newTimeYMDHM(2019, 9, 2, 9, 0), recurID: 1})
	assert.Equal(t, jobTimes[9], jobTime{start: newTimeYMDHM(2019, 9, 10, 9, 0), recurID: 9})

	// 多日任务，10日， 无重复
	job = &Job{
		Start: newTimeYMDHM(2019, 9, 1, 9, 0),
		End:   newTimeYMDHM(2019, 9, 10, 9, 0),
	}
	startDate = libdate.New(2019, 9, 1)
	endDate = libdate.New(2019, 9, 12)
	jobTimes, err = job.between(startDate, endDate)
	assert.Nil(t, err)
	assert.Len(t, jobTimes, 1)
	assert.Equal(t, jobTimes[0], jobTime{start: newTimeYMDHM(2019, 9, 1, 9, 0)})

	startDate = libdate.New(2019, 9, 5)
	endDate = libdate.New(2019, 9, 12)
	jobTimes, err = job.between(startDate, endDate)
	assert.Nil(t, err)
	assert.Len(t, jobTimes, 1)
	assert.Equal(t, jobTimes[0], jobTime{start: newTimeYMDHM(2019, 9, 1, 9, 0)})

	startDate = libdate.New(2019, 8, 1)
	endDate = libdate.New(2019, 8, 31)
	jobTimes, err = job.between(startDate, endDate)
	assert.Nil(t, err)
	assert.Len(t, jobTimes, 0)

	startDate = libdate.New(2019, 9, 11)
	endDate = libdate.New(2019, 9, 30)
	jobTimes, err = job.between(startDate, endDate)
	assert.Nil(t, err)
	assert.Len(t, jobTimes, 0)
}

func TestGetBodyTimePart(t *testing.T) {
	today := libdate.New(2019, 10, 15)
	t1 := newTimeYMDHM(2019, 10, 15, 9, 0)
	assert.Equal(t, "Today 09:00", getBodyTimePart(today, false, t1, true))
	assert.Equal(t, "today 09:00", getBodyTimePart(today, false, t1, false))

	t1 = newTimeYMDHM(2019, 10, 16, 9, 0)
	assert.Equal(t, "Tomorrow 09:00", getBodyTimePart(today, false, t1, true))
	assert.Equal(t, "tomorrow 09:00", getBodyTimePart(today, false, t1, false))

	t1 = newTimeYMDHM(2019, 10, 17, 9, 0)
	assert.Equal(t, "10/17/19 09:00", getBodyTimePart(today, false, t1, true))
	assert.Equal(t, "10/17/19 09:00", getBodyTimePart(today, false, t1, false))

	t1 = newTimeYMDHM(2019, 10, 15, 0, 0)
	assert.Equal(t, "Today", getBodyTimePart(today, true, t1, true))
	assert.Equal(t, "today", getBodyTimePart(today, true, t1, false))

	t1 = newTimeYMDHM(2019, 10, 16, 9, 0)
	assert.Equal(t, "Tomorrow", getBodyTimePart(today, true, t1, true))
	assert.Equal(t, "tomorrow", getBodyTimePart(today, true, t1, false))

	t1 = newTimeYMDHM(2019, 10, 17, 0, 0)
	assert.Equal(t, "10/17/19", getBodyTimePart(today, true, t1, true))
	assert.Equal(t, "10/17/19", getBodyTimePart(today, true, t1, false))
}

func TestGetRemindJobBody(t *testing.T) {
	now := newTimeYMDHM(2019, 10, 15, 18, 59)

	tests := []struct {
		start  time.Time
		end    time.Time
		allDay bool
		result string
	}{
		{
			start:  newTimeYMDHM(2019, 10, 15, 9, 10),
			end:    newTimeYMDHM(2019, 10, 15, 10, 20),
			result: "Today 09:10 to 10:20  Job Title",
		},
		{
			start:  newTimeYMDHM(2019, 10, 15, 9, 10),
			end:    newTimeYMDHM(2019, 10, 16, 10, 20),
			result: "Today 09:10 to tomorrow 10:20  Job Title",
		},
		{
			start:  newTimeYMDHM(2019, 10, 17, 9, 10),
			end:    newTimeYMDHM(2019, 10, 17, 10, 20),
			result: "10/17/19 09:10 to 10:20  Job Title",
		},
		{
			start:  newTimeYMDHM(2019, 10, 17, 9, 10),
			end:    newTimeYMDHM(2019, 10, 18, 10, 20),
			result: "10/17/19 09:10 to 10/18/19 10:20  Job Title",
		},
		{
			start:  newTimeYMDHM(2019, 10, 15, 0, 0),
			end:    newTimeYMDHM(2019, 10, 15, 23, 59),
			allDay: true,
			result: "Today  Job Title",
		},
		{
			start:  newTimeYMDHM(2019, 10, 15, 0, 0),
			end:    newTimeYMDHM(2019, 10, 16, 23, 59),
			allDay: true,
			result: "Today to tomorrow  Job Title",
		},
		{
			start:  newTimeYMDHM(2019, 10, 17, 0, 0),
			end:    newTimeYMDHM(2019, 10, 18, 23, 59),
			allDay: true,
			result: "10/17/19 to 10/18/19  Job Title",
		},
	}

	for idx, testData := range tests {
		job := &JobJSON{
			Start:  testData.start,
			End:    testData.end,
			Title:  "Job Title",
			AllDay: testData.allDay,
		}
		assert.Equal(t, testData.result, job.getRemindBody(now),
			"idx %d", idx)
	}
}
