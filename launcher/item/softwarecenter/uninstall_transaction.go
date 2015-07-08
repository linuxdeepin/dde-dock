package softwarecenter

import (
	"errors"
	"fmt"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"strings"
	"time"
)

// TODO: add Cancel
type UninstallTransaction struct {
	pkgName         string
	purge           bool
	timeoutDuration time.Duration
	timeout         <-chan time.Time
	done            chan struct{}
	failed          chan error
	soft            SoftwareCenterInterface
	disconnect      func()
}

func NewUninstallTransaction(soft SoftwareCenterInterface, pkgName string, purge bool, timeout time.Duration) *UninstallTransaction {
	return &UninstallTransaction{
		pkgName:         pkgName,
		purge:           purge,
		timeoutDuration: timeout,
		timeout:         nil,
		done:            make(chan struct{}, 1),
		failed:          make(chan error, 1),
		soft:            soft,
	}
}

func (t *UninstallTransaction) run() {
	t.disconnect = t.soft.Connectupdate_signal(func(message [][]interface{}) {
		switch message[0][0].(string) {
		case ActionStart, ActionUpdate, ActionFinish, ActionFailed:
			msgs := UpdateSignalTranslator(message)
			for _, action := range msgs {
				if action.Name == ActionFailed {
					detail := action.Detail.Value().(ActionFailedDetail)
					if strings.TrimRight(detail.PkgName, ":i386") == t.pkgName {
						err := fmt.Errorf("uninstall %q failed: %s", detail.PkgName, detail.Description)
						t.failed <- err
						t.disconnect()
						return
					}
				} else if action.Name == ActionFinish {
					defer func() {
						if err := recover(); err != nil {
							fmt.Println(err)
						}
					}()
					detail := action.Detail.Value().(ActionFinishDetail)
					if strings.TrimRight(detail.PkgName, ":i386") == t.pkgName {
						close(t.done)
						t.disconnect()
						return
					}
				}
			}
		default:
			t.disconnect()
			return
		}
	})

	t.timeout = time.After(t.timeoutDuration)
	if err := t.soft.UninstallPkg(t.pkgName, t.purge); err != nil {
		t.failed <- err
	}
}

func (t *UninstallTransaction) wait() error {
	select {
	case <-t.done:
		return nil
	case err := <-t.failed:
		return err
	case <-t.timeout:
		return errors.New("timeout")
	}
}

func (t *UninstallTransaction) Exec() error {
	t.run()
	return t.wait()
}
