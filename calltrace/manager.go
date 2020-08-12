package calltrace

import (
	"fmt"
	"os"
	"path/filepath"
	"pkg.deepin.io/lib/xdg/basedir"
	"runtime/pprof"
	"time"
)

// Manager manage calltrace files
type Manager struct {
	cpuFile   *os.File
	stackFile *os.File

	duration uint32

	quit chan bool
}

// NewManager create Manager and launch calltrace module
func NewManager(duration uint32) (*Manager, error) {
	dir, err := ensureDirExist()
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().UnixNano()
	cpu, err := os.Create(filepath.Join(dir,
		fmt.Sprintf("cpu_%v.prof", timestamp)))
	if err != nil {
		return nil, err
	}

	stack, err := os.OpenFile(filepath.Join(dir,
		fmt.Sprintf("stack_%v.log", timestamp)),
		os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		_ = cpu.Close()
		return nil, err
	}

	ct := new(Manager)
	ct.cpuFile = cpu
	ct.stackFile = stack
	ct.quit = make(chan bool)
	ct.duration = duration

	err = pprof.StartCPUProfile(ct.cpuFile)
	if err != nil {
		logger.Warning("Failed to start cpu profile:", err)
		_ = ct.cpuFile.Close()
		ct.cpuFile = nil
	}

	logger.Infof("[Manager] Will record profiles, once per %d second", ct.duration)
	ct.writeHeap()
	ct.recordStack()
	go ct.loop()

	return ct, nil
}

// SetAutoDestroy auto destroy manager after the special seconds
func (ct *Manager) SetAutoDestroy(seconds uint32) {
	if seconds == 0 {
		return
	}

	go func() {
		time.Sleep(time.Second * time.Duration(seconds))
		ct.stop()
	}()
}

// Stop terminate calltrace module
func (ct *Manager) stop() {
	if ct.quit != nil {
		ct.quit <- true
	}
	if ct.cpuFile != nil {
		pprof.StopCPUProfile()
		_ = ct.cpuFile.Close()
	}

	if ct.stackFile != nil {
		_ = ct.stackFile.Close()
	}
	logger.Info("[Manager] Terminated!")
	ct = nil
}

func (ct *Manager) loop() {
	var ticker = time.NewTicker(time.Second * time.Duration(ct.duration))
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				logger.Error("Invalid ticker event")
				return
			}

			ct.writeHeap()
			ct.recordStack()
		case <-ct.quit:
			ticker.Stop()
			close(ct.quit)
			ct.quit = nil
			return
		}
	}
}

func (ct *Manager) writeHeap() {
	dir, _ := ensureDirExist()
	mem, err := os.Create(filepath.Join(dir,
		fmt.Sprintf("memory_%v.prof", time.Now().UnixNano())))
	if err != nil {
		logger.Warning("Failed to create memory file:", err)
		return
	}

	err = pprof.WriteHeapProfile(mem)
	_ = mem.Close()
	if err != nil {
		logger.Warning("Failed to write head profile:", err)
	}
}

func (ct *Manager) recordStack() {
	_, err := ct.stackFile.WriteString("=== BEGIN " + time.Now().String() + " ===\n")
	if err != nil {
		logger.Warning("Failed to start record stack:", err)
		return
	}
	_, _ = ct.stackFile.WriteString("--- DUMP [goroutine] ---\n")
	_ = pprof.Lookup("goroutine").WriteTo(ct.stackFile, 2)
	_, _ =ct.stackFile.WriteString("\n--- DUMP [heap] ---\n")
	_ = pprof.Lookup("heap").WriteTo(ct.stackFile, 2)
	_, _ =ct.stackFile.WriteString("\n--- DUMP [threadcreate] ---\n")
	_ = pprof.Lookup("threadcreate").WriteTo(ct.stackFile, 2)
	_, _ =ct.stackFile.WriteString("\n--- DUMP [block] ---\n")
	_ = pprof.Lookup("block").WriteTo(ct.stackFile, 2)
	_, _ =ct.stackFile.WriteString("\n--- DUMP [mutex] ---\n")
	_ = pprof.Lookup("mutex").WriteTo(ct.stackFile, 2)
	_, _ =ct.stackFile.WriteString("=== END " + time.Now().String() + " ===\n\n\n")
	_ = ct.stackFile.Sync()
}

func ensureDirExist() (string, error) {
	dir := filepath.Join(basedir.GetUserCacheDir(),
		"deepin", "dde-daemon", "calltrace")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}
	return dir, nil
}
