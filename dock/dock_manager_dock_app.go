package dock

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
	"text/template"
)

const dockedItemTemplate string = `[Desktop Entry]
Name={{ .Name }}
Exec={{ .Exec }}
Icon={{ .Icon }}
Type=Application
Terminal=false
StartupNotify=false
`

type dockedItemInfo struct {
	Name, Icon, Exec string
}

func createScratchDesktopFile(id, title, icon, cmd string) (string, error) {
	logger.Debugf("create scratch file for %q", id)
	file := filepath.Join(scratchDir, addDesktopExt(id))
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		logger.Warning("Open file for write failed:", err)
		return "", err
	}

	defer f.Close()
	temp := template.Must(template.New("docked_item_temp").Parse(dockedItemTemplate))
	dockedItem := dockedItemInfo{title, icon, cmd}
	logger.Debugf("dockedItem: %#v", dockedItem)
	err = temp.Execute(f, dockedItem)
	if err != nil {
		return "", err
	}
	return file, nil
}

func removeScratchFiles(desktopFile string) {
	fileNoExt := trimDesktopExt(desktopFile)
	logger.Debug("removeScratchFiles", fileNoExt)
	extList := []string{".desktop", ".sh", ".png"}
	for _, ext := range extList {
		file := fileNoExt + ext
		if dutils.IsFileExist(file) {
			logger.Debugf("remove scratch file %q", file)
			err := os.Remove(file)
			if err != nil {
				logger.Warning("remove scratch file %q failed:", file, err)
			}
		}
	}
}

func createScratchDesktopFileWithAppEntry(entry *AppEntry) (string, error) {
	if entry.appInfo != nil {
		desktopFile := entry.appInfo.GetFileName()
		newDesktopFile := filepath.Join(scratchDir, entry.appInfo.innerId+".desktop")
		err := copyFileContents(desktopFile, newDesktopFile)
		if err != nil {
			return "", err
		}
		return newDesktopFile, nil
	}

	if entry.current == nil {
		return "", errors.New("entry.current is nil")
	}
	if err := os.MkdirAll(scratchDir, 0755); err != nil {
		return "", err
	}
	appId := entry.current.innerId
	title := entry.current.getDisplayName()
	// icon
	icon := entry.current.getIcon()
	if strings.HasPrefix(icon, "data:image") {
		path, err := dataUriToFile(icon, filepath.Join(scratchDir, appId+".png"))
		if err != nil {
			logger.Warning(err)
			icon = ""
		} else {
			icon = path
		}
	}
	if icon == "" {
		icon = "application-default-icon"
	}

	// cmd
	scriptContent := entry.getExec(false)
	scriptFile := filepath.Join(scratchDir, appId+".sh")
	err := ioutil.WriteFile(scriptFile, []byte(scriptContent), 0744)
	if err != nil {
		return "", err
	}
	cmd := scriptFile + " %U"

	file, err := createScratchDesktopFile(appId, title, icon, cmd)
	if err != nil {
		return "", err
	}
	return file, nil
}

func (m *DockManager) getDockedAppEntryByDesktopFilePath(desktopFilePath string) (*AppEntry, error) {
	return m.Entries.FilterDocked().GetByDesktopFilePath(desktopFilePath)
}

func (m *DockManager) saveDockedApps() {
	var list []string
	for _, entry := range m.Entries.FilterDocked() {
		path := entry.appInfo.GetFileName()
		list = append(list, zipDesktopPath(path))
	}
	m.DockedApps.Set(list)
}

func needScratchDesktop(appInfo *AppInfo) bool {
	if appInfo == nil {
		logger.Debug("needScratchDesktop: yes, appInfo is nil")
		return true
	}
	if appInfo.IsInstalled() {
		logger.Debug("needScratchDesktop: no, desktop is installed")
		return false
	}
	file := appInfo.GetFileName()
	if isFileInDir(file, scratchDir) {
		logger.Debug("needScratchDesktop: no, desktop in scratchDir")
		return false
	}
	logger.Debug("needScratchDesktop: yes")
	return true
}

func (m *DockManager) dockEntry(entry *AppEntry) bool {
	entry.dockMutex.Lock()
	defer entry.dockMutex.Unlock()

	if entry.IsDocked {
		logger.Warningf("dockEntry failed: entry %v is docked", entry.Id)
		return false
	}
	if needScratchDesktop(entry.appInfo) {
		file, err := createScratchDesktopFileWithAppEntry(entry)
		if err == nil {
			logger.Debug("dockEntry: createScratchDesktopFile successfully", file)
			appInfo := NewAppInfoFromFile(file)
			entry.setAppInfo(appInfo)
			entry.updateIcon()
			entry.innerId = entry.appInfo.innerId
		} else {
			logger.Warning("createScratchDesktopFileWithAppEntry failed", err)
			return false
		}
	}

	entry.setIsDocked(true)
	entry.updateMenu()
	m.saveDockedApps()
	return true
}

func isFileInDir(file, dir string) bool {
	fileDir := filepath.Dir(file)
	return fileDir == dir
}

func (m *DockManager) undockEntry(entry *AppEntry) {
	entry.dockMutex.Lock()
	defer entry.dockMutex.Unlock()

	if !entry.IsDocked {
		logger.Warningf("undockEntry failed: entry %v is not docked", entry.Id)
		return
	}

	if entry.appInfo == nil {
		logger.Warning("undockEntry failed: entry.appInfo is nil")
		return
	}
	desktop := entry.appInfo.GetFileName()
	logger.Debugf("undockEntry desktop: %q", desktop)
	isDesktopInScratchDir := false
	if isFileInDir(desktop, scratchDir) {
		isDesktopInScratchDir = true
		removeScratchFiles(entry.appInfo.GetFileName())
	}

	if !entry.hasWindow() {
		m.removeAppEntry(entry)
	} else {
		if isDesktopInScratchDir && entry.current != nil {
			if strings.HasPrefix(filepath.Base(desktop), windowHashPrefix) {
				// desktop base starts with w:
				// 由于有 Pid 识别方法在，在这里不能用 m.identifyWindow 再次识别
				entry.innerId = entry.current.innerId
				entry.setAppInfo(nil)
			} else {
				// desktop base starts with d:
				var newAppInfo *AppInfo
				logger.Debug("re-identify window", entry.current.innerId)
				entry.innerId, newAppInfo = m.identifyWindow(entry.current)
				entry.setAppInfo(newAppInfo)
			}
		}
		entry.updateIcon()
		entry.setIsDocked(false)
		entry.updateName()
		entry.updateMenu()
	}
	m.saveDockedApps()
}
