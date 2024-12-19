-- Write your migrate up statements here

CREATE TYPE role AS ENUM ('admin', 'editor', 'author', 'user', 'guest');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(128) UNIQUE NOT NULL,
    avatar_url VARCHAR(255),
    username VARCHAR(32) UNIQUE NOT NULL,
    full_name VARCHAR(64),
    bio VARCHAR(2000),
    pass_hash VARCHAR(64) NOT NULL,
    last_login_at TIMESTAMP,
    role role NOT NULL DEFAULT 'user',
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

---- create above / drop below ----

DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS role;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
