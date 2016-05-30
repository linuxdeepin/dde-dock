package dock

import (
	"encoding/gob"
	"os"
)

const desktopIdFileMapVersion = "0.1"

type desktopIdFileMap struct {
	content         *desktopIdFileMapContent
	file            *os.File
	hasChanged      bool
	autoSaveEnabled bool
}

type desktopIdFileMapContent struct {
	Version string
	Data    map[string][]string
}

func newDesktopIdFileMapFromFile(filePath string) (*desktopIdFileMap, error) {
	m := &desktopIdFileMap{
		content:         &desktopIdFileMapContent{},
		autoSaveEnabled: true,
	}

	var err error
	m.file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(m.file)
	err = dec.Decode(m.content)
	if err != nil {
		logger.Warning("Decode error:", err)
		m.cleanContent()
		return m, nil
	}
	// check version match
	if m.content.Version != desktopIdFileMapVersion {
		logger.Warning("version not match")
		m.cleanContent()
	} else {
		logger.Info("load file ok")
		logger.Debugf("load content: %#v", m.content)
	}
	return m, nil
}

func (m *desktopIdFileMap) cleanContent() {
	m.content.Version = windowDesktopMapVersion
	m.content.Data = make(map[string][]string)
	m.hasChanged = true
}

func (m *desktopIdFileMap) SetAutoSaveEnabled(val bool) {
	m.autoSaveEnabled = val
}

func (m *desktopIdFileMap) Save() error {
	logger.Info("call save")
	enc := gob.NewEncoder(m.file)
	// TODO: seek head
	_, err := m.file.Seek(0, 0)
	if err != nil {
		logger.Warning("save error", err)
		return err
	}
	logger.Debugf("save content: %#v", m.content)
	err = enc.Encode(m.content)
	if err != nil {
		logger.Warning("save error:", err)
		return err
	}
	err = m.file.Sync()
	if err != nil {
		logger.Warning("save error:", err)
		return err
	}
	m.hasChanged = false
	return nil
}

func (m *desktopIdFileMap) AutoSave() error {
	if !m.autoSaveEnabled {
		logger.Debug("Skip save")
		return nil
	}
	if m.hasChanged {
		return m.Save()
	}
	return nil
}

func (m *desktopIdFileMap) NewRel(desktopFile, desktopId string) {
	data := m.content.Data
	slice, ok := data[desktopId]
	if ok {
		if !isStrInSlice(desktopFile, slice) {
			data[desktopId] = append(slice, desktopFile)
			m.hasChanged = true
		}
	} else {
		data[desktopId] = []string{desktopFile}
		m.hasChanged = true
	}
}

// TODO 清理无效的关系
func (m *desktopIdFileMap) FindRelAppInfo(desktopId string) *AppInfo {
	slice, ok := m.content.Data[desktopId]
	if ok {
		for _, file := range slice {
			appInfo := NewAppInfoFromFile(file)
			if appInfo != nil {
				return appInfo
			}
		}
	}
	return nil
}
