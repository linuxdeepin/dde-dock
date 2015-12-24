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
	t := &DQueryPkgNameTransaction{data: map[string]string{}}
	decoder := json.NewDecoder(bufio.NewReader(f))
	err = decoder.Decode(&t.data)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *DQueryPkgNameTransaction) Query(desktopID string) string {
	if t.data != nil {
		pkg := t.data[desktopID]
		return pkg
	}
	return ""
}
