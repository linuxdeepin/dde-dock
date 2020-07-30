package calendar

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"
	libdate "github.com/rickb777/date"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
)

type Scheduler struct {
	signalLoop          *dbusutil.SignalLoop
	db                  *gorm.DB
	service             *dbusutil.Service
	notifications       *notifications.Notifications
	notifyJobMap        map[uint32]*JobJSON // key is notification id
	notifyJobMapMu      sync.Mutex
	timerGroup          timerGroup
	remindLaterTimers   map[uint]*time.Timer // key is job id
	remindLaterTimersMu sync.Mutex
	festivalJobEnabled  bool

	changeChan chan []uint
	quitChan   chan struct{}

	methods *struct {
		GetJobs          func() `in:"startYear,startMonth,startDay,endYear,endMonth,endDay" out:"jobs"`
		GetJob           func() `in:"id" out:"job"`
		GetJobsWithLimit func() `in:"startYear,startMonth,startDay,endYear,endMonth,endDay,maxNum" out:"jobs"`
		GetJobsWithRule  func() `in:"startYear,startMonth,startDay,endYear,endMonth,endDay,rule" out:"jobs"`
		QueryJobs        func() `in:"params" out:"jobs"`
		DeleteJob        func() `in:"id"`
		UpdateJob        func() `in:"jobInfo"`
		CreateJob        func() `in:"jobInfo" out:"id"`

		GetTypes   func() `out:"types"`
		GetType    func() `in:"id" out:"type"`
		DeleteType func() `in:"id"`
		UpdateType func() `in:"typeInfo"`
		CreateType func() `in:"typeInfo" out:"id"`

		DebugRemindJob func() `in:"id"`
	}

	signals *struct {
		JobsUpdated struct {
			Ids []int64
		}
	}
}

func newScheduler(db *gorm.DB, service *dbusutil.Service) *Scheduler {
	sessionBus := service.Conn()
	s := &Scheduler{
		db:                db,
		service:           service,
		changeChan:        make(chan []uint),
		quitChan:          make(chan struct{}),
		notifications:     notifications.NewNotifications(sessionBus),
		notifyJobMap:      make(map[uint32]*JobJSON),
		remindLaterTimers: make(map[uint]*time.Timer),
	}
	if isZH() {
		s.festivalJobEnabled = true
	}
	s.signalLoop = dbusutil.NewSignalLoop(sessionBus, 10)
	s.signalLoop.Start()
	s.listenDBusSignals()
	return s
}

func (s *Scheduler) destroy() {
	s.notifications.RemoveAllHandlers()
	s.signalLoop.Stop()
	close(s.quitChan)
}

func (s *Scheduler) GetInterfaceName() string {
	return dbusInterface
}

const (
	notifyCloseReasonDismissedByUser = 2
)

func (s *Scheduler) listenDBusSignals() {
	s.notifications.InitSignalExt(s.signalLoop, true)
	_, err := s.notifications.ConnectNotificationClosed(func(id uint32, reason uint32) {
		logger.Debug("signal notification closed", id, reason)
		defer func() {
			time.AfterFunc(100*time.Millisecond, func() {
				s.notifyJobMapMu.Lock()
				delete(s.notifyJobMap, id)
				s.notifyJobMapMu.Unlock()
			})
		}()

		if reason != notifyCloseReasonDismissedByUser {
			return
		}
		s.notifyJobMapMu.Lock()
		job := s.notifyJobMap[id]
		s.notifyJobMapMu.Unlock()

		if job == nil {
			return
		}

		err := callUIOpenSchedule(job)
		if err != nil {
			logger.Warning("failed to show job:", err)
		}
	})
	if err != nil {
		logger.Warning(err)
	}

	_, err = s.notifications.ConnectActionInvoked(func(notifyId uint32, actionKey string) {
		logger.Debug("signal action invoked", notifyId, actionKey)
		s.notifyJobMapMu.Lock()
		job := s.notifyJobMap[notifyId]
		s.notifyJobMapMu.Unlock()

		if job == nil {
			return
		}

		switch actionKey {
		case notifyActKeyRemindLater:

			logger.Debug("remind later", job.ID)
			job.remindLaterCount++
			s.remindJobLater(job)

		case notifyActKeyRemind1DayBefore:
			err = s.setJobRemindOneDayBefore(job)
			if err != nil {
				logger.Warning(err)
			}

		case notifyActKeyRemindTomorrow:
			err = s.setJobRemindTomorrow(job)
			if err != nil {
				logger.Warning(err)
			}
		}
	})
}

