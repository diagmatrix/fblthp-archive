#include <iostream>
#include <sqlite3.h>
#include <regex>
#include <string.h>

#include "migrations.h"

const char* MIGRATIONS_DB = "archive.db"; // TODO: Add from env
const char* MIGRATIONS_FOLDER = "migrations"; // TODO: Add from env
const char* HELP =
    "doorkeeper <option> [argument]\n"
    "Options:\n"
    "  -h, --help\t\tShow this help message\n"
    "  -s, --status\t\tShow the migrations status\n"
    "  -u, --upgrade\t\tMigrate the database upwards\n"
    "  -d, --downgrade\tMigrate the database downwards\n"
    "  -g, --generate\tGenerate a new migration\n"
    "Arguments for upgrade and downgrade:\n"
    "  head\t\t\tMigrate to the latest version\n"
    "  base\t\t\tMigrate to the initial version\n"
    "  <number>\t\tMigrate up or down <number> of versions\n"
    "Arguments for generate:\n"
    "  <name>\t\tName to give to migration\n"; ///< Help message

enum OPTIONS {STATUS, UPGRADE, DOWNGRADE, GENERATE}; ///< Migration manager options
bool str_eq(const char* str1, const char* str2) {return strcmp(str1, str2) == 0;} ///< Compare if two strings are equal

int main(int argc, const char* argv[]) {
    // Parse command line arguments
    if (argc < 2) {
        std::cout << HELP;
        return EXIT_FAILURE;
    }
    if (argc > 3) {
        std::cout << "Error: Too many arguments\n";
        return EXIT_FAILURE;
    }
    int option;
    if (str_eq(argv[1], "-h") || str_eq(argv[1], "--help")) {
        std::cout << HELP;
        return EXIT_SUCCESS;
    }
    if (str_eq(argv[1], "-u") || str_eq(argv[1], "--upgrade")) {
        option = UPGRADE;
    } else if (str_eq(argv[1], "-d") || str_eq(argv[1], "--downgrade")) {
        option = DOWNGRADE;
    } else if (str_eq(argv[1], "-s") || str_eq(argv[1], "--status")) {
        option = STATUS;
    } else if (str_eq(argv[1], "-g") || str_eq(argv[1], "--generate")) {
        option = GENERATE;
    } else {
        std::cout << "Error: Invalid option\n";
        return EXIT_FAILURE;
    }
    std::regex digit("\\d+");
    std::string argument;
    if (option != STATUS) {
        if (argc < 3) {
            std::cout << "Error: Missing argument\n";
            return EXIT_FAILURE;
        }
        if (str_eq(argv[2], "head") && option == UPGRADE) {
            argument = migration_constants::HEAD;
        } else if (str_eq(argv[2], "base") && option == DOWNGRADE) {
            argument = migration_constants::BASE;
        } else if (std::regex_match(argv[2], digit) && (option == UPGRADE || option == DOWNGRADE)) {
            argument = argv[2];
        } else if (option == GENERATE) {
            argument = argv[2];
        } else {
            std::cout << "Error: Invalid argument\n";
            return EXIT_FAILURE;
        }
    }

    // Open the database
    sqlite3* DB;
    if (int error = sqlite3_open(MIGRATIONS_DB, &DB); error != SQLITE_OK) {
        std::cout << sqlite3_errmsg(DB) << "\n";
        return EXIT_FAILURE;
    }

    // Initialize the migration table and manager
    if (!init_migration_table(DB)) {
        std::cout << "Error: Unable to initialize migration table: " << sqlite3_errmsg(DB) << std::endl;
        sqlite3_close(DB);
        return EXIT_FAILURE;
    }
    manager migrations_manager = create_migration_manager(MIGRATIONS_FOLDER, DB);

    // Perform the requested operation
    switch (option) {
        case UPGRADE:
            execute_migration(migrations_manager, migration_constants::UPGRADE, argument);
            break;
        case DOWNGRADE:
            execute_migration(migrations_manager, migration_constants::DOWNGRADE, argument);
            break;
        case STATUS:
            std::cout << print_migrations(migrations_manager.migrations);
            break;
        case GENERATE:
            if (!generate_migration(migrations_manager, argument)) {
                std::cout << "Error: Unable to generate migration\n";
                sqlite3_close(DB);
                return EXIT_FAILURE;
            }
            std::cout << "Migration generated successfully\n";
            break;
        default:
            std::cout << "Error: This should be unreachable\n";
            sqlite3_close(DB);
            return EXIT_FAILURE;
    }

    sqlite3_close(DB);
    return EXIT_SUCCESS;
}
