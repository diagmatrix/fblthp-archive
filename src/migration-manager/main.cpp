/**
 * Migration manager for handling database migrations
 * @author diagmatrix
 * @date 2024
 * @version 1.1
 */

#include <iostream>
#include <sqlite3.h>
#include <regex>
#include <string.h>
#include <map>
#include <cstdlib>

#include "migrations.h"
#include "env.h"

const std::map<std::string, std::string> DEFAULT_ENV = {
    {"MIGRATIONS_DB", "default.db"},
    {"MIGRATIONS_FOLDER", "migrations"}    
}; ///< Default environment variables
const char* HELP_MESSAGE =
    "doorkeeper <option> [argument]\n"
    "Options:\n"
    "  -h, --help\t\tShow this help message\n"
    "  -e, --environment\t\tUse custom environment\n"
    "  -s, --status\t\tShow the migrations status\n"
    "  -u, --upgrade\t\tMigrate the database upwards\n"
    "  -d, --downgrade\tMigrate the database downwards\n"
    "  -g, --generate\tGenerate a new migration\n"
    "Arguments for environment:\n"
    "  <file>\t\tPath to environment file\n"
    "Arguments for upgrade and downgrade:\n"
    "  head\t\t\tMigrate to the latest version\n"
    "  base\t\t\tMigrate to the initial version\n"
    "  <number>\t\tMigrate up or down <number> of versions\n"
    "Arguments for generate:\n"
    "  <name>\t\tName to give to migration\n"; ///< Help message
enum OPTIONS {HELP, ENVIRONMENT, STATUS, UPGRADE, DOWNGRADE, GENERATE}; ///< Migration manager options

bool str_eq(const char* str1, const char* str2) {return strcmp(str1, str2) == 0;} ///< Compare if two strings are equal
int get_option(const char* argument); ///< Check if an argument is an option and returns the option or -1 if false

int main(int argc, const char* argv[]) {
    // Parse command line arguments
    if (argc < 2) {
        std::cout << HELP_MESSAGE;
        return EXIT_FAILURE;
    }
    
    std::regex digit("\\d+");
    int option = -1;
    std::map<int, std::string> commands;
    for (int i = 1; i < argc; i++) {
        option = get_option(argv[i]);
        if (option == HELP || option == STATUS) {
            try {
                commands.insert({option, ""});
            } catch(const std::exception& e) {
                std::cout << "Error: Duplicate option" << std::endl;
                return EXIT_FAILURE;
            }
        } else if (option == UPGRADE || option == DOWNGRADE) {
            if (++i >= argc) { // Next argument
                std::cout << "Error: Missing argument" << std::endl;
                return EXIT_FAILURE;
            }
            std::string argument;
            if (option == UPGRADE && str_eq(argv[i], "head")) {
                argument = migration_constants::HEAD;
            } else if (option == DOWNGRADE && str_eq(argv[i], "base")) {
                argument = migration_constants::BASE;
            } else if (std::regex_match(argv[i], digit)) {
                argument = argv[i];
            } else {
                std::cout << "Error: Invalid argument: " << argv[i] << std::endl;
                return EXIT_FAILURE;
            }
            try {
                commands.insert({option, argument});
            } catch(const std::exception& e) {
                std::cout << "Error: Duplicate option" << std::endl;
                return EXIT_FAILURE;
            }
        } else if (option == GENERATE || option == ENVIRONMENT) {
            if (++i >= argc) { // Next argument
                std::cout << "Error: Missing argument" << std::endl;
                return EXIT_FAILURE;
            }
            try {
                commands.insert({option, argv[i]});
            } catch(const std::exception& e) {
                std::cout << "Error: Duplicate option" << std::endl;
                return EXIT_FAILURE;
            }
        }
    }

    // Execute help if present and exit
    if (commands.find(HELP) != commands.end()) {
        std::cout << HELP_MESSAGE;
        return EXIT_SUCCESS;
    }

    // Execute environment
    std::string env_file;
    if (const auto env = commands.find(ENVIRONMENT); env != commands.end()) {
        env_file = env->second;
        commands.erase(env);
    }
    load_env(env_file, DEFAULT_ENV);

    // Check only 1 option and get option
    if (commands.size() > 1) {
        std::cout << "Error: Too many options" << std::endl;
        return EXIT_FAILURE;
    } else if (commands.empty()) {
        std::cout << "Error: No option provided" << std::endl;
        return EXIT_FAILURE;
    }
    const auto command = commands.begin();
    option = command->first;
    std::string argument = command->second;

    // Open the database
    sqlite3* DB;
    if (int error = sqlite3_open(std::getenv("MIGRATIONS_DB"), &DB); error != SQLITE_OK) {
        std::cout << sqlite3_errmsg(DB) << "\n";
        return EXIT_FAILURE;
    }

    // Initialize the migration table and manager
    if (!init_migration_table(DB)) {
        std::cout << "Error: Unable to initialize migration table: " << sqlite3_errmsg(DB) << std::endl;
        sqlite3_close(DB);
        return EXIT_FAILURE;
    }
    manager migrations_manager = create_migration_manager(std::getenv("MIGRATIONS_FOLDER"), DB);

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

int get_option(const char* argument) {
    if (str_eq(argument, "-h") || str_eq(argument, "--help")) {
        return HELP;
    }
    if (str_eq(argument, "-e") || str_eq(argument, "--environment")) {
        return ENVIRONMENT;
    }
    if (str_eq(argument, "-u") || str_eq(argument, "--upgrade")) {
        return UPGRADE;
    }
    if (str_eq(argument, "-d") || str_eq(argument, "--downgrade")) {
        return DOWNGRADE;
    }
    if (str_eq(argument, "-s") || str_eq(argument, "--status")) {
        return STATUS;
    }
    if (str_eq(argument, "-g") || str_eq(argument, "--generate")) {
        return GENERATE;
    }
    return -1;
}       