func (s *Scheduler) emitJobsUpdated(ids ...uint) {
	ids0 := make([]int64, len(ids))
	for idx, value := range ids {
		ids0[idx] = int64(value)
	}
	err := s.service.Emit(s, "JobsUpdated", ids0)
	if err != nil {
		logger.Warning(err)
	}
}

const (
	uiDBusPath      = "/com/deepin/Calendar"
	uiDBusInterface = "com.deepin.Calendar"
	uiDBusService   = uiDBusInterface
)

func callUIOpenSchedule(job *JobJSON) error {
	bus, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	jobStr, err := toJson(job)
	if err != nil {
		return err
	}
	obj := bus.Object(uiDBusService, uiDBusPath)
	err = obj.Call(uiDBusInterface+".OpenSchedule", 0, jobStr).Err
	return err
}

type dateJobsWrap struct {
	date       libdate.Date
	jobs       []*Job
	extendJobs []*Job
}

type dateJobsWrapJSON struct {
	Date string
	Jobs []*JobJSON
}

func (w *dateJobsWrap) MarshalJSON() ([]byte, error) {
	var wj dateJobsWrapJSON
	wj.Date = w.date.String()
	wj.Jobs = make([]*JobJSON, len(w.jobs)+len(w.extendJobs))
	var err error
	for idx, j := range w.jobs {
		wj.Jobs[idx], err = j.toJobJSON()
		if err != nil {
			return nil, err
		}
	}
	baseIdx := len(w.jobs)
	for idx, j := range w.extendJobs {
		wj.Jobs[idx+baseIdx], err = j.toJobJSON()
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(wj)
}

func (s *Scheduler) getJobs(startDate, endDate libdate.Date) ([]dateJobsWrap, error) {
	var allJobs []*Job
	err := s.db.Find(&allJobs).Error
	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	result := getJobsBetween(startDate, endDate, allJobs, true, "", s.festivalJobEnabled)
	logger.Debug("cost time:", time.Since(t0))
	return result, nil
}

func (s *Scheduler) getJob(id uint) (*JobJSON, error) {
	var job Job
	err := s.db.First(&job, id).Error
	if err != nil {
		return nil, err
	}
	return job.toJobJSON()
}

func (s *Scheduler) queryJobs(key string, startTime, endTime time.Time) ([]dateJobsWrap, error) {
	var allJobs []*Job
	db := s.db

	key = strings.TrimSpace(key)
	if canQueryByPinyin(key) {
		var pinyin = createPinyinQuery(strings.ToLower(key))
		db = db.Where("instr(UPPER(title), UPPER(?)) OR title_pinyin LIKE ?", key, pinyin)
	} else if key != "" {
		db = db.Where("instr(UPPER(title), UPPER(?))", key)
	}

	err := db.Find(&allJobs).Error
	if err != nil {
		return nil, err
	}

	startDate := libdate.NewAt(startTime)
	endDate := libdate.NewAt(endTime)

	result := getJobsBetween(startDate, endDate, allJobs, true, key, s.festivalJobEnabled)

	timeRange := TimeRange{
		start: startTime,
		end:   endTime,
	}
	result = filterDateJobsWrap(result, timeRange)
	return result, nil
}

func filterDateJobsWrap(wraps []dateJobsWrap, timeRange TimeRange) []dateJobsWrap {
	var result []dateJobsWrap
	for _, wrap := range wraps {
		wrap.jobs = filterJobs(wrap.jobs, timeRange)
		wrap.extendJobs = filterJobs(wrap.extendJobs, timeRange)
		if len(wrap.jobs)+len(wrap.extendJobs) > 0 {
			result = append(result, wrap)
		}
	}
	return result
}

func filterJobs(jobs []*Job, timeRange TimeRange) []*Job {
	var result []*Job
	for _, job := range jobs {
		if job.timeRange().overlap(timeRange) {
			result = append(result, job)
		}
	}
	return result
}

func (s *Scheduler) deleteJob(id uint) error {
	var job Job
	err := s.db.Select("id").First(&job, id).Error
	if err != nil {
		return err
	}

	return s.db.Unscoped().Delete(&job).Error
}

func (s *Scheduler) updateJob(job *Job) error {
	err := job.validate()
	if err != nil {
		return err
	}
	var job0 Job
	err = s.db.Find(&job0, job.ID).Error
	if err != nil {
		return err
	}

	diffMap := make(map[string]interface{})

	if job0.Type != job.Type {
		diffMap["Type"] = job.Type
	}
	if job0.Title != job.Title {
		diffMap["Title"] = job.Title
	}
	if job0.TitlePinyin != job.TitlePinyin {
		diffMap["TitlePinyin"] = job.TitlePinyin
	}
	if job0.Description != job.Description {
		diffMap["Description"] = job.Description
	}
	if job0.AllDay != job.AllDay {
		diffMap["AllDay"] = job.AllDay
	}
	if !job0.Start.Equal(job.Start) {
		diffMap["Start"] = job.Start
	}
	if !job0.End.Equal(job.End) {
		diffMap["End"] = job.End
	}
	if job0.RRule != job.RRule {
		diffMap["RRule"] = job.RRule
	}
	if job0.Remind != job.Remind {
		diffMap["Remind"] = job.Remind
	}
	if job0.Ignore != job.Ignore {
		diffMap["Ignore"] = job.Ignore
	}

	if len(diffMap) > 0 {
		err = s.db.Model(job).Updates(diffMap).Error
	}
	return err
}

func (s *Scheduler) createJob(job *Job) error {
	err := job.validate()
	if err != nil {
		return err
	}
	job.ID = 0

	err = s.db.Create(job).Error
	return err
}

func (s *Scheduler) getTypes() ([]*JobTypeJSON, error) {
	return globalPredefinedTypes, nil
	//var types []JobType
	//err := s.db.Find(&types).Error
	//if err != nil {
	//	return nil, err
	//}
	//
	//result := make([]*JobTypeJSON, len(types))
	//for idx, t := range types {
	//	result[idx] = t.toJobTypeJSON()
	//}
	//return result, nil
}

func (s *Scheduler) getType(id uint) (*JobTypeJSON, error) {
	for _, jobType := range globalPredefinedTypes {
		if jobType.ID == id {
			return jobType, nil
		}
	}
	return nil, errors.New("job type not found")
	//var jobType JobType
	//err := s.db.First(&jobType, id).Error
	//if err != nil {
	//	return nil, err
	//}
	//return jobType.toJobTypeJSON(), nil
}

func (s *Scheduler) deleteType(id uint) error {
	var jobType JobType
	err := s.db.Select("id").First(&jobType, id).Error
	if err != nil {
		return err
	}

	return s.db.Delete(&jobType).Error
}

func (s *Scheduler) createType(jobType *JobType) error {
	err := jobType.validate()
	if err != nil {
		return err
	}
	jobType.ID = 0
	return s.db.Create(jobType).Error
}

func (s *Scheduler) updateType(jobType *JobType) error {
	err := jobType.validate()
	if err != nil {
		return err
	}
	var jobType0 JobType
	err = s.db.Find(&jobType0, jobType.ID).Error
	if err != nil {
		return err
	}

	diffMap := make(map[string]interface{})

	if jobType0.Name != jobType.Name {
		diffMap["Name"] = jobType.Name
	}
	if jobType0.Color != jobType.Color {
		diffMap["Color"] = jobType.Color
	}

	if len(diffMap) > 0 {
		err = s.db.Model(jobType).Updates(diffMap).Error
	}
	return err
}

type timerGroup struct {
	timers []*time.Timer
}

func (tg *timerGroup) addJob(job *Job, m *Scheduler) {
	now := time.Now()
	duration := job.remindTime.Sub(now)
	logger.Debugf("notify job %d %q after %v", job.ID, job.Title, duration)
	tg.timers = append(tg.timers, time.AfterFunc(duration, func() {
		jj, err := job.toJobJSON()
		if err != nil {
			logger.Warning(err)
			return
		}
		m.remindJob(jj)
	}))
}

func (tg *timerGroup) reset() {
	if tg.timers == nil {
		return
	}
	for _, timer := range tg.timers {
		timer.Stop()
	}
	tg.timers = nil
}

const (
	notifyActKeyClose            = "close"
	notifyActKeyRemindLater      = "later"
	notifyActKeyRemind1DayBefore = "one-day-before"
	notifyActKeyRemindTomorrow   = "tomorrow"

	layoutHM = "15:04"
)

func (s *Scheduler) remindJob(job *JobJSON) {
	if job.Remind == "" {
		logger.Warning("job.Remind is empty")
		return
	}

	now := time.Now()

	nDays, err := getRemindAdvanceDays(job.Remind)
	if err != nil {
		logger.Warning(err)
		return
	}

	var actions []string
	duration, durationMax := getRemindLaterDuration(job.remindLaterCount + 1)
	if nDays >= 3 && job.remindLaterCount == 1 {
		actions = []string{
			notifyActKeyRemind1DayBefore, gettext.Tr("One day before start"),
			notifyActKeyClose, gettext.Tr("Close"),
		}
	} else if (nDays == 1 || nDays == 2) && durationMax {
		actions = []string{
			notifyActKeyRemindTomorrow, gettext.Tr("Remind me tomorrow"),
			notifyActKeyClose, gettext.Tr("Close"),
		}
	} else {
		nextRemindTime := now.Add(duration)
		logger.Debug("nextRemindTime:", nextRemindTime)
		if nextRemindTime.Before(job.Start) {
			actions = []string{
				notifyActKeyRemindLater, gettext.Tr("Remind me later"),
				notifyActKeyClose, gettext.Tr("Close"),
			}
		} else {
			actions = []string{
				notifyActKeyClose, gettext.Tr("Close"),
			}
		}
	}

	title := gettext.Tr("Schedule Reminder")
	body := job.getRemindBody(now)
	logger.Debugf("remind now: %v, title: %v, body: %v, actions: %#v",
		now, title, body, actions)
	id, err := s.notifications.Notify(0, "dde-calendar", 0,
		"dde-calendar", title,
		body, actions, nil, 0)
	if err != nil {
		logger.Warning(err)
		return
	}
	logger.Debug("notify id:", id)

	s.notifyJobMapMu.Lock()
	s.notifyJobMap[id] = job
	s.notifyJobMapMu.Unlock()
}

func (s *Scheduler) setJobRemind(jj *JobJSON, remind string) error {
	var job Job
	err := s.db.Find(&job, jj.ID).Error
	if err != nil {
		return err
	}

	if job.RRule != "" {
		newJob := Job{
			Type:        jj.Type,
			Title:       jj.Title,
			Description: jj.Description,
			AllDay:      jj.AllDay,
			Start:       jj.Start,
			End:         jj.End,
			Remind:      remind,
		}
		err = s.withTx(func(tx *gorm.DB) error {
			ignore, err := job.getIgnore()
			if err != nil {
				return err
			}

			if !timeSliceContains(ignore, jj.Start) {
				ignore = append(ignore, jj.Start)
				err = job.setIgnore(ignore)
				if err != nil {
					return err
				}

				err = tx.Model(&job).Update("Ignore", job.Ignore).Error
				if err != nil {
					return err
				}

			} else {
				logger.Warning("job.ignore already contains jj.Start")
			}

			err = tx.Create(&newJob).Error
			return err
		})
		if err != nil {
			return err
		}
		s.notifyJobsChange(job.ID, newJob.ID)
		s.emitJobsUpdated(job.ID, newJob.ID)

	} else {
		err = s.db.Model(&job).Update("Remind", remind).Error
		if err != nil {
			return err
		}
		s.notifyJobsChange(job.ID)
		s.emitJobsUpdated(job.ID)
	}

	return nil
}

func (s *Scheduler) withTx(fn func(db *gorm.DB) error) (err error) {
	tx := s.db.Begin()
	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and re-panic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit().Error
		}
	}()

	err = fn(tx)
	return err
}

