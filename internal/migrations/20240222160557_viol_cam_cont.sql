-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS violations (
    id VARCHAR UNIQUE PRIMARY KEY,
    type VARCHAR UNIQUE NOT NULL,
    amount REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS cameras (
    id VARCHAR UNIQUE PRIMARY KEY,
    type VARCHAR(10) NOT NULL,
    description VARCHAR NOT NULL,
    coordinates VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS contacts (
    transport VARCHAR(20) UNIQUE PRIMARY KEY ,
    contacts JSONB NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS violations, cameras, contacts;
-- +goose StatementEnd
