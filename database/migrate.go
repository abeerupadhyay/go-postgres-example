package database

import (
    "fmt"
    "os"
    "strings"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)


func runMigrations(connUri pgConnectionUri, mode string) error {

    filePath := fmt.Sprintf("file://%s", os.Getenv("PG_MIGRATIONS_FILE_PATH"))
    m, err := migrate.New(filePath, string(connUri))
    if err != nil { panic(err) }

    defer m.Close()

    mode = strings.ToLower(mode)
    if mode == "up" {
        err = m.Up()
    } else if mode == "down" {
        err = m.Down()
    } else {
        err = fmt.Errorf("invalid migration mode '%v'", mode)
    }

    // DO NOT consider ErrNoChange as migration failure.
    // Simply ignore it.
    if err == migrate.ErrNoChange { err = nil }
    return err
}

// Apply all migrations under migrations file path
func RunMigrations() error {

    // Migrations is always performed on default database
    connUri := pgConnectionUri(os.Getenv("PG_DEFAULT_CONN_URI"))
    err := runMigrations(connUri, "up")
    if err != nil {
        pgLogger.Panicf("Migrations failed! Error: %v", err)
        return err
    } else {
        pgLogger.Infof("Migrations applied successfully")
        return nil
    }
}

// Revert all migrations under migrations file path
// NOTE: This function is used only while testing to cleanup database
func RevertMigrations() error {

    // Migrations is always performed on default database
    connUri := pgConnectionUri(os.Getenv("PG_DEFAULT_CONN_URI"))
    err := runMigrations(connUri, "down")
    if err != nil {
        pgLogger.Panicf("Migrations drop failed! Error: %v", err)
        return err
    } else {
        pgLogger.Infof("Migrations dropped successfully")
        return nil
    }
}
