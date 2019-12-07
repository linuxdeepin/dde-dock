package apps

import (
	"os"
	"sync"
	"time"

	"pkg.deepin.io/lib/fsnotify"
	"pkg.deepin.io/lib/strv"
)

type fsWatcher struct {
	*fsnotify.Watcher
	try          strv.Strv
	roots        strv.Strv
	timer        *time.Timer
	interval     time.Duration
	mu           sync.Mutex
	trySuccessCb func(string)
}

func newFsWatcher(interval time.Duration) (*fsWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	timer := time.NewTimer(interval)
	timer.Stop()
	w := &fsWatcher{
		Watcher:  watcher,
		timer:    timer,
		interval: interval,
	}
	go w.loopCheck()
	return w, nil
}

func (w *fsWatcher) loopCheck() {
	for {
		<-w.timer.C
		var newTryFiles []string

		w.mu.Lock()
		for _, file := range w.try {
			err := w.Watch(file)
			//logger.Debug("try watch", file)
			if os.IsNotExist(err) {
				newTryFiles = append(newTryFiles, file)
			} else if err == nil {
				logger.Debug("watch success", file)
				if w.trySuccessCb != nil {
					w.trySuccessCb(file)
				} else {
					logger.Warning("fsWatcher.trySuccessCb is nil")
				}
			}
		}

		w.try = newTryFiles

		if len(newTryFiles) > 0 {
			//logger.Debug(newTryFiles)
			w.timer.Reset(w.interval)
			//logger.Debug("reset timer")
		} else {
			logger.Debug("stop timer")
		}
		w.mu.Unlock()
	}
}

func (w *fsWatcher) addRoot(root string) {
	logger.Debug("fsWatcher.addRoot", root)
	w.mu.Lock()
	var added bool
	w.roots, added = w.roots.Add(root)
	w.mu.Unlock()

	if !added {
		return
	}

	err := w.Watch(root)
	if os.IsNotExist(err) {
		w.addTry(root)
	}
}

func (w *fsWatcher) removeRoot(root string) {
	logger.Debug("fsWatcher.removeRoot", root)
	w.mu.Lock()
	var deleted bool
	w.roots, deleted = w.roots.Delete(root)
	w.mu.Unlock()

	if !deleted {
		return
	}

	w.removeTry(root)
}

func (w *fsWatcher) addTry(file string) {
	w.mu.Lock()
	w.try, _ = w.try.Add(file)
	if len(w.try) == 1 {
		w.timer.Reset(w.interval)
	}
	w.mu.Unlock()
}

func (w *fsWatcher) removeTry(file string) {
	w.mu.Lock()
	w.try, _ = w.try.Delete(file)
	w.mu.Unlock()
}

func (w *fsWatcher) handleEvent(event *fsnotify.FileEvent) {
	if event.IsDelete() && strv.Strv(w.roots).Contains(event.Name) {
		w.addTry(event.Name)
	}
}
