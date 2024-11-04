package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/diagmatrix/fblthp-archive/exceptions"
)

type MigrationManager struct {
	DB          *sql.DB
	LastVersion uint
}

// ---------------------------------------------------------------------------------------------------------------------
// Public methods

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
		direction := MigrationDirection(argument)
		if direction == "" {
			return errors.New("No operation provided") // TODO: Implement custom error
		}
		return m.upgrade(direction)
	case Downgrade:
		direction := MigrationDirection(argument)
		if direction == "" {
			return errors.New("No operation provided") // TODO: Implement custom error
		}
		return m.downgrade(direction)
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

// ---------------------------------------------------------------------------------------------------------------------
// Private methods

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
		newMigrations, err := getNewMigrations(localMigrations, dbMigrations)
		if err != nil {
			return err
		}
		if len(newMigrations) == 0 {
			return errors.New("No new migrations found") // TODO: Implement custom error
		}
		for _, migration := range newMigrations {
			log.Println("Upgrading to " + migration.Name + "...")
			upgradeSQL, err := migration.GetOperation(Upgrade)
			if err != nil {
				return err
			}
			// Add migration to migrations table SQL statement
			upgradeSQL += "\n" + insertMigration(migration.Name) + ";"
			// Execute migration
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

func (m *MigrationManager) downgrade(direction MigrationDirection) error {
	log.Println("Downgrading to " + direction + "...")
	switch direction {
	case Head:
		return errors.New("Cannot downgrade to head") // TODO: Implement custom error
	case Tail:
		dbMigrations, err := m.scanMigrationsTable()
		if err != nil {
			return err
		}
		if len(dbMigrations) == 0 {
			return errors.New("No migrations to downgrade") // TODO: Implement custom error
		}
		dbMigrations = dbMigrations[:len(dbMigrations)-1]
		for _, migration := range dbMigrations {
			log.Println("Downgrading from " + migration.Name + "...")
			downgradeSQL, err := migration.GetOperation(Downgrade)
			if err != nil {
				return err
			}
			// Add delete from migrations table SQL statement
			downgradeSQL += "\n" + deleteMigration(migration.Name) + ";"
			// Execute migration
			_, err = m.DB.Exec(downgradeSQL)
			if err != nil {
				return err
			}
		}
		m.LastVersion = 0
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
	newMigrations, err := getNewMigrations(localMigrations, dbMigrations)
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
