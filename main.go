package main

import (
	"database/sql"
	"github.com/diagmatrix/fblthp-archive/db"
	_ "github.com/jackc/pgx/v5/stdlib" // Postgres SQL driver
	"log"
)

func main() {
	// TODO: Get from environment
	const POSTGRES_CONNECTION_STRING = "host=localhost port=5432 user=postgres password=postgres dbname=fblthp-archive sslmode=disable"

	// Get command line arguments
	//argv := os.Args
	//argc := len(argv)
	//if argc < 2 {
	//	log.Fatalf("No migration option provided")
	//}

	// Initialize database connection
	conn, err := sql.Open("pgx", POSTGRES_CONNECTION_STRING)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Initialize migration manager
	migrationManager := db.NewMigrationManager(conn)

	// Get operations
	err = migrationManager.Process("status", "")
	if err != nil {
		log.Fatalf("Failed to process: %v", err)
	}
	err = migrationManager.Process("upgrade", "head")
	if err != nil {
		log.Fatalf("Failed to process: %v", err)
	}
	err = migrationManager.Process("generate", "")
	if err != nil {
		log.Fatalf("Failed to process: %v", err)
	}
	err = migrationManager.Process("help", "")
	if err != nil {
		log.Fatalf("Failed to process: %v", err)
	}
}
