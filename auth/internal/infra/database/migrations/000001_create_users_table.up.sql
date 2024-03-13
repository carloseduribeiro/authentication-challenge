CREATE SCHEMA auth;

CREATE TYPE auth.USER_TYPE AS ENUM ('default', 'admin');

CREATE TABLE IF NOT EXISTS auth.users
(
    id        UUID PRIMARY KEY,
    document  VARCHAR(11) UNIQUE  NOT NULL,
    name      VARCHAR(255)        NOT NULL,
    email     VARCHAR(300) UNIQUE NOT NULL,
    password  VARCHAR(72)         NOT NULL,
    birthDate DATE                NOT NULL,
    type      auth.USER_TYPE      NOT NULL default 'default'
);