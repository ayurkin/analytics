//go:build approve || reject
// +build approve reject

package tests

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

const layout string = "2006-01-02 15:04:05.999999999 -0700 MST"

const (
	dbName = "app"
	dbUser = "app"
	dbPass = "secret"
)

func applyMigrations(connStr string) error {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	err = goose.Up(conn, "../../db/changelog/")
	if err != nil {
		return err
	}
	return nil
}
