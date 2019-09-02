package calendar

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/stephens2424/rrule"

	"github.com/jinzhu/gorm"
)

type Job struct {
	gorm.Model

	Type        int    // 类型
	Title       string // 标题
	Description string // 描述

	AllDay bool // 全天

	Start time.Time // 开始
	End   time.Time // 结束

	RRule  string // 重复规则
	Remind string // 提醒

	RecurID int    `gorm:"-"`
	Ignore  string // 忽略，JSON

	remindTime time.Time
}

type JobJSON struct {
	ID          uint
	Type        int
	Title       string
	Description string
	AllDay      bool
	Start, End  time.Time
	RRule       string
	Remind      string
	RecurID     int
	Ignore      []time.Time
}

func (j *Job) toJobJSON() (*JobJSON, error) {
	if j == nil {
		return nil, nil
	}
	ignore, err := j.getIgnore()
	if err != nil {
		return nil, err
	}
	return &JobJSON{
		ID:          j.ID,
		Type:        j.Type,
		Title:       j.Title,
		Description: j.Description,
		AllDay:      j.AllDay,
		Start:       j.Start,
		End:         j.End,
		RRule:       j.RRule,
		Remind:      j.Remind,
		RecurID:     j.RecurID,
		Ignore:      ignore,
	}, nil
}

func (j *JobJSON) toJob() (*Job, error) {
	if j == nil {
		return nil, nil
	}

	ignore, err := toJson(j.Ignore)
	if err != nil {
		return nil, err
	}

	job := &Job{
		Type:        j.Type,
		Title:       j.Title,
		Description: j.Description,
		AllDay:      j.AllDay,
		Start:       j.Start,
		End:         j.End,
		RRule:       j.RRule,
		Remind:      j.Remind,
		RecurID:     j.RecurID,
		Ignore:      ignore,
	}
	job.ID = j.ID
	return job, nil
}

func (j *Job) validate() error {
	// TODO
	return nil
}

func (j *Job) getIgnore() (result []time.Time, err error) {
	if j.Ignore == "" {
		return nil, nil
	}
	err = fromJson(j.Ignore, &result)
	return
}

func timeSliceContains(timeSlice []time.Time, t time.Time) bool {
	for _, t0 := range timeSlice {
		if t.Equal(t0) {
			return true
		}
	}
	return false
}

func getDateJobs(allJobs []*Job, date Date) ([]*Job, error) {

	var result []*Job
	dateRange := date.toTimeRange()
	for _, job := range allJobs {
		r := TimeRange{job.Start, job.End}

		ignore, err := job.getIgnore()
		if err != nil {
			return nil, err
		}
		if !timeSliceContains(ignore, job.Start) {
			if r.overlap(dateRange) {
				result = append(result, job)
			}
		}

		rJobs := job.getRecurrenceJobs(dateRange)
		for _, rJob := range rJobs {
			if timeSliceContains(ignore, rJob.Start) {
				continue
			}
			result = append(result, rJob)
		}
	}
	return result, nil
}

var remindReg1 = regexp.MustCompile(`\d;\d\d:\d\d`)

// 提醒的提前时间最大为 7 天
func parseRemind(startTime time.Time, remind string) (t time.Time, err error) {
	if remind == "" {
		return
	}

	if remindReg1.MatchString(remind) {
		var nDays, hour, min int
		_, err = fmt.Sscanf(remind, "%d;%d:%d", &nDays, &hour, &min)
		if err != nil {
			return
		}

		if nDays < 0 || nDays > 7 {
			err = errors.New("invalid value")
			return
		}

		if hour < 0 || hour > 23 {
			err = errors.New("invalid value")
			return
		}

		if min < 0 || min > 59 {
			err = errors.New("invalid value")
			return
		}

		t = startTime.AddDate(0, 0, -nDays)
		t = setClock(t, Clock{Hour: hour, Minute: min})
		return
	}
	var nMinutes int
	nMinutes, err = parseInt(remind)
	if err != nil {
		return
	}
	if nMinutes < 0 || nMinutes > 60*24*7 {
		err = errors.New("invalid value")
		return
	}

	t = startTime.Add(-time.Minute * time.Duration(nMinutes))
	return
}

func (j *Job) getRecurrenceJobs(dateRange TimeRange) []*Job {
	if j.RRule == "" {
		return nil
	}

	if dateRange.start.Before(j.Start) {
		return nil
	}

	rule, err := rrule.ParseRRule(j.RRule)
	if err != nil {
		logger.Warningf("failed to parse rrule %q: %v", j.RRule, err)
		return nil
	}
	rule.Dtstart = j.Start
	iter := rule.Iterator()
	iter.Next()

	count := 0
	var result []*Job
	for {
		count++
		t := iter.Next()
		if t == nil {
			break
		}
		start := *t
		if dateRange.end.Before(start) {
			break
		}
		interval := start.Sub(j.Start)
		end := j.End.Add(interval)

		r := TimeRange{start, end}
		if r.overlap(dateRange) {
			result = append(result, j.clone(start, end, count))
		}
		if count == 2000 {
			break
		}
	}
	return result
}

func (j *Job) clone(start, end time.Time, recurID int) *Job {
	j1 := &Job{
		Type:        j.Type,
		Title:       j.Title,
		Description: j.Description,
		AllDay:      j.AllDay,
		Start:       start,
		End:         end,
		RRule:       j.RRule,
		Remind:      j.Remind,
		RecurID:     recurID,
	}
	j1.ID = j.ID
	return j1
}

func (j *Job) String() string {
	var buf bytes.Buffer

	idDesc := strconv.Itoa(int(j.ID))
	if j.RecurID != 0 {
		idDesc += "/" + strconv.Itoa(j.RecurID)
	}

	buf.WriteString(fmt.Sprintf("job [%s] title: %q\n", idDesc, j.Title))
	buf.WriteString(fmt.Sprintf("start: %s, end: %s\n",
		formatTime(j.Start), formatTime(j.End)))
	buf.WriteString("rrule: " + j.RRule + "\n")

	return buf.String()
}
