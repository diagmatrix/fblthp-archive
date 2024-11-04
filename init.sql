-- Creates the skeleton of the database

-- Schema and tables for core functionality
CREATE SCHEMA IF NOT EXISTS doorkeeper;
CREATE TABLE IF NOT EXISTS doorkeeper.migrations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS doorkeeper.current_migration (
    id SERIAL PRIMARY KEY,
    migration_id INT REFERENCES doorkeeper.migrations(id),
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Schema for the cards
CREATE SCHEMA IF NOT EXISTS archive;
