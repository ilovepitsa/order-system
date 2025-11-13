-- +goose Up
-- +goose StatementBegin
CREATE TABLE inventory (
  store_id    BIGINT,
  sku         BIGINT,
  qty         INT,
  created_at  TIMESTAMP WITH TIME ZONE,
  updated_at  TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY (store_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE inventory;
-- +goose StatementEnd
