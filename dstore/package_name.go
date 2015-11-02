package dstore

import (
	"bufio"
	"encoding/json"
	"os"
)

type DQueryPkgNameTransaction struct {
	data map[string]string
}

// NewDQueryPkgNameTransaction returns package name of given desktop file.
func NewDQueryPkgNameTransaction(path string) (*DQueryPkgNameTransaction, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	transition := &DQueryPkgNameTransaction{data: map[string]string{}}
	decoder := json.NewDecoder(bufio.NewReader(f))
	err = decoder.Decode(&transition.data)
	if err != nil {
		return nil, err
	}
	return transition, nil
}

func (transition *DQueryPkgNameTransaction) Query(desktopID string) string {
	if transition.data != nil {
		pkg := transition.data[desktopID]
		return pkg
	}
	return ""
}
