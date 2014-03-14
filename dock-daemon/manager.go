package main

import "dlib/dbus"
import pkgbus "dbus/org/freedesktop/dbus"
import "time"
import "sync"

var busdaemon *pkgbus.DBusDaemon

type Manager struct {
	Entries       []*EntryProxyer
	entrireLocker sync.Mutex

	Added   func(dbus.ObjectPath)
	Removed func(string)
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"dde.dock.EntryManager",
		"/dde/dock/EntryManager",
		"dde.dock.EntryManager",
	}
}

func NewManager() *Manager {
	m := &Manager{}
	return m
}

func (m *Manager) watchEntries() {
	var err error
	busdaemon, err = pkgbus.NewDBusDaemon("org.freedesktop.DBus", "/org/freedesktop/DBus")
	if err != nil {
		panic(err)
	}

	// register existed entries
	names, err := busdaemon.ListNames()
	if err != nil {
		panic(err)
	}
	for _, n := range names {
		m.registerEntry(n)
	}

	// monitor name lost, name acquire
	busdaemon.ConnectNameOwnerChanged(func(name, oldOwner, newOwner string) {
		// if a new dbus session was installed, the name and newOwner
		// will be not empty, if a dbus session was uninstalled, the
		// name and oldOwner will be not empty
		if len(newOwner) != 0 {
			go func() {
				// FIXME: how long time should to wait for
				time.Sleep(500 * time.Millisecond)
				m.entrireLocker.Lock()
				m.registerEntry(name)
				m.entrireLocker.Unlock()
			}()
		} else {
			m.entrireLocker.Lock()
			m.unregisterEntry(name)
			m.entrireLocker.Unlock()
		}
	})
}

func (m *Manager) registerEntry(name string) {
	if !isEntryNameValid(name) {
		return
	}
	logger.Debug("register entry: %s", name)
	entryId, ok := getEntryId(name)
	if !ok {
		return
	}
	logger.Debug("register entry id: %s", entryId)
	entry, err := NewEntryProxyer(entryId)
	if err != nil {
		logger.Error("register entry failed: %v", err)
		return
	}
	err = dbus.InstallOnSession(entry)
	if err != nil {
		logger.Error("register entry failed: %v", err)
		return
	}
	m.Entries = append(m.Entries, entry)
	logger.Info("register entry success: %s", name)
}

func (m *Manager) unregisterEntry(name string) {
	if !isEntryNameValid(name) {
		return
	}
	logger.Debug("unregister entry: %s", name)
	entryId, ok := getEntryId(name)
	if !ok {
		return
	}
	logger.Debug("unregister entry id: %s", entryId)

	// find the index
	index := -1
	var entry *EntryProxyer = nil
	for i, e := range m.Entries {
		if e.entryId == entryId {
			index = i
			entry = e
			break
		}
	}

	if index < 0 {
		logger.Warning("slice out of bounds, entry len: %d, index: %d", len(m.Entries), index)
		return
	}
	logger.Debug("entry len: %d, index: %d", len(m.Entries), index)

	if entry != nil {
		dbus.UnInstallObject(entry)
	}

	// remove the entry from slice
	copy(m.Entries[index:], m.Entries[index+1:])
	m.Entries[len(m.Entries)-1] = nil
	m.Entries = m.Entries[:len(m.Entries)-1]

	logger.Info("unregister entry success: %s", name)
}
