-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS managers (
    id SERIAL PRIMARY KEY,
    login VARCHAR UNIQUE NOT NULL,
    hashed_password VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS specialists (
    id SERIAL PRIMARY KEY,
    login VARCHAR UNIQUE NOT NULL,
    hashed_password VARCHAR NOT NULL,
    fullname VARCHAR,
    level INTEGER DEFAULT (1),
    photo_url VARCHAR NOT NULL,
    is_verified BOOLEAN DEFAULT(FALSE) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS managers, specialists;
-- +goose StatementEnd
