CREATE TABLE IF NOT EXISTS auth.sessions
(
    id         UUID      NOT NULL,
    user_id    UUID      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES auth.users (id)
);