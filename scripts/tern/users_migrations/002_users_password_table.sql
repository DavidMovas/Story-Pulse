-- Write your migrate up statements here

CREATE TABLE users_password (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    pass_hash VARCHAR(64) NOT NULL,
    updated_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_password_user_id ON users_password(user_id);

---- create above / drop below ----

DROP INDEX idx_users_password_user_id;
DROP TABLE users_password;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
