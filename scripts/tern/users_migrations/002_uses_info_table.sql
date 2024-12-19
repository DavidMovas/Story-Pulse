-- Write your migrate up statements here

CREATE TABLE users_info (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    articles_count SMALLINT,
    comments_count SMALLINT,
    followers_count SMALLINT,
    following_count SMALLINT
);

---- create above / drop below ----

DROP TABLE IF EXISTS users_info;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
