-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title TEXT,
    date TIMESTAMP,
    duration INTERVAL,
    user_id INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events
-- +goose StatementEnd
