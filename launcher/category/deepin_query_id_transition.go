package category

import (
	"database/sql"
	"errors"
	"path"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/lib/gio-2.0"
)

type DeepinQueryIDTransition struct {
	db *sql.DB
}

func NewDeepinQueryIDTransition(dbPath string) (*DeepinQueryIDTransition, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	transition := &DeepinQueryIDTransition{db: db}
	return transition, nil
}

func (transition *DeepinQueryIDTransition) Query(app *gio.DesktopAppInfo) (CategoryID, error) {
	db := transition.db
	if db == nil {
		return OthersID, errors.New("invalid db")
	}

	filename := app.GetFilename()
	basename := path.Base(filename)
	var categoryName string
	err := db.QueryRow(`
	select first_category_name
	from desktop
	where desktop_name = ?`,
		basename,
	).Scan(&categoryName)
	if err != nil {
		return OthersID, err
	}

	if categoryName == "" {
		return OthersID, errors.New("get empty category")
	}

	return getCategoryID(categoryName)
}

func (transition *DeepinQueryIDTransition) Free() {
	if transition.db != nil {
		transition.db.Close()
	}
}
