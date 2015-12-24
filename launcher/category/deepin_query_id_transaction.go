package category

import (
	"fmt"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
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

func (t *DeepinQueryIDTransaction) Query(pkgName string) (CategoryID, error) {
	cidName, ok := t.pkgToCategory[pkgName]
	if !ok {
		return OthersID, fmt.Errorf("No such a category for package %q", pkgName)
	}
	return getCategoryID(cidName)
}

func (t *DeepinQueryIDTransaction) Free() {
}
