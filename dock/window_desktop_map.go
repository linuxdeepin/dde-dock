package dock

import (
	"encoding/gob"
	"os"
)

const windowDesktopMapVersion = "0.1"

type windowDesktopMap struct {
	content         *windowDesktopMapContent
	file            *os.File
	hasChanged      bool
	autoSaveEnabled bool
}

type windowDesktopMapContent struct {
	Version string
	Data    map[string][]string
}

func newWindowDesktopMapFromFile(filePath string) (*windowDesktopMap, error) {
	m := &windowDesktopMap{
		content:         &windowDesktopMapContent{},
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
	if m.content.Version != windowDesktopMapVersion {
		logger.Warning("version not match")
		m.cleanContent()
	} else {
		logger.Info("load file ok")
		logger.Debugf("load content: %#v", m.content)
	}
	return m, nil
}

func (m *windowDesktopMap) SetAutoSaveEnabled(val bool) {
	m.autoSaveEnabled = val
}

func (m *windowDesktopMap) Save() error {
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

func (m *windowDesktopMap) AutoSave() error {
	if !m.autoSaveEnabled {
		logger.Debug("Skip save")
		return nil
	}
	if m.hasChanged {
		return m.Save()
	}
	return nil
}

func (m *windowDesktopMap) Destroy() {
	// close file
	m.file.Close()
}

func (m *windowDesktopMap) cleanContent() {
	m.content.Version = windowDesktopMapVersion
	m.content.Data = make(map[string][]string)
	m.hasChanged = true
}

func (m *windowDesktopMap) NewRel(windowHash, desktopHash string) {
	logger.Debugf("NewRel win hash %v => desktop hash %v", windowHash, desktopHash)
	const windowHashSliceLenLimit = 50
	dHash := m.FindRel(windowHash)
	if dHash != "" {
		logger.Warningf("NewRel failed: win hash %v => desktop hash %v exist", windowHash, dHash)
		return
	}

	data := m.content.Data
	slice, ok := data[desktopHash]
	if ok {
		// 如果尺寸超过限制值，说明缓存不起作用，清除之前的缓存
		if len(slice) > windowHashSliceLenLimit {
			data[desktopHash] = []string{windowHash}
			m.hasChanged = true
			return
		}

		data[desktopHash] = append(slice, windowHash)
		m.hasChanged = true
	} else {
		data[desktopHash] = []string{windowHash}
		m.hasChanged = true
	}
}

func (m *windowDesktopMap) FindRel(windowHash string) string {
	for desktopHash, slice := range m.content.Data {
		if isStrInSlice(windowHash, slice) {
			return desktopHash
		}
	}
	return ""
}

func (m *windowDesktopMap) DelRel(windowHash, desktopHash string) {
	data := m.content.Data
	slice, ok := data[desktopHash]
	if ok {
		var index int = -1
		for i, v := range slice {
			if v == windowHash {
				index = i
			}
		}
		if index != -1 {
			data[desktopHash] = append(slice[:index], slice[index+1:]...)
			logger.Debugf("DelRel: window hash %v => desktop hash %v", windowHash, desktopHash)
			m.hasChanged = true
		}
	}
}
