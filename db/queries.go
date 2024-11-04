package db

const (
	// Migrations
	SELECT_MIGRATIONS        string = "SELECT * FROM doorkeeper.migrations"
	SELECT_CURRENT_MIGRATION string = "SELECT * FROM doorkeeper.migrations WHERE id = (SELECT migration_id FROM doorkeeper.current_migration LIMIT 1)"
)
