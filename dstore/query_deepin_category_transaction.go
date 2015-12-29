package dstore

import (
	"fmt"
)

type DeepinQueryIDTransaction struct {
	pkgToCategory map[string]string
}

func NewDeepinQueryIDTransaction(pkgToCategoryFile string) (*DeepinQueryIDTransaction, error) {
	t := &DeepinQueryIDTransaction{
		pkgToCategory: map[string]string{},
	}

	var err error
	t.pkgToCategory, err = getCategoryInfo(pkgToCategoryFile)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *DeepinQueryIDTransaction) Query(pkgName string) (string, error) {
	cidName, ok := t.pkgToCategory[pkgName]
	if !ok {
		return OthersName, fmt.Errorf("No such a category for package %q", pkgName)
	}
	return cidName, nil
}
