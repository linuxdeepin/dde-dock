// profiling of your Go application. TODO: move to common lib.
package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
)

type _Profile struct {
	writer io.Writer
	name   string
}

func (prof *_Profile) Name() string {
	return prof.name
}

type _CPUProfile struct {
	_Profile
}

func newCPUProfile(name string) *_CPUProfile {
	return &_CPUProfile{_Profile{name: name}}
}

func (prof *_CPUProfile) Start(writer io.Writer) {
	pprof.StartCPUProfile(writer)
}

func (prof *_CPUProfile) Stop() {
	pprof.StopCPUProfile()
}

type _MemProfile struct {
	_Profile
}

func newMemPrifle(name string) *_MemProfile {
	return &_MemProfile{_Profile{name: name}}
}

func (prof *_MemProfile) Start(writer io.Writer) {
	prof.writer = writer
}

func (prof *_MemProfile) Stop() {
	pprof.Lookup("heap").WriteTo(prof.writer, 0)
}

type _BlockProfile struct {
	_Profile
}

func newBlockProfile(name string) *_BlockProfile {
	return &_BlockProfile{_Profile{name: name}}
}

func (prof *_BlockProfile) Start(writer io.Writer) {
	prof.writer = writer
}

func (prof *_BlockProfile) Stop() {
	pprof.Lookup("block").WriteTo(prof.writer, 0)
}

// Config controls the operation of the profile package.
type Config struct {
	// CPUProfile is the name of cpu profile which controls if cpu profiling will be enabled.
	// It defaults to false.
	CPUProfile string

	// MemProfile is the name of memory profile which controls if cpu profiling will be enabled.
	// It defaults to false.
	MemProfile string

	// MemProfile is the name of memory profile which controls if cpu profiling will be enabled.
	// It defaults to false.
	BlockProfile string

	// NoShutdownHook controls whether the profiling package should
	// hook SIGINT to write profiles cleanly.
	// Programs with more sophisticated signal handling should set
	// this to true and ensure the Stop() function returned from Start()
	// is called during shutdown.
	NoShutdownHook bool

	closers []func()
}

func (cfg *Config) enableProfile(prof interface {
	Name() string
	Start(io.Writer)
	Stop()
}) error {
	if prof.Name() == "" {
		return nil
	}

	var fn string
	var err error
	fn = prof.Name()
	profilePath := filepath.Dir(fn)
	if profilePath == "" {
		profilePath, err = ioutil.TempDir("", "profile")
	}
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(fn), 0777)
	if err != nil {
		return err
	}

	f, err := os.Create(fn)
	if err != nil {
		cfg.Stop()
		return err
	}

	prof.Start(f)
	cfg.closers = append(cfg.closers, func() {
		prof.Stop()
		f.Close()
	})

	return nil
}

// Start starts a new profiling session configured using *Config.
// The caller should call the Stop method to cleanly stop profiling.
func (cfg *Config) Start() error {
	var err error

	if err = cfg.enableProfile(newCPUProfile(cfg.CPUProfile)); err != nil {
		return err
	}
	if err = cfg.enableProfile(newMemPrifle(cfg.MemProfile)); err != nil {
		return err
	}
	if err = cfg.enableProfile(newBlockProfile(cfg.BlockProfile)); err != nil {
		return err
	}

	if !cfg.NoShutdownHook {
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c

			cfg.Stop()

			os.Exit(0)
		}()
	}
	return nil
}

// Stop stops all profile.
func (cfg *Config) Stop() {
	for _, c := range cfg.closers {
		c()
	}
}
