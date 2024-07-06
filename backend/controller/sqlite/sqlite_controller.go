package controller

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	ctrl "backend/controller"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	file string
	db   *sql.DB
}

type eventRow struct {
	id          int
	name        string
	description string
	venue       string
	date        time.Time
	kind        string
}

const DateFormat string = "2006-01-02T15:04:05Z"

func NewSqliteController(file string) (ctrl.Controller, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return Sqlite{}, err
	}
	return Sqlite{
		file,
		db,
	}, nil
}

func SqliteFromSqlFile(dbFile string, sqlFile string) (ctrl.Controller, error) {
	if _, err := os.Stat(dbFile); err == nil {
		return nil, fmt.Errorf("cannot create db on existing db %s", dbFile)
	}
	db, dbErr := sql.Open("sqlite3", dbFile)
	if dbErr != nil {
		return nil, dbErr
	}
	schema, e := os.ReadFile(sqlFile)
	if e != nil {
		return nil, e
	}
	statements := strings.Split(string(schema), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, er := db.Exec(stmt)
		if er != nil {
			return nil, er
		}
	}
	return Sqlite{
		dbFile,
		db,
	}, nil
}
