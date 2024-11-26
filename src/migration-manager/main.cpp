#include <iostream>
#include <sqlite3.h>
#include "migrations.h"

int main() {
    const char* MIGRATIONS_DB = "..\\doorkeeper.db"; // TODO: Add from args / env
    sqlite3* DB;
    if (int error = sqlite3_open(MIGRATIONS_DB, &DB); error != SQLITE_OK) {
        std::cout << sqlite3_errmsg(DB) << "\n";
        return EXIT_FAILURE;
    }

    if (!init_migration_table(DB)) {
        std::cout << "Error: Unable to initialize migration table: " << sqlite3_errmsg(DB) << std::endl;
        sqlite3_close(DB);
        return EXIT_FAILURE;
    }

    sqlite3_close(DB);
    return EXIT_SUCCESS;
}