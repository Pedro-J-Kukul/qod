-- Up migration for quotes table
CREATE TABLE IF NOT EXISTS quotes (
	id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
	quote TEXT NOT NULL,
	author TEXT NOT NULL
);