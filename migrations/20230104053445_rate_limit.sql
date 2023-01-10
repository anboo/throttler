-- +goose Up
-- +goose StatementBegin
CREATE TABLE rate_limiters (
    id VARCHAR(255),
    tokens INT CHECK (tokens BETWEEN 0 AND limit_tokens),
    limit_tokens INT,
    interval INT CHECK (interval > 0),
    last_reserved_at BIGINT,
    last_recalculated_at BIGINT,
    PRIMARY KEY (id)
);

-- Function
CREATE OR REPLACE FUNCTION take_token (id VARCHAR(100)) RETURNS RECORD AS $$
DECLARE
    limit_tokens INTEGER;
    tokens INTEGER;
    extra_tokens INTEGER;
    new_tokens INTEGER;
    interval INTEGER;
    last_recalculated_at BIGINT;
    this_recalculated_at BIGINT;
BEGIN
    SELECT b.tokens, b.last_recalculated_at, b.interval, b.limit_tokens
    INTO tokens, last_recalculated_at, interval, limit_tokens
    FROM rate_limiters b WHERE b.id = $1
        FOR UPDATE;

    IF limit_tokens IS NULL THEN
        raise notice 'Id % rate not configured', $1;
        RETURN (false, CAST(-1 AS BIGINT));
    END IF;

    extra_tokens := FLOOR(
        (EXTRACT(EPOCH FROM now()) * 1000 - last_recalculated_at) / interval
    )::int;
    this_recalculated_at := FLOOR(EXTRACT(EPOCH FROM now()) * 1000);
    new_tokens := LEAST(limit_tokens, tokens + extra_tokens);
    raise notice 'Id % has % tokens last refill %', $1, new_tokens, this_recalculated_at;

    IF new_tokens <= 0 THEN
        RETURN (false, @this_recalculated_at - last_recalculated_at - interval);
    END IF;

    UPDATE rate_limiters b SET (tokens, last_recalculated_at) = (new_tokens - 1, this_recalculated_at) WHERE b.id = $1;
    RETURN (true, CAST(-1 AS BIGINT));
END
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rate_limiters;
DROP FUNCTION IF EXISTS take_token(VARCHAR);
-- +goose StatementEnd
