-- +goose Up
-- +goose StatementBegin
CREATE TYPE status_type AS ENUM ('Correct', 'Incorrect', 'Unknown');

CREATE TABLE IF NOT EXISTS cases (
    id SERIAL PRIMARY KEY,
    camera_id VARCHAR NOT NULL,
    transport VARCHAR(20) NOT NULL,
    violation_id VARCHAR NOT NULL,
    violation_value VARCHAR NOT NULL,
    level INTEGER NOT NULL,
    datetime TIMESTAMP WITH TIME ZONE NOT NULL,
    photo_url VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS rated_cases (
    id SERIAL PRIMARY KEY,
    specialist_id INTEGER NOT NULL,
    case_id INTEGER NOT NULL,
    choice BOOLEAN NOT NULL,
    date DATE NOT NULL,
    status status_type DEFAULT('Unknown') NOT NULL,
    UNIQUE(specialist_id, case_id)
);

ALTER TABLE cases
    ADD CONSTRAINT fk_camera
        FOREIGN KEY (camera_id) REFERENCES cameras(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_transport
        FOREIGN KEY (transport) REFERENCES contacts(transport) ON DELETE CASCADE,
    ADD CONSTRAINT fk_violation
        FOREIGN KEY (violation_id) REFERENCES violations(id) ON DELETE CASCADE;

ALTER TABLE rated_cases
    ADD CONSTRAINT fk_specialist
        FOREIGN KEY (specialist_id) REFERENCES specialists(id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_case
        FOREIGN KEY (case_id) REFERENCES cases(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rated_cases, cases;
DROP TYPE IF EXISTS status_type;
-- +goose StatementEnd
