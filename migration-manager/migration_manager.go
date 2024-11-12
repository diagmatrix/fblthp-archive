package mm

import (
	"database/sql"
	"github.com/diagmatrix/fblthp-archive/queries"
	"log"
)

type MigrationManager struct {
	DB         *sql.DB     // Database connection (SQL based)
	Migrations []Migration // List of migrations
	Current    int         // Current migration version id
}

func NewMigrationManager(db *sql.DB) *MigrationManager {
	// Check if the database connection is valid
	err := db.Ping()
	if err != nil {
		log.Fatalf("Error: %v", err) // TODO: Implement custom error
	}

	return &MigrationManager{DB: db}
}

// ---------------------------------------------------------------------------------------------------------------------
// Public methods

// ---------------------------------------------------------------------------------------------------------------------
// Private methods

func (mm *MigrationManager) scanDBMigrations() error {
	// Scan the database for migrations
	res, err := mm.DB.Query(queries.GET_MIGRATIONS)
	if err != nil {
		return err
	}

	var migrations []Migration
	var createdAt, lastRun interface{}
	for res.Next() {
		var m Migration
		err = res.Scan(&m.ID, &m.Name, &createdAt, &lastRun, &m.Executed)
		if err != nil {
			return err
		}
		migrations = append(migrations, m)
	}
	mm.Migrations = migrations

	return nil
}

func (mm *MigrationManager) scanLocalMigrations() error {
	// Scan the database for migrations
	err := mm.scanDBMigrations()
	if err != nil {
		return err
	}

	// Get the current migration version id
	mm.Current = mm.Migrations[len(mm.Migrations)-1].ID

	return nil
}
