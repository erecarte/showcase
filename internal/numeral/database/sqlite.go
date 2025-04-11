package database

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var (
	DbName     = "foo"
	DiskDbPath = fmt.Sprintf("%s.db", DbName)
	DiskDbURI  = fmt.Sprintf("sqlite://./%s", DiskDbPath)
	MemDbURI   = fmt.Sprintf("file:%s?mode=memory&cache=shared", DbName)
)

func GetConnection(dbPath string) (*sql.DB, error) {
	migrateDB(dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func migrateDB(diskDbPath string) {
	dbDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		log.Println(err)
		return
	}
	dbURI := fmt.Sprintf("sqlite://%s", diskDbPath)
	migrations, err := migrate.NewWithSourceInstance("iofs", dbDriver, dbURI)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Running migrations")
	err = migrations.Up()
	if err != nil && err.Error() != "no change" {
		log.Println(err)
		return
	}
}
