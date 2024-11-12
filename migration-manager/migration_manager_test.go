package mm

import (
	"database/sql"
	"github.com/diagmatrix/fblthp-archive/queries"
	_ "github.com/jackc/pgx/v5/stdlib" // Postgres SQL driver
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

// ---------------------------------------------------------------------------------------------------------------------
// Tests

// Test to check if the scanning of the database for migrations works
func TestMigrationManagerScanDBOk(t *testing.T) {
	err := insertMigrations(MIGRATIONS)
	assert.Nil(t, err)

	manager := NewMigrationManager(db)

	err = manager.scanDBMigrations()
	assert.Nil(t, err)
	assert.Equal(t, N_MIGRATIONS, len(manager.Migrations))

	t.Cleanup(func() {
		err = dropMigrations()
		assert.Nil(t, err)
	})
}

// Test to check if the scanning of migrations when some migrations are added works
func TestMigrationManagerScanNew(t *testing.T) {
	t.Error("Not implemented")
}

// Test to check if the scanning of migrations when some migrations correctly returns an error
func TestMigrationManagerScanMissing(t *testing.T) {
	t.Error("Not implemented")
}

// Test to check if the scanning of migrations when some migrations are unordered correctly returns an error
func TestMigrationManagerScanUnordered(t *testing.T) {
	t.Error("Not implemented")
}

// ---------------------------------------------------------------------------------------------------------------------
// Fixtures

var db *sql.DB

const MIGRATIONS string = "../testdata/migrations.sql"
const N_MIGRATIONS int = 3

func setup() {
	// Connect to the database
	const POSTGRES_CONNECTION string = "host=localhost port=5432 user=postgres password=postgres dbname=fblthp-test sslmode=disable" // TODO: Get from env
	var err error
	db, err = sql.Open("pgx", POSTGRES_CONNECTION)
	if err != nil {
		log.Fatalf("Error in test setup: %v", err) // TODO: Implement custom error
	}
}

func teardown() {
	// Delete all the data
	err := dropMigrations()
	if err != nil {
		log.Fatalf("Error in test teardown: %v", err) // TODO: Implement custom error
	}

	// Close the database connection
	err = db.Close()
	if err != nil {
		log.Fatalf("Error in test teardown: %v", err) // TODO: Implement custom error
	}
}

func dropMigrations() error {
	_, err := db.Exec(queries.DROP_MIGRATIONS)
	return err
}

func insertMigrations(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(file))
	if err != nil {
		return err
	}
	return nil
}
