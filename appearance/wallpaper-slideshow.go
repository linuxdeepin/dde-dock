package appearance

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"pkg.deepin.io/dde/daemon/appearance/background"

	"pkg.deepin.io/lib/utils"
)

type changeBgFunc func(monitorSpace string, t time.Time)

// wallpaper slideshow scheduler
type WSScheduler struct {
	mu              sync.Mutex
	lastSetBg       time.Time //last set wallpaper time
	interval        time.Duration
	quit            chan chan struct{}
	intervalChanged chan struct{}
	running         bool
	fn              changeBgFunc
}

func newWSScheduler(fun changeBgFunc) *WSScheduler {
	s := &WSScheduler{
		quit:            make(chan chan struct{}),
		intervalChanged: make(chan struct{}),
		fn:              fun,
	}
	return s
}

func (s *WSScheduler) remainDuration() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	elapsed := time.Since(s.lastSetBg)
	if elapsed < 0 {
		elapsed = 0
	}
	result := s.interval - elapsed
	return result
}

func (s *WSScheduler) loopCheck(mointorSpace string) {
	s.running = true
	for {
		select {
		case <-s.intervalChanged:
			continue
		case t := <-time.After(s.remainDuration()):
			if s.fn != nil {
				go s.fn(mointorSpace, t)
			}
			s.lastSetBg = t
		case ch := <-s.quit:
			s.running = false
			close(ch)
			return
		}
	}
}

func (s *WSScheduler) updateInterval(monitorSpace string, v time.Duration) {
	if v < time.Second {
		v = time.Second
	}
	logger.Debug("update interval", v)

	s.mu.Lock()
	s.interval = v
	if s.running {
		s.mu.Unlock()
		s.intervalChanged <- struct{}{}
		return
	}

	s.mu.Unlock()
	go s.loopCheck(monitorSpace)
}

func (s *WSScheduler) stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	logger.Debug("stop")
	ch := make(chan struct{})
	s.quit <- ch
	// wait quit
	<-ch
}

// wallpaper slideshow config
type WSConfig struct {
	LastChange time.Time
	Showed     []string
}

func loadWSConfig(filename string) (mapMonitorWorkspaceWSConfig, error) {
	var cfg mapMonitorWorkspaceWSConfig
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadWSConfigSafe(filename string) mapMonitorWorkspaceWSConfig {
	cfg, _ := loadWSConfig(filename)
	return cfg
}

func (c mapMonitorWorkspaceWSConfig) save(filename string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	dir := filepath.Dir(filename)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	return err
}

// wallpaper slideshow loop
type WSLoop struct {
	mu        sync.Mutex
	rand      *rand.Rand
	showed    map[string]struct{}
	all       []string
	fsChanged bool
}

func newWSLoop() *WSLoop {
	return &WSLoop{
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
		showed:    make(map[string]struct{}),
		fsChanged: true,
	}
}

func (wrl *WSLoop) GetShowed() []string {
	wrl.mu.Lock()

	result := make([]string, 0, len(wrl.showed))
	for file := range wrl.showed {
		result = append(result, file)
	}

	wrl.mu.Unlock()
	return result
}

func (wrl *WSLoop) getNotShowed() []string {
	if wrl.fsChanged {
		bgs := background.ListBackground()
		bgFiles := make([]string, 0, len(bgs))
		for _, bg := range bgs {
			bgFiles = append(bgFiles, utils.DecodeURI(bg.Id))
		}
		wrl.all = bgFiles
	}

	var result []string
	for _, file := range wrl.all {
		_, ok := wrl.showed[file]
		if !ok {
			result = append(result, file)
		}
	}
	return result
}

func (wrl *WSLoop) getNext() string {
	notShowed := wrl.getNotShowed()
	if len(notShowed) == 0 {
		return ""
	}
	idx := wrl.rand.Intn(len(notShowed))
	next := notShowed[idx]
	wrl.showed[next] = struct{}{}
	return next
}

func (wrl *WSLoop) reset() {
	logger.Debug("WSLoop.reset")
	wrl.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	// clear showed
	for key := range wrl.showed {
		delete(wrl.showed, key)
	}
}

func (wrl *WSLoop) AddToShowed(file string) {
	file = utils.DecodeURI(file)
	wrl.mu.Lock()
	wrl.showed[file] = struct{}{}
	wrl.mu.Unlock()
}

func (wrl *WSLoop) GetNext() string {
	wrl.mu.Lock()
	defer wrl.mu.Unlock()

	next := wrl.getNext()
	if next != "" {
		return next
	}

	if len(wrl.all) > 0 {
		wrl.reset()
		next = wrl.getNext()
	}

	return next
}

func (wrl *WSLoop) NotifyFsChanged() {
	wrl.mu.Lock()
	wrl.fsChanged = true
	wrl.mu.Unlock()
}

func isValidWSPolicy(policy string) bool {
	if policy == wsPolicyWakeup || policy == wsPolicyLogin || policy == "" {
		return true
	}

	_, err := strconv.ParseUint(policy, 10, 32)
	return err == nil
}
