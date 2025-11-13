-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMP WITH TIME ZONE,
    user_id         BIGINT,
    status          TEXT,
    source          TEXT,
    payment         TEXT,
    metadata        JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
