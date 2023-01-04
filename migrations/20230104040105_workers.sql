-- +goose Up
-- +goose StatementBegin
CREATE TABLE workers (
    id CHAR(36),
    last_ping_at BIGINT,
    PRIMARY KEY (id)
);
ALTER TABLE requests ADD COLUMN worker_id CHAR(36) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workers;
ALTER TABLE requests DROP COLUMN worker_id;
-- +goose StatementEnd
