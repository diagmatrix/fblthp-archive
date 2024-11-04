-- Create table for colors
CREATE TABLE archive.colors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    white BOOLEAN NOT NULL,
    blue BOOLEAN NOT NULL,
    black BOOLEAN NOT NULL,
    red BOOLEAN NOT NULL,
    green BOOLEAN NOT NULL
);

-- Insert colors
INSERT INTO archive.colors (name, white, blue, black, red, green) VALUES
-- Colorless
('Colorless', FALSE, FALSE, FALSE, FALSE, FALSE),
-- Mono colors
('White', TRUE, FALSE, FALSE, FALSE, FALSE),
('Blue', FALSE, TRUE, FALSE, FALSE, FALSE),
('Black', FALSE, FALSE, TRUE, FALSE, FALSE),
('Red', FALSE, FALSE, FALSE, TRUE, FALSE),
('Green', FALSE, FALSE, FALSE, FALSE, TRUE),
-- Guilds
('Azorius', TRUE, TRUE, FALSE, FALSE, FALSE),
('Dimir', FALSE, TRUE, TRUE, FALSE, FALSE),
('Rakdos', FALSE, FALSE, TRUE, TRUE, FALSE),
('Gruul', TRUE, FALSE, FALSE, TRUE, TRUE),
('Selesnya', TRUE, FALSE, TRUE, FALSE, TRUE),
('Orzhov', TRUE, FALSE, TRUE, FALSE, FALSE),
('Izzet', FALSE, TRUE, FALSE, TRUE, TRUE),
('Golgari', TRUE, FALSE, TRUE, FALSE, TRUE),
('Boros', TRUE, TRUE, FALSE, TRUE, FALSE),
('Simic', TRUE, TRUE, FALSE, FALSE, TRUE),
-- Shards
('Bant', TRUE, TRUE, TRUE, FALSE, TRUE),
('Esper', TRUE, TRUE, TRUE, TRUE, FALSE),
('Grixis', FALSE, TRUE, TRUE, TRUE, FALSE),
('Jund', TRUE, FALSE, TRUE, TRUE, TRUE),
('Naya', TRUE, TRUE, FALSE, TRUE, TRUE),
-- Wedges
('Abzan', TRUE, TRUE, TRUE, FALSE, TRUE),
('Jeskai', TRUE, TRUE, FALSE, TRUE, TRUE),
('Sultai', TRUE, FALSE, TRUE, TRUE, TRUE),
('Mardu', TRUE, FALSE, TRUE, TRUE, FALSE),
('Temur', FALSE, TRUE, TRUE, TRUE, TRUE),
-- 4 colors
('Yore-Tiller', TRUE, TRUE, TRUE, TRUE, FALSE),
('Glint-Eye', FALSE, TRUE, TRUE, TRUE, TRUE),
('Dune-Brood', TRUE, FALSE, TRUE, TRUE, TRUE),
('Ink-Treader', TRUE, TRUE, FALSE, TRUE, TRUE),
('Witch-Maw', TRUE, TRUE, TRUE, FALSE, TRUE),
-- 5 colors
('WUBRG', TRUE, TRUE, TRUE, TRUE, TRUE);