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
    level INTEGER DEFAULT (1) NOT NULL ,
    photo_url VARCHAR,
    row INTEGER DEFAULT (0) NOT NULL,
    current_row INTEGER DEFAULT (0) NOT NULL,
    is_verified BOOLEAN DEFAULT(FALSE) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS managers, specialists;
-- +goose StatementEnd
