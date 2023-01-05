-- +goose Up
-- +goose StatementBegin
CREATE TABLE rate_limiters (
    id VARCHAR(255),
    tokens INT,
    limit_tokens INT,
    interval INT,
    last_reserved_at BIGINT,
    last_recalculated_at BIGINT,
    PRIMARY KEY (id)
);

-- UPDATE rate_limiters SET
--     tokens = tokens + TRUNC((? - last_recalculated_at) / interval),
--     last_recalculated_at = ?
-- WHERE id IN (
--     SELECT id FROM rate_limiters WHERE id = '2' FOR UPDATE
-- ) RETURNING *;
--
-- UPDATE rate_limiters SET tokens = tokens - 1, last_reserved_at = ? WHERE id = '2';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rate_limiters;
-- +goose StatementEnd
