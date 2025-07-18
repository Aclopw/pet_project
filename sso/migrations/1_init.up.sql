CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	email VARCHAR(255) UNIQUE,
	tg_user_id VARCHAR(255) UNIQUE,
	is_active BOOLEAN NOT NULL DEFAULT FALSE,
	activation_link VARCHAR(255) NOT NULL,
	password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL,
	token VARCHAR(255) NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
   