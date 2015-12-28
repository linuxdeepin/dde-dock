package appearance

import (
	"os"
	"path"
	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/dde/daemon/appearance/subthemes"
	"strings"
	"time"
)

var (
	gtkDirs  []string
	iconDirs []string
	bgDirs   []string
)

func (m *Manager) handleThemeChanged() {
	if m.watcher == nil {
		return
	}

	m.watchGtkDirs()
	m.watchIconDirs()
	m.watchBgDirs()

	for {
		select {
		case <-m.endWatcher:
			return
		case err := <-m.watcher.Error:
			logger.Warning("Recieve file watcher error:", err)
			return
		case ev := <-m.watcher.Event:
			<-time.After(time.Millisecond * 800)
			file := ev.Name
			switch {
			case hasEventOccurred(file, bgDirs):
				background.RefreshBackground()
			case hasEventOccurred(file, gtkDirs):
				subthemes.RefreshGtkThemes()
			case hasEventOccurred(file, iconDirs):
				subthemes.RefreshIconThemes()
				subthemes.RefreshCursorThemes()
			}
		}
	}
}

func (m *Manager) watchGtkDirs() {
	var home = os.Getenv("HOME")
	gtkDirs = []string{
		path.Join(home, ".local/share/themes"),
		path.Join(home, ".themes"),
		"/usr/local/share/themes",
		"/usr/share/themes",
	}

	m.watchFiles(gtkDirs)
}

func (m *Manager) watchIconDirs() {
	var home = os.Getenv("HOME")
	iconDirs = []string{
		path.Join(home, ".local/share/icons"),
		path.Join(home, ".icons"),
		"/usr/local/share/icons",
		"/usr/share/icons",
	}

	m.watchFiles(iconDirs)
}

func (m *Manager) watchBgDirs() {
	bgDirs = background.ListDirs()
	m.watchFiles(bgDirs)
}

func (m *Manager) watchFiles(files []string) {
	for _, file := range files {
		err := m.watcher.Watch(file)
		if err != nil {
			logger.Debugf("Watch file '%s' failed: %v", file, err)
		}
	}
}

func hasEventOccurred(ev string, list []string) bool {
	for _, v := range list {
		if strings.Contains(ev, v) {
			return true
		}
	}
	return false
}
