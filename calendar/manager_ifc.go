package calendar

import (
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName = "com.deepin.daemon.Calendar"
	dbusPath        = "/com/deepin/daemon/Calendar/Scheduler"
	dbusInterface   = "com.deepin.daemon.Calendar.Scheduler"
)

func (m *Manager) GetJob(id int64) (string, *dbus.Error) {
	job, err := m.getJob(uint(id))
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	result, err := toJson(job)
	return result, dbusutil.ToError(err)
}

func (m *Manager) DeleteJob(id int64) *dbus.Error {
	err := m.deleteJob(uint(id))
	if err == nil {
		m.notifyJobsChange()
	}
	return dbusutil.ToError(err)
}

func (m *Manager) UpdateJob(jobStr string) *dbus.Error {
	var jj JobJSON
	err := fromJson(jobStr, &jj)
	if err != nil {
		return dbusutil.ToError(err)
	}

	job, err := jj.toJob()
	if err != nil {
		return dbusutil.ToError(err)
	}
	err = m.updateJob(job)
	if err == nil {
		m.notifyJobsChange()
	}
	return dbusutil.ToError(err)
}

func (m *Manager) CreateJob(jobStr string) (int64, *dbus.Error) {
	var jj JobJSON
	err := fromJson(jobStr, &jj)
	if err != nil {
		return 0, dbusutil.ToError(err)
	}

	job, err := jj.toJob()
	if err != nil {
		return 0, dbusutil.ToError(err)
	}
	err = m.createJob(job)
	if err != nil {
		return 0, dbusutil.ToError(err)
	}
	m.notifyJobsChange()
	return int64(job.ID), nil
}

func (m *Manager) GetTypes() (string, *dbus.Error) {
	types, err := m.getTypes()
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	result, err := toJson(types)
	return result, dbusutil.ToError(err)
}

func (m *Manager) GetType(id int64) (string, *dbus.Error) {
	t, err := m.getType(uint(id))
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	result, err := toJson(t)
	return result, dbusutil.ToError(err)
}

func (m *Manager) DeleteType(id int64) *dbus.Error {
	err := m.deleteType(uint(id))
	return dbusutil.ToError(err)
}

func (m *Manager) CreateType(typeStr string) (int64, *dbus.Error) {
	var jobType JobTypeJSON
	err := fromJson(typeStr, &jobType)
	if err != nil {
		return 0, dbusutil.ToError(err)
	}

	jt := jobType.toJobType()
	err = m.createType(jt)
	if err != nil {
		return 0, dbusutil.ToError(err)
	}
	return int64(jt.ID), nil
}

func (m *Manager) UpdateType(typeStr string) *dbus.Error {
	var jobType JobTypeJSON
	err := fromJson(typeStr, &jobType)
	if err != nil {
		return dbusutil.ToError(err)
	}

	jt := jobType.toJobType()
	err = m.updateType(jt)
	return dbusutil.ToError(err)
}
