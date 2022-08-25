package clients

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	Exec(query string, args ...any) error
	QueryRow(destination interface{}, query string, args ...any) error
	IsEmptyResult(err error) bool
}

type database struct {
	db *sql.DB
}

const databaseFile string = ".github-notify.cache.db"

const create string = `
  CREATE TABLE IF NOT EXISTS cache (
  key TEXT NOT NULL PRIMARY KEY,
  value TEXT,
  createdAt DATETIME NOT NULL,
  deleteAt DATETIME
);`

func NewSQLiteClient() (Database, error) {
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(create); err != nil {
		return nil, err
	}

	return &database{
		db: db,
	}, nil
}

func (d *database) Exec(query string, args ...any) error {
	_, err := d.db.Exec(query, args...)

	return err
}

func (d *database) QueryRow(destination interface{}, query string, args ...any) error {
	row := d.db.QueryRow(query, args...)

	return row.Scan(destination)
}

func (d *database) IsEmptyResult(err error) bool {
	return err == sql.ErrNoRows
}
