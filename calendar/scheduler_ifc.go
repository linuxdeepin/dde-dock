package calendar

import (
	"time"

	dbus "github.com/godbus/dbus"
	libdate "github.com/rickb777/date"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName = "com.deepin.daemon.Calendar"
	dbusPath        = "/com/deepin/daemon/Calendar/Scheduler"
	dbusInterface   = "com.deepin.daemon.Calendar.Scheduler"
)

type queryJobsParams struct {
	Key   string
	Start time.Time
	End   time.Time
}

func (s *Scheduler) QueryJobs(paramsStr string) (string, *dbus.Error) {
	var params queryJobsParams
	err := fromJson(paramsStr, &params)
	if err != nil {
		return "", dbusutil.ToError(err)
	}

	jobs, err := s.queryJobs(params.Key, params.Start, params.End)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	result, err := toJson(jobs)
	return result, dbusutil.ToError(err)
}

func (s *Scheduler) QueryJobsWithLimit(paramsStr string, maxNum int32) (string, *dbus.Error) {
	var params queryJobsParams
	err := fromJson(paramsStr, &params)
	if err != nil {
		return "", dbusutil.ToError(err)
	}

	jobs, err := s.queryJobsWithLimit(params.Key, params.Start, params.End, maxNum)
	if err != nil {
		return "", dbusutil.ToError(err)
	}

	result, err := toJson(jobs)
	return result, dbusutil.ToError(err)
}

func (s *Scheduler) QueryJobsWithRule(paramsStr string, rule string) (string, *dbus.Error) {
	var params queryJobsParams
	err := fromJson(paramsStr, &params)
	if err != nil {
		return "", dbusutil.ToError(err)
	}

	jobs, err := s.queryJobsWithRule(params.Key, params.Start, params.End, rule)
	if err != nil {
		return "", dbusutil.ToError(err)
	}

	result, err := toJson(jobs)
	return result, dbusutil.ToError(err)
}

func (s *Scheduler) GetJobs(startYear, startMonth, startDay, endYear, endMonth, endDay int32) (string, *dbus.Error) {
	startDate := libdate.New(int(startYear), time.Month(startMonth), int(startDay))
	endDate := libdate.New(int(endYear), time.Month(endMonth), int(endDay))
	jobs, err := s.getJobs(startDate, endDate)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	result, err := toJson(jobs)
	return result, dbusutil.ToError(err)
}

func (s *Scheduler) GetJobsWithLimit(startYear, startMonth, startDay, endYear, endMonth, endDay int32, maxNum int32) (string, *dbus.Error) {
	startDate := libdate.New(int(startYear), time.Month(startMonth), int(startDay))
	endDate := libdate.New(int(endYear), time.Month(endMonth), int(endDay))
	jobs, err := s.getJobsWithLimit(startDate, endDate, maxNum)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	result, err := toJson(jobs)
	return result, dbusutil.ToError(err)
}

func (s *Scheduler) getJobsWithLimit(startDate, endDate libdate.Date, maxNum int32) ([]dateJobsWrap, error) {
	var allJobs []*Job
	err := s.db.Find(&allJobs).Error
	if err != nil {
		return nil, err
	}
	t0 := time.Now()
	wraps := getJobsBetween(startDate, endDate, allJobs, false, "", false)
	result := takeFirstNJobs(wraps, maxNum)
	logger.Debug("cost time:", time.Since(t0))
	return result, nil
}

func (s *Scheduler) GetJobsWithRule(startYear, startMonth, startDay, endYear, endMonth, endDay int32, rule string) (string, *dbus.Error) {
	startDate := libdate.New(int(startYear), time.Month(startMonth), int(startDay))
	endDate := libdate.New(int(endYear), time.Month(endMonth), int(endDay))
	jobs, err := s.getJobsWithRule(startDate, endDate, rule)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	result, err := toJson(jobs)
	return result, dbusutil.ToError(err)
}

func (s *Scheduler) getJobsWithRule(startDate, endDate libdate.Date, rule string) ([]dateJobsWrap, error) {
	allJobs, err := getJobsWithRule(s.db, rule)
	if err != nil {
		return nil, err
	}
	var result []dateJobsWrap
	t0 := time.Now()
	wraps := getJobsBetween(startDate, endDate, allJobs, true, "", false)
	for _, item := range wraps {
		if item.jobs != nil {
			result = append(result, item)
		}
	}
	logger.Debug("cost time:", time.Since(t0))
	return result, nil
}

func (s *Scheduler) GetJob(id int64) (string, *dbus.Error) {
	job, err := s.getJob(uint(id))
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	result, err := toJson(job)
	return result, dbusutil.ToError(err)
}

func (s *Scheduler) DeleteJob(id int64) *dbus.Error {
	err := s.deleteJob(uint(id))
	if err == nil {
		s.notifyJobsChange(uint(id))
	}
	return dbusutil.ToError(err)
}

func (s *Scheduler) UpdateJob(jobStr string) *dbus.Error {
	var jj JobJSON
	err := fromJson(jobStr, &jj)
	if err != nil {
		return dbusutil.ToError(err)
	}

	job, err := jj.toJob()
	if err != nil {
		return dbusutil.ToError(err)
	}
	err = s.updateJob(job)
	if err == nil {
		s.notifyJobsChange(job.ID)
	}
	return dbusutil.ToError(err)
}

func (s *Scheduler) CreateJob(jobStr string) (int64, *dbus.Error) {
	var jj JobJSON
	err := fromJson(jobStr, &jj)
	if err != nil {
		return 0, dbusutil.ToError(err)
	}

	job, err := jj.toJob()
	if err != nil {
		return 0, dbusutil.ToError(err)
	}
	err = s.createJob(job)
	if err != nil {
		return 0, dbusutil.ToError(err)
	}
	s.notifyJobsChange(job.ID)
	return int64(job.ID), nil
}

func (s *Scheduler) GetTypes() (string, *dbus.Error) {
	types, err := s.getTypes()
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	result, err := toJson(types)
	return result, dbusutil.ToError(err)
}

func (s *Scheduler) GetType(id int64) (string, *dbus.Error) {
	t, err := s.getType(uint(id))
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	result, err := toJson(t)
	return result, dbusutil.ToError(err)
}

func (s *Scheduler) DeleteType(id int64) *dbus.Error {
	err := s.deleteType(uint(id))
	return dbusutil.ToError(err)
}

func (s *Scheduler) CreateType(typeStr string) (int64, *dbus.Error) {
	var jobType JobTypeJSON
	err := fromJson(typeStr, &jobType)
	if err != nil {
		return 0, dbusutil.ToError(err)
	}

	jt := jobType.toJobType()
	err = s.createType(jt)
	if err != nil {
		return 0, dbusutil.ToError(err)
	}
	return int64(jt.ID), nil
}

func (s *Scheduler) UpdateType(typeStr string) *dbus.Error {
	var jobType JobTypeJSON
	err := fromJson(typeStr, &jobType)
	if err != nil {
		return dbusutil.ToError(err)
	}

	jt := jobType.toJobType()
	err = s.updateType(jt)
	return dbusutil.ToError(err)
}
