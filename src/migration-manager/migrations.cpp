#include <ranges>
#include <filesystem>
#include <fstream>
#include <iostream>

#include "migrations.h"
#include "exceptions.h"

namespace fs = std::filesystem;
using namespace migration_constants;

// Migration functions
// ---------------------------------------------------------------------------------------------------------------------
std::vector<migration> scan_migrations(const std::string& path, sqlite3* DB) {
    std::vector<migration> migrations;
    str_map local_migrations = scan_local_migrations(path);
    str_map db_migrations = scan_db_migrations(DB);

    size_t found_migrations = 0;
    for (const auto& local_key : local_migrations | std::views::keys) {
        std::string name = fs::path(local_key).stem().string();
        std::string exec_time;
        const auto [up_stmt, down_stmt] = parse_migration(local_key);
        const auto& db_migration = db_migrations.find(fs::path(local_key).stem().string());
        if (db_migration != db_migrations.end()) {
            found_migrations++;
            name = db_migration->first;
            exec_time = db_migration->second;
        }
        migrations.push_back(migration(name, exec_time,up_stmt, down_stmt));
    }

    if (found_migrations != db_migrations.size()) {
        throw inconsistent_migrations_error();
    }

    return migrations;
}

int find_last_executed(const std::vector<migration>& migrations) {
    for (int i = static_cast<int>(migrations.size()) - 1; i >= 0; i--) {
        if (!migrations[i].exec_time.empty()) {
            return i;
        }
    }
    return -1;
}

std::string print_migrations(const std::vector<migration>& migrations) {
    // Center string utility
    auto center_string = [](const std::string& str, size_t width) {
        if (str.size() >= width) {
            return str; // No need to pad if string is already larger
        }

        size_t padding = width - str.size();
        size_t padLeft = padding / 2; // Left padding
        size_t padRight = padding - padLeft; // Right padding

        return std::string(padLeft, ' ') + str + std::string(padRight, ' ');
    };

    // Get the maximum size of the migration name and execution time strings
    const std::string migration = "Migration";
    const std::string executed_at = "Executed At";
    int mig_name_size = migration.size();
    int exec_time_size = executed_at.size();
    for (const auto& mig : migrations) {
        mig_name_size = std::max(static_cast<int>(mig.name.size()), mig_name_size);
        exec_time_size = std::max(static_cast<int>(mig.exec_time.size()), exec_time_size);
    }

    // Create the table
    const std::string separator = "+-" + std::string(mig_name_size, '-') + "-+-" + std::string(exec_time_size, '-') + "-+";
    const std::string header = "| " + center_string(migration, mig_name_size) + " | " + center_string(executed_at, exec_time_size) + " |";
    std::string table = separator + "\n" + header + "\n" + separator + "\n";
    for (const auto& mig : migrations) {
        table += "| " + center_string(mig.name, mig_name_size) + " | " + center_string(mig.exec_time, exec_time_size) + " |\n";
        table += separator + "\n";
    }

    return table;
}

// Migration manager functions
// ---------------------------------------------------------------------------------------------------------------------
manager create_migration_manager(const std::string& path, sqlite3* DB) {
    manager manager;
    manager.path = path;
    manager.DB = DB;
    manager.migrations = scan_migrations(path, DB);
    manager.last_executed_idx = find_last_executed(manager.migrations);
    return manager;
}

void execute_migration(manager& manager, const std::string& operation, const std::string& arg) {
    int target_idx;
    if (arg == HEAD) {
        target_idx = static_cast<int>(manager.migrations.size());
    } else if (arg == BASE) {
        target_idx = INT32_MAX; // This "should" be large enough
    } else {
        target_idx = std::stoi(arg);
    }

    if (operation == UPGRADE) {
        unsigned int idx = std::min(manager.last_executed_idx + target_idx + 1, static_cast<int>(manager.migrations.size()));
        for (size_t i = std::min(manager.last_executed_idx + 1, static_cast<int>(manager.migrations.size())); i < idx; i++) {;
            sqlite3_exec(manager.DB, manager.migrations[i].up_stmt.c_str(), nullptr, nullptr, nullptr);
            if (sqlite3_errcode(manager.DB) != SQLITE_OK) {
                const std::string err_msg = "Upgrading to " + manager.migrations[i].name + "(" + sqlite3_errmsg(manager.DB) + ")";
                throw migration_execution_error(err_msg);
            }
            std::cout << "Upgraded to migration: " << manager.migrations[i].name << std::endl;
            const std::string sql = "INSERT INTO migrations (name) VALUES ('" + manager.migrations[i].name + "');";
            sqlite3_exec(manager.DB, sql.c_str(), nullptr, nullptr, nullptr);
            if (sqlite3_errcode(manager.DB) != SQLITE_OK) {
                const std::string err_msg = "Adding to migrations (" + std::string(sqlite3_errmsg(manager.DB)) + ")";
                throw migration_execution_error(err_msg);
            }
            manager.last_executed_idx++;
        }
    } else if (operation == DOWNGRADE) {
        int idx = std::max(manager.last_executed_idx - target_idx + 1, 0);
        for (int i = manager.last_executed_idx; i >= idx; i--) {
            sqlite3_exec(manager.DB, manager.migrations[i].down_stmt.c_str(), nullptr, nullptr, nullptr);
            if (sqlite3_errcode(manager.DB) != SQLITE_OK) {
                std::string err_msg = "Downgrading from " + manager.migrations[i].name + "(" + sqlite3_errmsg(manager.DB) + ")";
                throw migration_execution_error(err_msg);
            }
            std::cout << "Downgraded migration " << manager.migrations[i].name << std::endl;
            const std::string sql = "DELETE FROM migrations WHERE name = '" + manager.migrations[i].name + "';";
            sqlite3_exec(manager.DB, sql.c_str(), nullptr, nullptr, nullptr);
            if (sqlite3_errcode(manager.DB) != SQLITE_OK) {
                const std::string err_msg = "Removing from migrations (" + std::string(sqlite3_errmsg(manager.DB)) + ")";
                throw migration_execution_error(err_msg);
            }
            manager.last_executed_idx--;
        }
    }
}

