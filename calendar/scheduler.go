package calendar

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"
	libdate "github.com/rickb777/date"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

type Scheduler struct {
	signalLoop     *dbusutil.SignalLoop
	db             *gorm.DB
	notifications  *notifications.Notifications
	notifyJobMap   map[uint32]*JobJSON // key is notification id
	notifyJobMapMu sync.Mutex
	timerGroup     timerGroup
	changeChan     chan struct{}
	quitChan       chan struct{}

	methods *struct {
		GetJobs   func() `in:"startYear,startMonth,startDay,endYear,endMonth,endDay" out:"jobs"`
		GetJob    func() `in:"id" out:"job"`
		DeleteJob func() `in:"id"`
		UpdateJob func() `in:"jobInfo"`
		CreateJob func() `in:"jobInfo" out:"id"`

		GetTypes   func() `out:"types"`
		GetType    func() `in:"id" out:"type"`
		DeleteType func() `in:"id"`
		UpdateType func() `in:"typeInfo"`
		CreateType func() `in:"typeInfo" out:"id"`

		DebugRemindJob func() `in:"id"`
	}
}

func newScheduler(db *gorm.DB, service *dbusutil.Service) *Scheduler {
	sessionBus := service.Conn()
	s := &Scheduler{
		db:            db,
		changeChan:    make(chan struct{}),
		quitChan:      make(chan struct{}),
		notifications: notifications.NewNotifications(sessionBus),
		notifyJobMap:  make(map[uint32]*JobJSON),
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
			s.notifyJobMapMu.Lock()
			delete(s.notifyJobMap, id)
			s.notifyJobMapMu.Unlock()
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
	result := getJobsBetween(startDate, endDate, allJobs, true)
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
	var types []JobType
	err := s.db.Find(&types).Error
	if err != nil {
		return nil, err
	}

	result := make([]*JobTypeJSON, len(types))
	for idx, t := range types {
		result[idx] = t.toJobTypeJSON()
	}
	return result, nil
}

func (s *Scheduler) getType(id uint) (*JobTypeJSON, error) {
	var jobType JobType
	err := s.db.First(&jobType, id).Error
	if err != nil {
		return nil, err
	}
	return jobType.toJobTypeJSON(), nil
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
		m.remindJob(job)
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

func (s *Scheduler) remindJob(job *Job) {
	layout := "2006-01-02 15:04"
	body := job.Start.Format(layout) + " ~ " + job.End.Format(layout)
	logger.Debugf("remind now: %v, title: %v, body: %v", time.Now(), job.Title, body)
	id, err := s.notifications.Notify(0, "dde-daemon", 0,
		"dde-calendar", job.Title,
		body, nil, nil, 0)
	if err != nil {
		logger.Warning(err)
		return
	}
	logger.Debug("id:", id)
	jj, err := job.toJobJSON()
	if err != nil {
		logger.Warning(err)
		return
	}
	s.notifyJobMapMu.Lock()
	s.notifyJobMap[id] = jj
	s.notifyJobMapMu.Unlock()
}

func (s *Scheduler) DebugRemindJob(id int64) *dbus.Error {
	var job Job
	err := s.db.First(&job, id).Error
	if err != nil {
		return dbusutil.ToError(err)
	}

	s.remindJob(&job)
	return nil
}

func (s *Scheduler) notifyJobsChange() {
	s.changeChan <- struct{}{}
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

			case <-s.changeChan:
				now := time.Now()
				setTimerGroup(now)

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

	wraps := getJobsBetween(startDate, endDate, allJobs, false)
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
