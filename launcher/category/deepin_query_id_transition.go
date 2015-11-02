package category

import (
	"fmt"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

type DeepinQueryIDTransaction struct {
	pkgToCategory map[string]string
}

func NewDeepinQueryIDTransaction(pkgToCategoryFile string) (*DeepinQueryIDTransaction, error) {
	transition := &DeepinQueryIDTransaction{
		pkgToCategory: map[string]string{},
	}

	var err error
	transition.pkgToCategory, err = getCategoryInfo(pkgToCategoryFile)
	if err != nil {
		return nil, err
	}

	return transition, nil
}

func (transition *DeepinQueryIDTransaction) Query(pkgName string) (CategoryID, error) {
	cidName, ok := transition.pkgToCategory[pkgName]
	if !ok {
		return OthersID, fmt.Errorf("No such a category for package %q", pkgName)
	}
	return getCategoryID(cidName)
}

func (transition *DeepinQueryIDTransaction) Free() {
}
