package main

import (
	"database/sql"
	// "fmt"
	"os/exec"
	p "path"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SOFTWARE_STATUS_CREATED  string = "created"
	SOFTWARE_STATUS_MODIFIED string = "updated"
	SOFTWARE_STATUS_DELETED  string = "deleted"
)

func getPackageNames(path string) []string {
	// logger.Info("getPackageNames")
	names, err := getPackageNamesFromDatabase(path)
	if err != nil {
		logger.Info(err)
		names, err = getPackageNamesFromCommandline(path)
		if err != nil {
			logger.Info(err)
			return names
		}
	}
	// logger.Info(names)
	return names
}

func getPackageNamesFromDatabase(path string) ([]string, error) {
	dbPath, err := getDBPath(CategoryNameDBPath)
	if err != nil {
		return make([]string, 0), err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return make([]string, 0), err
	}
	defer db.Close()

	basename := p.Base(path)
	var names string
	err = db.QueryRow(
		`select pkg_names
		from desktop
		where desktop_name = ?`,
		basename,
	).Scan(&names)
	if err != nil {
		return make([]string, 0), err
	}
	// logger.Info(names)
	return strings.Split(names, ","), nil
}

func getPackageNamesFromCommandline(path string) ([]string, error) {
	names, err := exec.Command("dpkg", "-S", path).Output()
	if err != nil {
		return make([]string, 0), err
	}
	name := strings.Split(string(names), ":")[0]
	return []string{name}, nil
}
