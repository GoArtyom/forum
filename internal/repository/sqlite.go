package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum/config"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteDB(config *config.Config) (*sql.DB, error) {
	db, err := sql.Open(config.DB.Driver, config.DB.DSN)
	if err != nil {
		return nil, err
	}

	ctx, calcel := context.WithTimeout(context.Background(), time.Second*5)
	defer calcel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB, migrate string) error {
	_, err := db.Exec("PRAGMA foreign_keys=ON;")
	if err != nil {
		return err
	}
	file, err := os.ReadFile(migrate)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		fmt.Println(request)
		_, err := tx.Exec(request)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
