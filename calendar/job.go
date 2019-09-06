package calendar

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	libdate "github.com/rickb777/date"
	"github.com/stephens2424/rrule"
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
	// j.End < j.Start
	if j.End.Before(j.Start) {
		return errors.New("job end time before start time")
	}

	_, err := rrule.ParseRRule(j.RRule)
	if err != nil {
		return fmt.Errorf("invalid RRule: %v", err)
	}

	_, err = parseRemind(j.Start, j.Remind)
	if err != nil {
		return fmt.Errorf("invalid Remind: %v", err)
	}

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

type jobTime struct {
	start   time.Time
	recurID int
}

func getJobsBetween(startDate, endDate libdate.Date, jobs []*Job, extend bool) (wraps []dateJobsWrap) {
	days := endDate.Sub(startDate)
	wraps = make([]dateJobsWrap, days+1)
	date := startDate
	for idx := range wraps {
		wraps[idx].date = date
		date = date.Add(1)
	}

	for _, job := range jobs {
		interval := job.End.Sub(job.Start)
		jobTimes, err := job.between(startDate, endDate)
		if err != nil {
			if logger != nil {
				logger.Warning(err)
			}
			continue
		}
		for _, jobTime := range jobTimes {
			var j *Job
			if jobTime.recurID == 0 {
				j = job
			} else {
				j = job.clone(jobTime.start, jobTime.start.Add(interval), jobTime.recurID)
			}
			d := libdate.NewAt(jobTime.start)
			idx := d.Sub(startDate)
			wraps[idx].jobs = append(wraps[idx].jobs, j)
		}
	}

	if !extend {
		return
	}
	for idx, wrap := range wraps {
		for _, job := range wrap.jobs {
			jStartDate := libdate.NewAt(job.Start)
			jEndDate := libdate.NewAt(job.End)

			days := int(jEndDate.Sub(jStartDate))

			for i := 0; i < days; i++ {
				tIdx := idx + i + 1
				if tIdx == len(wraps) {
					break
				}
				w := &wraps[tIdx]
				w.extendJobs = append(w.extendJobs, job)
			}
		}
	}

	return
}

const recurrenceLimit = 3650

func (j *Job) between(startDate, endDate libdate.Date) ([]jobTime, error) {
	jStartDate := libdate.NewAt(j.Start)
	if endDate.Before(jStartDate) {
		// endDate < jStartDate
		return nil, nil
	}

	ignore, err := j.getIgnore()
	if err != nil {
		return nil, err
	}
	// 此次满足条件 jStartDate <= endDate
	if j.RRule == "" {
		if !startDate.After(jStartDate) {
			// startDate <= jStartDate <= endDate
			if timeSliceContains(ignore, j.Start) {
				// ignore this job
				return nil, nil
			}
			return []jobTime{
				{start: j.Start},
			}, nil
		}

		return nil, nil
	}

	rule, err := rrule.ParseRRule(j.RRule)
	if err != nil {
		return nil, err
	}
	rule.Dtstart = j.Start
	iter := rule.Iterator()

	count := 0
	var result []jobTime

	for {
		if count == recurrenceLimit {
			break
		}
		t := iter.Next()
		if t == nil {
			break
		}
		start := *t
		d := libdate.NewAt(start)
		if endDate.Before(d) {
			// endDate < d
			break
		}

		if !startDate.After(d) &&
			!timeSliceContains(ignore, start) {
			// startDate <= d <= endDate and not ignored
			result = append(result, jobTime{start: start, recurID: count})
		}

		count++
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
	nMinutes, err = strconv.Atoi(remind)
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

func (j *Job) getRemindTime() (time.Time, error) {
	start := j.Start
	if j.AllDay {
		start = setClock(j.Start, Clock{})
	}

	return parseRemind(start, j.Remind)
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
