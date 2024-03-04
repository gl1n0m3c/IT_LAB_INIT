-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS violation (
    id VARCHAR UNIQUE PRIMARY KEY,
    type VARCHAR NOT NULL,
    amount REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS cameras (
    id VARCHAR UNIQUE PRIMARY KEY,
    type VARCHAR(10) NOT NULL,
    description VARCHAR NOT NULL,
    coordinates VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS contacts (
    transport VARCHAR(8) UNIQUE PRIMARY KEY ,
    contacts JSONB NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS violation, cameras, contacts;
-- +goose StatementEnd