func (s *Scheduler) setJobRemindOneDayBefore(jj *JobJSON) error {
	// 非全天，提醒改成24小时前
	// 全天，提醒改成一天前的09:00。
	remind := "1440"
	if jj.AllDay {
		remind = "1;09:00"
	}
	return s.setJobRemind(jj, remind)
}

func (s *Scheduler) setJobRemindTomorrow(jj *JobJSON) error {
	// 非全天，提醒改成1小时前；
	// 全天，提醒改成当天09:00。
	remind := "60"
	if jj.AllDay {
		remind = "0;09:00"
	}
	return s.setJobRemind(jj, remind)
}

func getRemindLaterDuration(count int) (time.Duration, bool) {
	max := false
	duration := time.Duration(10+((count-1)*5)) * time.Minute
	if duration >= time.Hour {
		max = true
		duration = time.Hour
	}

	if logger.GetLogLevel() == log.LevelDebug {
		duration = duration / 60
		if count >= 3 {
			max = true
		}
	}

	return duration, max
}

func (s *Scheduler) remindJobLater(job *JobJSON) {
	duration, _ := getRemindLaterDuration(job.remindLaterCount)
	logger.Debug("remindJobLater duration:", duration)
	timer := time.AfterFunc(duration, func() {
		s.remindLaterTimersMu.Lock()
		delete(s.remindLaterTimers, job.ID)
		s.remindLaterTimersMu.Unlock()

		s.remindJob(job)
	})
	s.remindLaterTimersMu.Lock()
	s.remindLaterTimers[job.ID] = timer
	s.remindLaterTimersMu.Unlock()
}

