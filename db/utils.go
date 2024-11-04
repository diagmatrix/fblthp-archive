package db

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
// Constants

type MigrationOperation string
type MigrationDirection string

const (
	Upgrade   MigrationOperation = "up"
	Downgrade MigrationOperation = "down"
	Status    MigrationOperation = "status"
	Generate  MigrationOperation = "generate"
	Help      MigrationOperation = "help"
)

const (
	Head MigrationDirection = "head"
	Tail MigrationDirection = "tail"
	Next MigrationDirection = "next"
	Prev MigrationDirection = "prev"
	None MigrationDirection = "none"
)

const MIGRATION_DIRECTORY string = "db/migrations"
const MIGRATION_IDENTIFIER string = "M-"

// ---------------------------------------------------------------------------------------------------------------------
// Helpers

func getNewMigrations(localMigrations, dbMigrations []Migration) ([]Migration, error) {
	log.Println("Getting new migrations ...")
	// Check existing migrations
	for _, dbMigration := range dbMigrations {
		found := false
		for _, localMigrations := range localMigrations {
			if dbMigration.Name == localMigrations.Name {
				found = true
			}
		}
		if !found {
			return nil, errors.New("Migration not found: " + dbMigration.Name) // TODO: Implement custom error
		}
	}
	// Check new migrations
	newMigrations := []Migration{}
	for _, localMigration := range localMigrations {
		found := false
		for _, dbMigration := range dbMigrations {
			if localMigration.Name == dbMigration.Name {
				found = true
			}
		}
		if !found {
			newMigrations = append(newMigrations, localMigration)
		}
	}
	log.Println("New migrations:", newMigrations)
	return newMigrations, nil
}

func parseMigrationName(name string) (string, error) {
	parts := strings.Split(name, "M-")
	if len(parts) < 2 {
		return "", errors.New("Invalid migration name")
	}
	return parts[1], nil
}

func printMigrations(executedMigrations, pendingMigrations []Migration) {
	// Get max lengths
	maxNameLength := 0
	maxIDLength := 0
	for _, migration := range append(executedMigrations, pendingMigrations...) {
		if len(migration.Name) > maxNameLength {
			maxNameLength = len(migration.Name)
		}
		strID := strconv.Itoa(int(migration.ID))
		if len(strID) > maxIDLength {
			maxIDLength = len(strID)
		}
	}
	maxIDLength += 2
	maxNameLength += 2
	maxExecutedLength := 12

	// Get print strings
	title := "|" + centerString("ID", maxIDLength) + "|" + centerString("Name", maxNameLength) + "|" + centerString("Executed", maxExecutedLength) + "|"
	executedRows := migrationsToString(executedMigrations, true, maxNameLength, maxIDLength, maxExecutedLength)
	newRows := migrationsToString(pendingMigrations, false, maxNameLength, maxIDLength, maxExecutedLength)
	separator := getSeparator(maxIDLength, maxNameLength, maxExecutedLength)
	table := "\n" + separator + "\n" + title + "\n" + separator + "\n" + executedRows + newRows

	log.Println(table)
}

func migrationsToString(migrations []Migration, executed bool, maxNameLength, maxIDLength, maxExecutedLength int) string {
	var rows string
	for _, migration := range migrations {
		rows += "|" + centerString(strconv.Itoa(int(migration.ID)), maxIDLength) + "|" + centerString(migration.Name, maxNameLength) + "|" + centerString(strconv.FormatBool(executed), maxExecutedLength) + "|\n"
		rows += getSeparator(maxIDLength, maxNameLength, maxExecutedLength) + "\n"
	}
	return rows
}

func getSeparator(maxIDLength, maxNameLength, maxExecutedLength int) string {
	return "+" + strings.Repeat("-", maxIDLength) + "+" + strings.Repeat("-", maxNameLength) + "+" + strings.Repeat("-", maxExecutedLength) + "+"
}

func centerString(s string, width int) string {
	if len(s) >= width {
		return s
	} else {
		padding := width - len(s)
		leftPadding := padding / 2
		rightPadding := padding - leftPadding
		return strings.Repeat(" ", leftPadding) + s + strings.Repeat(" ", rightPadding)
	}
}

func reverse(array []any) {
	for i, j := 0, len(array)-1; i < j; i, j = i+1, j-1 {
		array[i], array[j] = array[j], array[i]
	}
}
