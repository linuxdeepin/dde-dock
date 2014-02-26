package main

import (
	"database/sql"
	"fmt"
	"log"
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
	names, err := getPackageNamesFromDatabase(path)
	if err != nil {
		fmt.Println(err)
		names, err = getPackageNamesFromCommandline(path)
		if err != nil {
			fmt.Println(err)
			return names
		}
	}
	return names
}

func getPackageNamesFromDatabase(path string) ([]string, error) {
	dbPath, err := getDBPath(CategoryNameDBPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
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
	return strings.Split(names, " "), nil
}

func getPackageNamesFromCommandline(path string) ([]string, error) {
	names, err := exec.Command("dpkg", "-S", path).Output()
	if err != nil {
		return make([]string, 0), err
	}
	return strings.Split(string(names), " "), nil
}
