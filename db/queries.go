package db

// ---------------------------------------------------------------------------------------------------------------------
// Constant queries
const (
	SELECT_MIGRATIONS        string = "SELECT * FROM doorkeeper.migrations"
	SELECT_CURRENT_MIGRATION string = "SELECT * FROM doorkeeper.migrations WHERE id = (SELECT migration_id FROM doorkeeper.current_migration LIMIT 1)"
)

// ---------------------------------------------------------------------------------------------------------------------
// Variable queries

func insertMigration(name string) string {
	return "INSERT INTO doorkeeper.migrations (name) VALUES ('" + name + "')"
}

func deleteMigration(name string) string {
	return "DELETE FROM doorkeeper.migrations WHERE name = '" + name + "'"
}
