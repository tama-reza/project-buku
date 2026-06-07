CREATE TABLE IF NOT EXISTS kategori (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(255) NOT NULL,
		modified_at TIMESTAMP,
		modified_by VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(255) NOT NULL,
		modified_at TIMESTAMP,
		modified_by VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS buku (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		image_url VARCHAR(255) NOT NULL,
		release_year INTEGER NOT NULL,
		price INTEGER NOT NULL,
		total_page INTEGER NOT NULL,
		thickness VARCHAR(255) NOT NULL,
		category_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(255) NOT NULL,
		modified_at TIMESTAMP,
		modified_by VARCHAR(255),
		CONSTRAINT fk_buku_kategori FOREIGN KEY (category_id) REFERENCES kategori(id) ON DELETE CASCADE
);

INSERT INTO users (username, password, created_by)
VALUES ('admin', 'password_rahasia', 'system')
ON CONFLICT (username) DO NOTHING;
