package queries

// ---------------------------------------------------------------------------------------------------------------------
// Database tables and schemas

const (
	MIGRATIONS_SCHEMA string = "doorkeeper"
	MIGRATIONS_TABLE  string = MIGRATIONS_SCHEMA + ".migrations"
)

// ---------------------------------------------------------------------------------------------------------------------
// Constant Queries

// Select queries
const (
	GET_MIGRATIONS           string = "SELECT * FROM " + MIGRATIONS_TABLE + " ORDER BY id"
	GET_CURRENT_MIGRATION_ID string = "SELECT id FROM " + MIGRATIONS_TABLE + " WHERE executed = True ORDER BY id DESC LIMIT 1"
)

// Drop queries
const DROP_MIGRATIONS = "TRUNCATE TABLE " + MIGRATIONS_TABLE
