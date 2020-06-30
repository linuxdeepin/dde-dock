package calendar

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	lunarcalendar "github.com/linuxdeepin/go-dbus-factory/com.deepin.api.lunarcalendar"
	libdate "github.com/rickb777/date"
	"github.com/teambition/rrule-go"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/gettext"
)

type Job struct {
	gorm.Model

	Type        int    // 类型
	Title       string // 标题
	Description string // 描述
	TitlePinyin string // 拼音

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

	remindLaterCount int
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
		TitlePinyin: createPinyin(j.Title),
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

func (j *Job) timeRange() TimeRange {
	return TimeRange{
		start: j.Start,
		end:   j.End,
	}
}

func (j *Job) validate() error {
	// j.End < j.Start
	if j.End.Before(j.Start) {
		return errors.New("job end time before start time")
	}

	var err error
	if j.RRule != "" {
		_, err = rrule.StrToROptionInLocation(j.RRule, time.Local)
		if err != nil {
			return fmt.Errorf("invalid RRule: %v", err)
		}
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

func (j *Job) setIgnore(ignore []time.Time) error {
	v, err := toJson(ignore)
	if err != nil {
		return err
	}
	j.Ignore = v
	return nil
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

func getJobsBetween(startDate, endDate libdate.Date, jobs []*Job, extend bool, queryKey string, needFestival bool) (wraps []dateJobsWrap) {
	days := endDate.Sub(startDate)
	wraps = make([]dateJobsWrap, days+1)
	date := startDate
	for idx := range wraps {
		wraps[idx].date = date
		date = date.Add(1)
	}

	if needFestival {
		fillFestivalJobs(startDate, endDate, queryKey, wraps)
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
			j = job.clone(jobTime.start, jobTime.start.Add(interval), jobTime.recurID)
			jStartDate := libdate.NewAt(jobTime.start)
			idx := int(jStartDate.Sub(startDate))

			if 0 <= idx && idx < len(wraps) {
				wraps[idx].jobs = append(wraps[idx].jobs, j)
			}

			if !extend {
				continue
			}
			// NOTE: extend 指把跨越多天的 job 给扩展开，比如开始日期为 2019-09-01 结束日期为 2019-09-02 的
			// job, 将会扩展出一个job 放入 2019-09-02 那天的 extendJobs 字段中。

			jEndDate := libdate.NewAt(j.End)
			days := int(jEndDate.Sub(jStartDate))
			for i := 0; i < days; i++ {
				tIdx := idx + i + 1
				if tIdx == len(wraps) {
					break
				}
				if tIdx < 0 {
					continue
				}

				w := &wraps[tIdx]
				w.extendJobs = append(w.extendJobs, j)
			}
		}
	}

	return
}

var festivalIdMap = make(map[string]uint)
var nextFestivalJobId uint = math.MaxUint32
var festivalIdMapMu sync.Mutex

func getFestivalJobId(name string) uint {
	festivalIdMapMu.Lock()
	defer festivalIdMapMu.Unlock()

	id, ok := festivalIdMap[name]
	if ok {
		return id
	}

	id = nextFestivalJobId
	festivalIdMap[name] = id
	nextFestivalJobId--
	return id
}

func fillFestivalJobs(startDate, endDate libdate.Date, queryKey string, wraps []dateJobsWrap) {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		logger.Warning(err)
		return
	}
	lc := lunarcalendar.NewLunarCalendar(sessionBus)
	dayFestivalSlice, err := lc.GetFestivalsInRange(0, startDate.String(), endDate.String())
	if err != nil {
		logger.Warning(err)
		return
	}

	if queryKey != "" {
		dayFestivalSlice = filterDayFestivalSlice(dayFestivalSlice, queryKey)
	}

	for _, dayFestival := range dayFestivalSlice {
		if len(dayFestival.Festivals) == 0 {
			continue
		}
		d := libdate.New(int(dayFestival.Year), time.Month(dayFestival.Month), int(dayFestival.Day))
		idx := int(d.Sub(startDate))
		//logger.Debugf("id: %d, d: %v, festivals: %v", idx, d, dayFestival.Festivals)
		if 0 <= idx && idx < len(wraps) {
			for _, festival := range dayFestival.Festivals {
				j := &Job{
					Title:  festival,
					Type:   JobTypeFestival,
					AllDay: true,
					Start:  newTimeYMDHM(d.Year(), d.Month(), d.Day(), 0, 0),
					End:    newTimeYMDHM(d.Year(), d.Month(), d.Day(), 23, 59),
					RRule:  "FREQ=YEARLY",
				}
				j.ID = getFestivalJobId(festival)
				wraps[idx].jobs = append(wraps[idx].jobs, j)
			}
		}
	}
}

func filterDayFestivalSlice(dayFestivalSlice []lunarcalendar.DayFestival, queryKey string) []lunarcalendar.DayFestival {
	var result []lunarcalendar.DayFestival
	for _, dayFestival := range dayFestivalSlice {
		var festivals []string
		for _, festival := range dayFestival.Festivals {
			if strings.Contains(festival, queryKey) {
				festivals = append(festivals, festival)
			}
			if canQueryByPinyin(queryKey) && pinyinMatch(festival, queryKey) {
				festivals = append(festivals, festival)
			}
		}
		if len(festivals) == 0 {
			continue
		}
		result = append(result, lunarcalendar.DayFestival{
			Year:      dayFestival.Year,
			Month:     dayFestival.Month,
			Day:       dayFestival.Day,
			Festivals: festivals,
		})
	}
	return result
}

const recurrenceLimit = 3650

func (j *Job) between(startDate, endDate libdate.Date) ([]jobTime, error) {
	jStartDate := libdate.NewAt(j.Start)
	jEndDate := libdate.NewAt(j.End)
	nDays := jEndDate.Sub(jStartDate)
	if endDate.Before(jStartDate) {
		// endDate < jStartDate
		return nil, nil
	}

	ignore, err := j.getIgnore()
	if err != nil {
		return nil, err
	}
	// 此处满足条件 jStartDate <= endDate
	if j.RRule == "" {
		if (dateRange{startDate, endDate}).overlap(dateRange{jStartDate, jEndDate}) {
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

	rOpt, err := rrule.StrToROptionInLocation(j.RRule, time.Local)
	if err != nil {
		return nil, err
	}
	rOpt.Dtstart = j.Start
	rule, err := rrule.NewRRule(*rOpt)
	if err != nil {
		return nil, err
	}
	next := rule.Iterator()

	count := 0
	var result []jobTime

	for {
		if count == recurrenceLimit {
			break
		}
		start, ok := next()
		if !ok {
			break
		}
		jStartDate := libdate.NewAt(start)
		if endDate.Before(jStartDate) {
			// endDate < jStartDate
			// jStartDate 会随着循环迭代而增加，所以应该 break。
			break
		}
		jEndDate := jStartDate.Add(nDays)

		if (dateRange{startDate, endDate}).overlap(dateRange{jStartDate, jEndDate}) {
			if timeSliceContains(ignore, start) {
				// ignore this job
			} else {
				result = append(result, jobTime{start: start, recurID: count})
			}
		}

		count++
	}

	return result, nil
}

// 返回提醒设置提前几天
func getRemindAdvanceDays(remind string) (int, error) {
	var err error
	if remindReg1.MatchString(remind) {
		var nDays, hour, min int
		_, err = fmt.Sscanf(remind, "%d;%d:%d", &nDays, &hour, &min)
		if err != nil {
			return 0, err
		}

		return nDays, nil
	}
	var nMinutes int
	nMinutes, err = strconv.Atoi(remind)
	if err != nil {
		return 0, err
	}

	return int(float64(nMinutes) / float64(24*60)), nil
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

func getBodyTimePart(today libdate.Date, allDay bool, t time.Time, start bool) (msgStr string) {
	title := func(v string) string {
		return v
	}
	if start {
		title = strings.Title
	}

	date := libdate.NewAt(t)
	diff := date.Sub(today)
	if allDay {
		msgStr = cFormatTime("%x", t)
		if diff == 0 {
			msgStr = title(gettext.Tr("today"))
		} else if diff == 1 {
			msgStr = title(gettext.Tr("tomorrow"))
		}
	} else {
		msgStr = cFormatTime("%x %H:%M", t)
		if diff == 0 {
			msgStr = title(gettext.Tr("today")) + " " + t.Format(layoutHM)
		} else if diff == 1 {
			msgStr = title(gettext.Tr("tomorrow")) + " " + t.Format(layoutHM)
		}
	}
	return msgStr
}

func (job *JobJSON) getRemindBody(now time.Time) (body string) {
	today := libdate.NewAt(now)
	msgStart := getBodyTimePart(today, job.AllDay, job.Start, true)
	msgEnd := getBodyTimePart(today, job.AllDay, job.End, false)

	startDate := libdate.NewAt(job.Start)
	endDate := libdate.NewAt(job.End)
	to := gettext.Tr("%s to %s")
	var prefix string
	if job.AllDay {
		if startDate.Equal(endDate) {
			prefix = msgStart
		} else {
			prefix = fmt.Sprintf(to, msgStart, msgEnd)
		}
	} else {
		if startDate.Equal(endDate) {
			msgEnd = job.End.Format(layoutHM)
		}
		prefix = fmt.Sprintf(to, msgStart, msgEnd)
	}

	body = prefix + "  " + job.Title
	return body
}
