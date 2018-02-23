package calltrace

import (
	"fmt"
	"os"
	"path/filepath"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/xdg/basedir"
	"runtime/pprof"
	"time"
)

type CallTrace struct {
	cpuFile   *os.File
	stackFile *os.File

	duration uint32

	quit   chan bool
	logger *log.Logger
}

// Start launch calltrace module
func Start(duration uint32, l *log.Logger) (*CallTrace, error) {
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
		cpu.Close()
		return nil, err
	}

	ct := new(CallTrace)
	ct.cpuFile = cpu
	ct.stackFile = stack
	ct.logger = l
	ct.quit = make(chan bool)
	ct.duration = duration

	err = pprof.StartCPUProfile(ct.cpuFile)
	if err != nil {
		ct.logger.Warning("Failed to start cpu profile:", err)
		ct.cpuFile.Close()
		ct.cpuFile = nil
	}

	ct.logger.Infof("[CallTrace] Will record profiles, once per %d second", ct.duration)
	ct.writeHeap()
	ct.recordStack()
	go ct.loop()

	return ct, nil
}

func (ct *CallTrace) SetAutoDestroy(seconds uint32) {
	if seconds == 0 {
		return
	}

	go func() {
		time.Sleep(time.Second * time.Duration(seconds))
		ct.stop()
	}()
}

// Stop terminate calltrace module
func (ct *CallTrace) stop() {
	if ct.quit != nil {
		ct.quit <- true
	}
	if ct.cpuFile != nil {
		pprof.StopCPUProfile()
		ct.cpuFile.Close()
	}

	if ct.stackFile != nil {
		ct.stackFile.Close()
	}
	ct.logger.Info("[CallTrace] Terminated!")
	ct = nil
}

func (ct *CallTrace) loop() {
	var ticker = time.NewTicker(time.Second * time.Duration(ct.duration))
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				ct.logger.Error("Invalid ticker event")
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

func (ct *CallTrace) writeHeap() {
	dir, _ := ensureDirExist()
	mem, err := os.Create(filepath.Join(dir,
		fmt.Sprintf("memory_%v.prof", time.Now().UnixNano())))
	if err != nil {
		ct.logger.Warning("Failed to create memory file:", err)
		return
	}

	err = pprof.WriteHeapProfile(mem)
	mem.Close()
	if err != nil {
		ct.logger.Warning("Failed to write head profile:", err)
	}
	return
}

func (ct *CallTrace) recordStack() {
	_, err := ct.stackFile.WriteString("=== BEGIN " + time.Now().String() + " ===\n")
	if err != nil {
		ct.logger.Warning("Failed to start record stack:", err)
		return
	}
	ct.stackFile.WriteString("--- DUMP [goroutine] ---\n")
	pprof.Lookup("goroutine").WriteTo(ct.stackFile, 2)
	ct.stackFile.WriteString("\n--- DUMP [heap] ---\n")
	pprof.Lookup("heap").WriteTo(ct.stackFile, 2)
	ct.stackFile.WriteString("\n--- DUMP [threadcreate] ---\n")
	pprof.Lookup("threadcreate").WriteTo(ct.stackFile, 2)
	ct.stackFile.WriteString("\n--- DUMP [block] ---\n")
	pprof.Lookup("block").WriteTo(ct.stackFile, 2)
	ct.stackFile.WriteString("\n--- DUMP [mutex] ---\n")
	pprof.Lookup("mutex").WriteTo(ct.stackFile, 2)
	ct.stackFile.WriteString("=== END " + time.Now().String() + " ===\n\n\n")
	ct.stackFile.Sync()
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