bool generate_migration(const manager& manager, const std::string& name) {
    // Get date
    std::time_t t = std::time(nullptr);
    std::tm* now = std::localtime(&t);
    char buffer[9];
    std::strftime(buffer, sizeof(buffer), "%Y%m%d", now);

    // Create path and file contents
    const std::string path = manager.path + "/" + buffer + "_" + name + ".sql";
    const std::string sql =
        "-- MIGRATION UP START\n"
        " << Add SQL statements for upgrade\n"
        "-- MIGRATION UP END\n\n"
        "-- MIGRATION DOWN START\n"
        " << Add SQL statements for downgrade\n"
        "-- MIGRATION DOWN END\n";

    // Write to file
    std::ofstream file(path);
    if (file.is_open()) {
        file << sql;
        file.close();
        return true;
    }

    return false;
}

// Helper functions
// ---------------------------------------------------------------------------------------------------------------------
bool init_migration_table(sqlite3* DB) {
    std::string sql = "CREATE TABLE IF NOT EXISTS migrations ("
                      "id INTEGER PRIMARY KEY AUTOINCREMENT,"
                      "name TEXT NOT NULL,"
                      "executed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"
                      ");";

    return sqlite3_exec(DB, sql.c_str(), nullptr, nullptr, nullptr) == SQLITE_OK;
}

str_map scan_local_migrations(const std::string& path) {
    str_map migrations;
    if (fs::exists(path) && fs::is_directory(path)) {
        for (const auto& entry : fs::directory_iterator(path)) {
            if (entry.is_regular_file() && entry.path().extension() == ".sql") {
                migrations.insert({entry.path().string(), ""});
            }
        }
    }

    return migrations;
}

str_map scan_db_migrations(sqlite3* DB) {
    str_map migrations;
    sqlite3_stmt* stmt;

    if (sqlite3_prepare_v2(DB, SELECT_ALL_MIGRATIONS, -1, &stmt, nullptr) == SQLITE_OK) {
        while (sqlite3_step(stmt) == SQLITE_ROW) {
            std::string name = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 1));
            std::string date = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 2));
            migrations.insert({name, date});
        }
    }

    return migrations;
}

str_pair parse_migration(const std::string& path) {
    str_pair migration_stmt;

    std::ifstream file(path);
    if (!file.is_open()) {
        const std::string err_msg = "Unable to open file " + path;
        throw parse_migration_error(err_msg);
    }

    bool is_upgrade = false;
    bool is_downgrade = false;
    std::string line;
    while(std::getline(file, line)) {
        if (line == UP_START_TAG) {
            if (is_downgrade) {
                file.close();
                const std::string err_msg = "In " + path + " (Start up statement before end down)";
                throw parse_migration_error(err_msg);
            }
            is_upgrade = true;
        } else if (line == DOWN_START_TAG) {
            if (is_upgrade) {
                file.close();
                const std::string err_msg = "In " + path + " (Start down statement before end up)";
                throw parse_migration_error(err_msg);
            }
            is_downgrade = true;
        } else if (line == UP_END_TAG) {
            if (!is_upgrade) {
                file.close();
                const std::string err_msg = "In " + path + " (End up statement before start)";
                throw parse_migration_error(err_msg);
            }
            is_upgrade = false;
        } else if (line == DOWN_END_TAG) {
            if (!is_downgrade) {
                file.close();
                const std::string err_msg = "In " + path + " (End down statement before start)";
                throw parse_migration_error(err_msg);
            }
            is_downgrade = false;
        } else {
            if (is_upgrade) {
                migration_stmt.first += line + "\n";
            } else if (is_downgrade) {
                migration_stmt.second += line + "\n";
            }
        }
    }

    file.close();
    return migration_stmt;
}

std::string import_sql(const std::string& path) {
    std::string sql;
    if (fs::exists(path) && fs::is_regular_file(path)) {
        std::ifstream file(path);
        if (file.is_open()) {
            std::string line;
            while (std::getline(file, line)) {
                sql += line + "\n";
            }
        }
    }

    return sql;
}