func (s *Scheduler) DebugRemindJob(id int64) *dbus.Error {
	var job Job
	err := s.db.First(&job, id).Error
	if err != nil {
		return dbusutil.ToError(err)
	}

	jj, err := job.toJobJSON()
	if err != nil {
		return dbusutil.ToError(err)
	}
	s.remindJob(jj)
	return nil
}

func (s *Scheduler) notifyJobsChange(ids ...uint) {
	s.changeChan <- ids
}

func (s *Scheduler) startRemindLoop() {
	const interval = 10 * time.Minute
	ticker := time.NewTicker(interval)

	setTimerGroup := func(now time.Time) {
		s.timerGroup.reset()
		tr := TimeRange{
			start: now,
			end:   now.Add(interval),
		}
		jobs, err := s.getRemindJobs(tr)
		if err != nil {
			logger.Warning(err)
			return
		}
		for _, job := range jobs {
			s.timerGroup.addJob(job, s)
		}
	}

	setTimerGroup(time.Now())
	go func() {
		for {
			select {
			case now := <-ticker.C:
				setTimerGroup(now)

			case ids := <-s.changeChan:
				now := time.Now()
				setTimerGroup(now)

				s.remindLaterTimersMu.Lock()
				for _, id := range ids {
					timer := s.remindLaterTimers[id]
					if timer != nil {
						logger.Debug("cancel remind later", id)
						timer.Stop()
						delete(s.remindLaterTimers, id)
					}
				}
				s.remindLaterTimersMu.Unlock()

			case <-s.quitChan:
				return
			}
		}
	}()
}

func (s *Scheduler) getRemindJobs(tr TimeRange) ([]*Job, error) {
	var allJobs []*Job
	err := s.db.Find(&allJobs, "remind IS NOT NULL AND remind != '' ").Error
	if err != nil {
		return nil, err
	}
	startDate := libdate.NewAt(tr.start)
	endDate := startDate.Add(8)

	var result []*Job

	wraps := getJobsBetween(startDate, endDate, allJobs, false, "", false)
	for _, wrap := range wraps {
		for _, job := range wrap.jobs {
			remindT, err := job.getRemindTime()
			if err != nil {
				continue
			}
			if !tr.start.After(remindT) &&
				!remindT.After(tr.end) {
				// tr.start <= remindT <= tr.end
				job.remindTime = remindT
				result = append(result, job)
			}
		}
	}

	return result, nil
}
