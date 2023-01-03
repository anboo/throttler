-- +goose Up
-- +goose StatementBegin
CREATE TABLE requests (
    id           CHAR(36),
    status       VARCHAR(50),
    created_at   TIMESTAMP,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE requests;
-- +goose StatementEnd
