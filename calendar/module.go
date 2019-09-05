package calendar

import (
	"os"
	"path/filepath"
	"time"

	"pkg.deepin.io/lib/xdg/basedir"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/calendar")
var dbFile = filepath.Join(basedir.GetUserConfigDir(), "deepin/dde-daemon/calendar/scheduler.db")

func init() {
	loader.Register(newModule())
}

type Module struct {
	scheduler *Scheduler
	*loader.ModuleBase
}

func (m *Module) GetDependencies() []string {
	return nil
}

func (m *Module) Start() error {
	if m.scheduler != nil {
		return nil
	}
	go func() {
		t0 := time.Now()
		err := m.start()
		if err != nil {
			logger.Warning("failed to start calendar module:", err)
		}
		logger.Info("start calendar module cost", time.Since(t0))
	}()
	return nil
}

func (m *Module) Stop() error {
	if m.scheduler == nil {
		return nil
	}

	service := loader.GetService()
	err := service.ReleaseName(dbusServiceName)
	if err != nil {
		return err
	}

	err = service.StopExport(m.scheduler)
	if err != nil {
		return err
	}

	close(m.scheduler.quitChan)
	m.scheduler = nil
	return nil
}

func (m *Module) start() error {
	err := os.MkdirAll(filepath.Dir(dbFile), 0755)
	if err != nil {
		return err
	}

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	if logger.GetLogLevel() == log.LevelDebug {
		db = db.Debug()
	}
	err = db.AutoMigrate(&Job{}).Error
	if err != nil {
		logger.Warning(err)
	}

	hasJobTypeTable := db.HasTable(&JobType{})

	err = db.AutoMigrate(&JobType{}).Error
	if err != nil {
		logger.Warning(err)
	}

	if !hasJobTypeTable {
		err = initJobTypeTable(db)
		if err != nil {
			logger.Warning(err)
		}
	}

	service := loader.GetService()
	m.scheduler = newScheduler(db, service)

	err = service.Export(dbusPath, m.scheduler)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	m.scheduler.startRemindLoop()
	return nil
}

func newModule() *Module {
	m := new(Module)
	m.ModuleBase = loader.NewModuleBase("calendar", m, logger)
	return m
}

func initJobTypeTable(db *gorm.DB) error {
	work := gettext.Tr("Work")
	life := gettext.Tr("Life")
	other := gettext.Tr("Other")

	types := []*JobType{
		{
			Name:  work,
			Color: "#FF0000", // red
		},
		{
			Name:  life,
			Color: "#00FF00", // green
		},
		{
			Name:  other,
			Color: "#800080", // purple
		},
	}
	for idx, t := range types {
		t.ID = uint(idx) + 1
		err := db.Create(t).Error
		if err != nil {
			return err
		}
	}
	return nil
}
