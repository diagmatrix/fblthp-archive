-- MIGRATION UP START
CREATE TABLE IF NOT EXISTS mtg_set (
    name VARCHAR,
    code VARCHAR,
    set_type VARCHAR,
    digital BOOLEAN,
    released_at DATE,
    card_count INTEGER,
    search_uri VARCHAR,
    icon_uri VARCHAR,
    _id INTEGER PRIMARY KEY AUTOINCREMENT,
    _created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    _modified_at TIMESTAMP
);
CREATE UNIQUE INDEX mtg_set_unique_set ON mtg_set (name, code);
-- MIGRATION UP END

-- MIGRATION DOWN START
DROP TABLE IF EXISTS mtg_set;
DROP INDEX IF EXISTS mtg_set_unique_set;
-- MIGRATION DOWN END
