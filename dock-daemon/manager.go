package main

import "dlib/dbus"
import pkgbus "dbus/org/freedesktop/dbus"

var busdaemon *pkgbus.DBusDaemon

// var busdaemon, _ = pkgbus.NewDBusDaemon("org.freedesktop.DBus", "/")

type Manager struct {
	Entries []*EntryProxyer

	Added   func(dbus.ObjectPath)
	Removed func(string)
}

func NewManager() *Manager {
	m := &Manager{}
	return m
}

func (m *Manager) watchEntries() {
	var err error
	busdaemon, err = pkgbus.NewDBusDaemon("org.freedesktop.DBus", "/")
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
	busdaemon.ConnectNameAcquired(func(name string) {
		logger.Debug("dbus name acquired: ", name)
		m.registerEntry(name)
	})
	busdaemon.ConnectNameLost(func(name string) {
		logger.Debug("dbus name lost: ", name)
		m.unregisterEntry(name)
	})
}

func (m *Manager) registerEntry(name string) {
	if !isEntryNameValid(name) {
		return
	}
	logger.Debug("register entry: ", name)
	entryId, ok := getEntryId(name)
	if !ok {
		return
	}
	logger.Debug("register entry id: ", entryId)
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
	logger.Info("register entry: ", name)
}

func (m *Manager) unregisterEntry(name string) {
	if !isEntryNameValid(name) {
		return
	}
	logger.Debug("unregister entry: ", name)
	entryId, ok := getEntryId(name)
	if !ok {
		return
	}
	logger.Debug("unregister entry id: ", entryId)

	// find the index
	var index int
	var entry *EntryProxyer
	for i, e := range m.Entries {
		if e.entryId == entryId {
			index = i
			entry = e
		}
	}

	dbus.UnInstallObject(entry)

	// remove the entry from slice
	copy(m.Entries[index:], m.Entries[index+1:])
	m.Entries[len(m.Entries)-1] = nil
	m.Entries = m.Entries[:len(m.Entries)-1]

	logger.Info("unregister entry: ", name)
}
