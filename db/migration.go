package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/diagmatrix/fblthp-archive/exceptions"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Constants

type MigrationOperation string

const (
	Upgrade   MigrationOperation = "up"
	Downgrade MigrationOperation = "down"
	Status    MigrationOperation = "status"
	Generate  MigrationOperation = "generate"
	Help      MigrationOperation = "help"
)

const MIGRATION_DIRECTORY string = "db/migrations"
const MIGRATION_IDENTIFIER string = "M-"

type MigrationDirection string

const (
	Head MigrationDirection = "head"
	Tail MigrationDirection = "tail"
	Next MigrationDirection = "next"
	Prev MigrationDirection = "prev"
	None MigrationDirection = "none"
)

// ---------------------------------------------------------------------------------------------------------------------
// Migrations

type Migration struct {
	ID        uint
	Name      string
	CreatedAt time.Time
}

func (m *Migration) GetOperation(operation MigrationOperation) (string, error) {
	fileName := MIGRATION_DIRECTORY + "/" + MIGRATION_IDENTIFIER + m.Name + "/" + string(operation) + ".sql"
	file, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Migration Manager

type MigrationManager struct {
	DB          *sql.DB
	LastVersion uint
}

func NewMigrationManager(db *sql.DB) *MigrationManager {
	log.Println("Initiating migration manager...")
	// Test connection
	err := db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	var lastMigration uint
	err = db.QueryRow(SELECT_CURRENT_MIGRATION).Scan(&lastMigration)
	if err != nil {
		log.Println("No current migration found")
		lastMigration = 0
	}
	return &MigrationManager{DB: db, LastVersion: lastMigration}
}

func (m *MigrationManager) Process(operation MigrationOperation, argument string) error {
	log.Println("Processing " + operation + "...")
	switch operation {
	case Upgrade:
		if argument == "" {
			return errors.New("No operation provided") // TODO: Implement custom error
		}
		return m.upgrade(MigrationDirection(argument))
	case Downgrade:
		return exceptions.NewNotImplementedError()
	case Status:
		return m.status()
	case Generate:
		return m.generate(argument)
	case Help:
		return m.help()
	default:
		return m.help()
	}
}

func (m *MigrationManager) upgrade(direction MigrationDirection) error {
	log.Println("Upgrading to " + direction + "...")
	switch direction {
	case Tail:
		return errors.New("Cannot upgrade to tail") // TODO: Implement custom error
	case Head:
		localMigrations, dbMigrations, err := m.scan()
		if err != nil {
			return err
		}
		newMigrations, err := m.getNewMigrations(localMigrations, dbMigrations)
		if err != nil {
			return err
		}
		for _, migration := range newMigrations {
			upgradeSQL, err := migration.GetOperation(Upgrade)
			if err != nil {
				return err
			}
			_, err = m.DB.Exec(upgradeSQL)
			if err != nil {
				return err
			}
		}
		return nil
	case Next:
		return exceptions.NewNotImplementedError()
	case Prev:
		return exceptions.NewNotImplementedError()
	default:
		return errors.New("Invalid direction") // TODO: Implement custom error
	}
}

func (m *MigrationManager) status() error {
	log.Println("Getting migration status...")
	localMigrations, dbMigrations, err := m.scan()
	if err != nil {
		return err
	}
	newMigrations, err := m.getNewMigrations(localMigrations, dbMigrations)
	if err != nil {
		return err
	}
	printMigrations(dbMigrations, newMigrations)
	return nil
}

func (m *MigrationManager) generate(name string) error {
	log.Println("Generating migration skeleton...")
	if name == "" {
		name = "NewMigration"
	}
	today := time.Now().Format("20060102")
	dirName := fmt.Sprintf("%s/%s%s-%s", MIGRATION_DIRECTORY, MIGRATION_IDENTIFIER, today, name)
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}
	files := []string{"up.sql", "down.sql"}
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", dirName, file)
		err = os.WriteFile(filePath, []byte("-- Add migration SQL code here"), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MigrationManager) help() error {
	log.Println("Showing help...")
	// TODO: Implement help
	return exceptions.NewNotImplementedError()
}

func (m *MigrationManager) getNewMigrations(localMigrations, dbMigrations []Migration) ([]Migration, error) {
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

func (m *MigrationManager) scan() ([]Migration, []Migration, error) {
	log.Println("Scanning migrations...")
	localMigrations, err := m.scanMigrationsDirectory()
	if err != nil {
		return nil, nil, err
	}
	dbMigrations, err := m.scanMigrationsTable()
	if err != nil {
		return nil, nil, err
	}
	log.Println("Local migrations:", localMigrations)
	log.Println("Database migrations:", dbMigrations)
	return localMigrations, dbMigrations, nil
}

func (m *MigrationManager) scanMigrationsDirectory() ([]Migration, error) {
	files, err := os.ReadDir(MIGRATION_DIRECTORY)
	if err != nil {
		return nil, err
	}
	var migrations []Migration
	id := uint(1)
	for _, file := range files {
		if file.IsDir() && strings.Contains(file.Name(), MIGRATION_IDENTIFIER) {
			name, err := parseMigrationName(file.Name())
			if err != nil {
				return nil, err
			}
			migrations = append(migrations, Migration{id, name, time.Time{}})
			id++
		}
	}
	return migrations, nil
}

func (m *MigrationManager) scanMigrationsTable() ([]Migration, error) {
	query, err := m.DB.Query(SELECT_MIGRATIONS)
	if err != nil {
		return nil, err
	}
	var migrations []Migration
	var id uint
	var name string
	var createdAt time.Time
	for query.Next() {
		err = query.Scan(&id, &name, &createdAt)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, Migration{id, name, createdAt})
	}
	if m.LastVersion > id {
		return nil, errors.New("Current migration does not exist in database") // TODO: Implement custom error
	}
	return migrations, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Helpers

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
