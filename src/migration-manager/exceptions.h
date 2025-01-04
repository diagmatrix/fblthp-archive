/**
 * Custom exceptions for the project
 * @author diagmatrix
 * @date 2024
 * @version 1.0
 */

#ifndef EXCEPTIONS_H
#define EXCEPTIONS_H
#include <exception>
#include <string>

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
    explicit migration_execution_error(const std::string& message) {
        this->msg = "Error: Migration execution failed - " + message;
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
    explicit parse_migration_error(const std::string& message) {
        this->msg = "Error: Migration parsing failed - " + message;
    }
    const char* what() const noexcept override {
        return this->msg.c_str();
    }
};

#endif //EXCEPTIONS_H
