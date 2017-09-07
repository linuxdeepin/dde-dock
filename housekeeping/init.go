package housekeeping

import (
	"dbus/org/freedesktop/notifications"
	"os"
	"pkg.deepin.io/dde/daemon/loader"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/utils"
	"time"
)

const (
	// 500MB
	fsMinLeftSpace = 1024 * 1024 * 500
)

func init() {
	loader.Register(NewDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
	ticker   *time.Ticker
	stopChan chan struct{}
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("housekeeping", daemon, logger)
	return daemon
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

var (
	logger = log.NewLogger("housekeeping")
)

func (d *Daemon) Start() error {
	if d.stopChan != nil {
		return nil
	}

	d.ticker = time.NewTicker(time.Minute * 1)
	d.stopChan = make(chan struct{})
	go func() {
		for {
			select {
			case <-d.ticker.C:
				fs, err := utils.QueryFilesytemInfo(os.Getenv("HOME"))
				if err != nil {
					logger.Error("Failed to get filesystem info:", err)
					break
				}
				logger.Debug("Home filesystem info(total, free, avail):",
					fs.TotalSize, fs.FreeSize, fs.AvailSize)
				if fs.AvailSize > fsMinLeftSpace {
					break
				}
				sendNotify("dialog-warning", "",
					Tr("Insufficient disk space, please clean up in time!"))
			case <-d.stopChan:
				logger.Debug("Stop housekeeping")
				if d.ticker != nil {
					d.ticker.Stop()
					d.ticker = nil
				}
				return
			}
		}
	}()
	return nil
}

func (d *Daemon) Stop() error {
	if d.stopChan != nil {
		close(d.stopChan)
		d.stopChan = nil
	}
	return nil
}

func sendNotify(icon, summary, body string) error {
	notifier, err := notifications.NewNotifier(
		"org.freedesktop.Notifications",
		"/org/freedesktop/Notifications")
	if err != nil {
		return err
	}

	_, err = notifier.Notify("housekeeping", 0,
		icon, summary, body,
		nil, nil, 0)

	return err
}
