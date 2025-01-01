-- MIGRATION UP START
CREATE TABLE IF NOT EXISTS raw_collection (
    name VARCHAR,
    set_code VARCHAR,
    number VARCHAR,
    rarity VARCHAR(1),
    quantity INTEGER,
    added TIMESTAMP,
    last_modified TIMESTAMP,
    foil BOOLEAN,
    colors json,
    color_id json,
    _id INTEGER PRIMARY KEY AUTOINCREMENT,
    _created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- MIGRATION UP END

-- MIGRATION DOWN START
DROP TABLE IF EXISTS raw_collection;
-- MIGRATION DOWN END
