
CREATE TABLE gists (
	id INTEGER NOT NULL PRIMARY KEY,
	user_id INTEGER, 
  updated_at TEXT,
	title TEXT,
	url TEXT
);

CREATE TABLE users (
	id INTEGER NOT NULL PRIMARY KEY,
	credential TEXT,
	authentication_service TEXT
);
