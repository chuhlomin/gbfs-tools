package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	// Import go-sqlite3 library
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/chuhlomin/gbfs-tools/pkg/structs"
)

const createSystemsTableSQL = `CREATE TABLE systems (
	"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"system_id" TEXT,
	"country_code" TEXT,
	"name" TEXT,
	"location" TEXT,
	"url" TEXT,
	"auto_discovery_url" TEXT,
	"is_enabled" BOOL NOT NULL DEFAULT 'true'
);`

type SQLite struct {
	db *sqlx.DB
}

func NewSQLite(pathToDatabase string) (*SQLite, error) {
	db, err := sqlx.Open("sqlite3", pathToDatabase)
	if err != nil {
		return nil, errors.Wrap(err, "open sqlite database")
	}

	_, err = db.Query("select 1 from systems;")
	if err != nil {
		log.Println("No systems table found in database, creating...")
		statement, err := db.Prepare(createSystemsTableSQL)
		if err != nil {
			return nil, errors.Wrap(err, "prepare statement: create systems table")
		}
		_, err = statement.Exec()
		if err != nil {
			return nil, errors.Wrap(err, "execute statement: create systems table")
		}
	}

	return &SQLite{db: db}, nil
}

func (sql *SQLite) AddSystem(system structs.System) error {
	_, err := sql.db.NamedExec(
		`INSERT INTO systems (system_id, country_code, name, location, url, auto_discovery_url)
			VALUES (:system_id, :country_code, :name, :location, :url, :auto_discovery_url)`,
		&system,
	)

	return err
}

func (sql *SQLite) DisableSystem(id string) error {
	_, err := sql.db.Exec(
		`UPDATE systems SET is_enabled = false WHERE system_id = $1`,
		id,
	)

	return err
}

func (sql *SQLite) GetSystems() ([]structs.System, error) {
	rows, err := sql.db.Queryx(
		`SELECT system_id, country_code, name, location, url, auto_discovery_url, is_enabled
			FROM systems`,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select systems")
	}

	var result []structs.System
	for rows.Next() {
		system := structs.System{}
		err := rows.StructScan(&system)
		if err != nil {
			return nil, errors.Wrap(err, "scan systems")
		}
		result = append(result, system)
	}

	return result, nil
}

func (sql *SQLite) Close() error {
	return sql.db.Close()
}
