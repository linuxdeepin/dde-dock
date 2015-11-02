package dstore

import (
	"time"
)

// DUninstallTransaction is command object for uninstalling package.
// TODO: add Cancel
type DUninstallTransaction struct {
	pkgName         string
	purge           bool // this is useless for new interface.
	timeoutDuration time.Duration
	timeout         <-chan time.Time
	done            chan error
	disconnect      func()
}

// NewDUninstallTransaction creates a new DUninstallTransaction.
func NewDUninstallTransaction(pkgName string, purge bool, timeout time.Duration) *DUninstallTransaction {
	return &DUninstallTransaction{
		pkgName:         pkgName,
		purge:           purge,
		timeoutDuration: timeout,
		timeout:         nil,
		done:            make(chan error, 1),
	}
}

func (t *DUninstallTransaction) run() {
	proxy, err := newDStoreManager()
	if err != nil {
		t.done <- err
		return
	}
	defer destroyDStoreManager(proxy)

	t.timeout = time.After(time.Second * t.timeoutDuration)
	jobPath, err := proxy.RemovePackage("", t.pkgName)
	if err != nil {
		t.done <- err
		return
	}

	go waitJobDone(jobPath, jobTypeRemove, t.timeout, &(t.done))
}

func (t *DUninstallTransaction) wait() error {
	err := <-t.done
	close(t.done)
	return err
}

// Exec executes this transaction.
func (t *DUninstallTransaction) Exec() error {
	t.run()
	return t.wait()
}
