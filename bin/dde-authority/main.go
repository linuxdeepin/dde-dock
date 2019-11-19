package main

import (
	"flag"
	"time"

	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

var noQuitFlag bool
var logger = log.NewLogger("dde-authority")

func init() {
	flag.BoolVar(&noQuitFlag, "no-quit", false, "do not auto quit")
}

const (
	dbusInterface   = "com.deepin.daemon.Authority"
	dbusServiceName = dbusInterface
	dbusPath        = "/com/deepin/daemon/Authority"

	dbusAgentInterface = dbusInterface + ".Agent"
)

func main() {
	flag.Parse()
	service, err := dbusutil.NewSystemService()
	if err != nil {
		logger.Fatal(err)
	}

	auth := newAuthority(service)
	err = service.Export(dbusPath, auth)
	if err != nil {
		logger.Fatal(err)
	}
	err = service.RequestName(dbusServiceName)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Debug("start service")
	if !noQuitFlag {
		service.SetAutoQuitHandler(3*time.Minute, func() bool {
			auth.mu.Lock()
			canQuit := len(auth.txs) == 0
			auth.mu.Unlock()
			return canQuit
		})
	}
	service.Wait()
}
