/**
 * Migration manager header file
 * @author diagmatrix
 * @date 2024
 * @version 1.0
 */

#ifndef MIGRATION_MANAGER_H
#define MIGRATION_MANAGER_H
#include <map>
#include <sqlite3.h>
#include <string>
#include <vector>

// Aliases
// -----------------------------------------------------------------------------------------------------------------
typedef std::pair<std::string, std::string> str_pair;
typedef std::map<std::string, std::string> str_map;

// Constants
// -----------------------------------------------------------------------------------------------------------------
namespace migration_constants {
    inline const char* SELECT_ALL_MIGRATIONS = "SELECT * FROM migrations;"; ///< SQL statement to retrieve all migrations
    inline const char* UP_START_TAG = "-- MIGRATION UP START"; ///< Start line for upgrade statement
    inline const char* UP_END_TAG = "-- MIGRATION UP END"; ///< End line for upgrade statement
    inline const char* DOWN_START_TAG = "-- MIGRATION DOWN START"; ///< Start line for downgrade statement
    inline const char* DOWN_END_TAG = "-- MIGRATION DOWN END"; ///< End line for downgrade statement
    inline const char* UPGRADE = "UP"; ///< Upgrade operation string
    inline const char* DOWNGRADE = "DOWN"; ///< Downgrade operation string
    inline const char* HEAD = "HEAD"; ///< Head migration string
    inline const char* BASE = "BASE"; ///< Base migration string
}

// Types
// -----------------------------------------------------------------------------------------------------------------
/**
 * Struct to hold migration information
 */
struct migration {
    std::string name; ///< Migration name
    std::string exec_time; ///< Execution time
    std::string up_stmt; ///< Upgrade SQL statement
    std::string down_stmt; ///< Downgrade SQL statement
};

/**
 * Scans the migrations found in a directory
 * @param path Path of the directory
 * @param DB Sqlite database object
 * @return List of migration records found
 * @throw inconsistent_migrations_error if the migration numbers are inconsistent
 */
std::vector<migration> scan_migrations(const std::string& path, sqlite3* DB);

/**
 * Finds the last executed migration
 * @param migrations List of migrations
 * @return Index of the last executed migration, -1 if none found
 */
int find_last_executed(const std::vector<migration>& migrations);

/**
 * Prints the list of migrations
 * @param migrations List of migrations
 * @return String with the list of migrations pretty printed
 */
std::string print_migrations(const std::vector<migration>& migrations);

/**
 * Struct to hold migration manager information
 */
struct manager {
    std::string path; ///< Path to the migrations directory
    sqlite3* DB = nullptr; ///< Sqlite database object
    std::vector<migration> migrations; ///< List of migrations
    int last_executed_idx = -1; ///< Index of the last executed migration
};

/**
 * Creates a migration manager object
 * @param path Path of the migrations
 * @param DB Sqlite database object
 * @return Manager object
 */
manager create_migration_manager(const std::string& path, sqlite3* DB);

/**
 * Executes a migration operation
 * @param manager Migration manager object
 * @param operation Operation to execute
 * @param arg Argument for the operation (should be checked if valid before calling this function)
 * @return True if the operation was successful, false otherwise
 */
void execute_migration(manager& manager, const std::string& operation, const std::string& arg);

/**
 * Exception raised when the migration numbers are inconsistent between the database and the local directory
 */
class inconsistent_migrations_error final: public std::exception {
public:
    const char* what() const noexcept override {
        return "Error: Inconsistent migrations between DB and migrations directory";
    }
};

/**
 * Exception raised when an error occurs while executing a migration
 */
class migration_execution_error final: public std::exception {
    std::string msg;
public:
    explicit migration_execution_error(const std::string& msg) {
        this->msg = "Error: Migration execution failed - " + msg;
    }
    const char* what() const noexcept override {
        return this->msg.c_str();
    }
};

/**
 * Exception raised when an error occurs while parsing a migration
 */
class parse_migration_error final: public std::exception {
    std::string msg;
public:
    explicit parse_migration_error(const std::string& msg) {
        this->msg = "Error: Migration parsing failed - " + msg;
    }
    const char* what() const noexcept override {
        return this->msg.c_str();
    }

};

// Functions
// -----------------------------------------------------------------------------------------------------------------

/**
 * Initializes the migration table in the database
 * @param DB Sqlite database object
 * @return True if the table was created successfully, false otherwise
 */
bool init_migration_table(sqlite3* DB);

/**
 * Retrieves the list of migrations from a local directory
 * @param path Path where the migrations exist
 * @return List of migrations found
 */
str_map scan_local_migrations(const std::string& path);

/**
 * Retrieves the list of migrations from the database
 * @param DB Sqlite database object
 * @return List of migrations found
 */
str_map scan_db_migrations(sqlite3* DB);

/**
 * Parses a migration file path into its upgrade and downgrade components
 * @param path Path of the migration file
 * @return Pair with the SQL strings for the upgrade and downgrade components
 */
str_pair parse_migration(const std::string& path);

/**
 * Imports a SQL file into a string
 * @param path Path for the SQL file
 * @return A string with the contents of the SQL file
 */
std::string import_sql(const std::string& path);

#endif //MIGRATION_MANAGER_H
