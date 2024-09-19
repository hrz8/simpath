package database

import (
	"database/sql"
	"embed"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func ConnectDB(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	goose.SetLogger(log.New(log.Writer(), "[goose] ", log.LstdFlags))
	goose.SetTableName("schema_migrations")
	goose.SetDialect("postgres")
	goose.SetBaseFS(MigrationsFS)

	return goose.Up(db, "migrations")
}
