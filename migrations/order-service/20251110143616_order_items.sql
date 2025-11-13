-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_items (
    id          BIGSERIAL PRIMARY KEY,
    order_id    BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    sku         TEXT NOT NULL,
    store_id    BIGINT NOT NULL, -- склад/точка, откуда списан товар
    qty         INT NOT NULL,
    unit_price  NUMERIC(10,2) NOT NULL,
    total_price NUMERIC(10,2) GENERATED ALWAYS AS (qty * unit_price) STORED
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items;
-- +goose StatementEnd
