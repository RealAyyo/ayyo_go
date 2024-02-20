-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
    ADD COLUMN description TEXT,
ADD COLUMN notification_time INTERVAL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
DROP COLUMN description,
DROP COLUMN notification_time;
-- +goose StatementEnd