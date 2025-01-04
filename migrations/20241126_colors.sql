-- MIGRATION UP START
------------------------------------------------------------------------------------------------------------------------
CREATE TABLE colors (
    name TEXT NOT NULL,
    white BOOLEAN NOT NULL,
    blue BOOLEAN NOT NULL,
    black BOOLEAN NOT NULL,
    red BOOLEAN NOT NULL,
    green BOOLEAN NOT NULL,
    _id INTEGER PRIMARY KEY AUTOINCREMENT
);
INSERT INTO colors(name, white, blue, black, red, green) VALUES
('Colorless', FALSE, FALSE, FALSE, FALSE, FALSE),
('White', TRUE, FALSE, FALSE, FALSE, FALSE),
('Blue', FALSE, TRUE, FALSE, FALSE, FALSE),
('Black', FALSE, FALSE, TRUE, FALSE, FALSE),
('Red', FALSE, FALSE, FALSE, TRUE, FALSE),
('Green', FALSE, FALSE, FALSE, FALSE, TRUE),
('Azorius', TRUE, TRUE, FALSE, FALSE, FALSE),
('Dimir', FALSE, TRUE, TRUE, FALSE, FALSE),
('Rakdos', FALSE, FALSE, TRUE, TRUE, FALSE),
('Gruul', FALSE, FALSE, FALSE, TRUE, TRUE),
('Selesnya', TRUE, FALSE, FALSE, FALSE, TRUE),
('Orzhov', TRUE, FALSE, TRUE, FALSE, FALSE),
('Izzet', FALSE, TRUE, FALSE, TRUE, FALSE),
('Golgari', TRUE, FALSE, TRUE, FALSE, TRUE),
('Boros', TRUE, FALSE, FALSE, TRUE, FALSE),
('Simic', TRUE, TRUE, FALSE, FALSE, TRUE),
('Esper', TRUE, TRUE, TRUE, FALSE, FALSE),
('Grixis', FALSE, TRUE, TRUE, TRUE, FALSE),
('Jund', FALSE, FALSE, TRUE, TRUE, TRUE),
('Naya', TRUE, FALSE, FALSE, TRUE, TRUE),
('Bant', TRUE, TRUE, FALSE, FALSE, TRUE),
('Abzan', TRUE, FALSE, TRUE, FALSE, TRUE),
('Jeskai', TRUE, TRUE, FALSE, TRUE, FALSE),
('Sultai', FALSE, TRUE, TRUE, FALSE, TRUE),
('Mardu', TRUE, FALSE, TRUE, TRUE, FALSE),
('Temur', FALSE, TRUE, FALSE, TRUE, TRUE),
('Yore-Tiller', TRUE, TRUE, TRUE, TRUE, FALSE),
('Glint-Eye', FALSE, TRUE, TRUE, TRUE, TRUE),
('Dune-Brood', TRUE, TRUE, TRUE, FALSE, TRUE),
('Ink-Treader', TRUE, TRUE, FALSE, TRUE, TRUE),
('Witch-Maw', TRUE, TRUE, TRUE, FALSE, TRUE),
('WUBRG', TRUE, TRUE, TRUE, TRUE, TRUE);
------------------------------------------------------------------------------------------------------------------------
-- MIGRATION UP END

-- MIGRATION DOWN START
------------------------------------------------------------------------------------------------------------------------
DROP TABLE colors;
------------------------------------------------------------------------------------------------------------------------
-- MIGRATION DOWN END