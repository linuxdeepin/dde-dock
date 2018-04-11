package main

import (
	"flag"
	"log"
	"time"

	"pkg.deepin.io/lib/dbusutil"
)

var noQuitFlag bool

func init() {
	log.SetFlags(log.Lshortfile)
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
		log.Fatal(err)
	}

	auth := newAuthority(service)
	err = service.Export(dbusPath, auth)
	if err != nil {
		log.Fatal(err)
	}
	err = service.RequestName(dbusServiceName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("start service")
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
