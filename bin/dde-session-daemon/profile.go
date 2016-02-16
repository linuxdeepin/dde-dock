/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"os"
	"os/signal"
	"runtime/pprof"
)

func createFile(name string) *os.File {
	f, err := os.Create(name)
	if err != nil {
		logger.Fatal(err)
	}
	return f
}

func startProfile(fn func(c <-chan os.Signal)) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			switch sig.String() {
			case "interrupt":
				fn(c)
				close(c)
				os.Exit(0)
			}
		}
	}()
}

func startMemProfile(name string) {
	logger.Info("Start memory profile")
	f := createFile(name)
	startProfile(func(c <-chan os.Signal) {
		logger.Info("Memory profile done.")
		pprof.WriteHeapProfile(f)
		f.Close()
	})
}

func startCPUProfile(name string) {
	logger.Info("Start CPU profile")
	f := createFile(name)
	pprof.StartCPUProfile(f)
	startProfile(func(c <-chan os.Signal) {
		logger.Info("CPU profile done.")
		pprof.StopCPUProfile()
		f.Close()
	})
}
