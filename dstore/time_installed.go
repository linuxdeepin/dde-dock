package dstore

import (
	"bufio"
	"encoding/json"
	"os"
)

type DQueryTimeInstalledTransaction struct {
	data map[string]int64
}

func NewDQueryTimeInstalledTransaction(file string) (*DQueryTimeInstalledTransaction, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	transition := &DQueryTimeInstalledTransaction{data: map[string]int64{}}
	decoder := json.NewDecoder(bufio.NewReader(f))
	err = decoder.Decode(&transition.data)
	if err != nil {
		return nil, err
	}
	return transition, nil
}

func (transition *DQueryTimeInstalledTransaction) Query(pkgName string) int64 {
	timeInstalled := transition.data[pkgName]
	return timeInstalled
}
