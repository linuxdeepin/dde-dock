package dstore

import (
	"gir/gio-2.0"
)

type QueryCategoryTransaction struct {
	pkgQuery       *DQueryPkgNameTransaction
	deepinQuery    *DeepinQueryIDTransaction
	xCategoryQuery *XCategoryQueryIDTransaction
}

func NewQueryCategoryTransaction(desktopToPkgFile string, appInfoFile string, xcategoryFile string) (*QueryCategoryTransaction, error) {
	t := &QueryCategoryTransaction{}
	var err1 error
	var err2 error

	t.pkgQuery, err1 = NewDQueryPkgNameTransaction(desktopToPkgFile)
	t.deepinQuery, err2 = NewDeepinQueryIDTransaction(appInfoFile)
	t.xCategoryQuery, _ = NewXCategoryQueryIDTransaction(xcategoryFile)

	if err1 != nil {
		return t, err1
	} else if err2 != nil {
		return t, err2
	}

	return t, nil
}

func (t *QueryCategoryTransaction) Query(app *gio.DesktopAppInfo) (string, error) {
	if t.pkgQuery != nil && t.deepinQuery != nil {
		pkgName := t.pkgQuery.Query(app.GetId())
		cid, err := t.deepinQuery.Query(pkgName)
		if err == nil {
			return cid, nil
		}
	}

	cid, err := t.xCategoryQuery.Query(app.GetCategories())
	return cid, err
}
