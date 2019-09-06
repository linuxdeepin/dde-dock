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
	job := &Job{
		Start: newTimeYMDHM(2019, 9, 1, 9, 0),
		End:   newTimeYMDHM(2019, 9, 1, 10, 0),
		RRule: "FREQ=DAILY",
	}
	startDate := libdate.New(2019, 9, 1)
	endDate := libdate.New(2019, 9, 30)
	timeCounts, err := job.between(startDate, endDate)
	assert.Nil(t, err)
	//t.Log(timeCounts)
	for _, timeCount := range timeCounts {
		t.Logf("%s %d\n", timeCount.start, timeCount.recurID)
	}
}
