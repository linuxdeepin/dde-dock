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
	fs, err := utils.QueryFilesytemInfo(os.Getenv("HOME"))
	if err != nil {
		logger.Error("Failed to get filesystem info:", err)
		return err
	}
	if fs.AvailSize > fsMinLeftSpace {
		return nil
	}

	go func() {
		time.Sleep(time.Second * 30)
		sendNotify("dialog-warning", "",
			Tr("Insufficient disk space, please clean up in time!"))
	}()
	return nil
}

func (*Daemon) Stop() error {
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
