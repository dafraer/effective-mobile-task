CREATE TABLE IF NOT EXISTS people (
		id SERIAL PRIMARY KEY,
		name TEXT,
		surname TEXT,
		patronymic TEXT,
        age INT,
        geneder TEXT,
        nationality TEXT
);