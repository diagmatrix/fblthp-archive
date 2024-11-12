-- Creates the skeleton of the database

-- Schema and tables for core functionality
CREATE SCHEMA IF NOT EXISTS doorkeeper;
-- Migrations table
CREATE TABLE IF NOT EXISTS doorkeeper.migrations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_run TIMESTAMP,
    executed BOOLEAN DEFAULT FALSE,
    CONSTRAINT migrations_name_unique UNIQUE (name)
);

-- Schema for the cards
CREATE SCHEMA IF NOT EXISTS archive;
