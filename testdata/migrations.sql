CREATE TABLE IF NOT EXISTS doorkeeper.migrations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_run TIMESTAMP,
    executed BOOLEAN DEFAULT FALSE,
    CONSTRAINT migrations_name_unique UNIQUE (name)
);

INSERT INTO doorkeeper.migrations (name, executed) VALUES
    ('Test migration 1', TRUE),
    ('Test migration 2', TRUE),
    ('Test migration 3', FALSE);
