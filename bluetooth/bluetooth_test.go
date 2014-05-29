package bluetooth

import (
	"testing"
	"time"
)

func TestModuleStartStop(t *testing.T) {
	// logger.SetLogLevel(liblogger.LEVEL_DEBUG)
	Stop()
	Stop()
	go func() {
		time.Sleep(30 * time.Second)
		Start()
		Stop()
		Stop()
	}()
	Start()
	Stop()
	Stop()
	time.Sleep(30 * time.Second)
}
