-- Write your migrate up statements here

CREATE TYPE role AS ENUM ('admin', 'editor', 'author', 'user', 'guest');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(128) UNIQUE NOT NULL,
    avatar_url VARCHAR(255),
    full_name VARCHAR(64),
    bio VARCHAR(2000),
    last_login_at TIMESTAMP,
    role role NOT NULL DEFAULT 'user',
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

---- create above / drop below ----

DROP INDEX idx_users_email;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS role;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
