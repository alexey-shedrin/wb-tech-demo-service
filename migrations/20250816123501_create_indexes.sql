-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_orders_order_uid ON orders (order_uid);
CREATE INDEX idx_deliveries_order_uid ON deliveries (order_uid);
CREATE INDEX idx_payments_order_uid ON payments (order_uid);
CREATE INDEX idx_items_order_uid ON items (order_uid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_deliveries_order_uid;
DROP INDEX IF EXISTS idx_payments_order_uid;
DROP INDEX IF EXISTS idx_items_order_uid;
DROP INDEX IF EXISTS idx_orders_order_uid;
-- +goose StatementEnd
