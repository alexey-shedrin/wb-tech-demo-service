-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    order_uid         VARCHAR(255) PRIMARY KEY,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    track_number      VARCHAR(255),
    entry             VARCHAR(255),
    locale            VARCHAR(255),
    internal_signature VARCHAR(255),
    customer_id       VARCHAR(255),
    delivery_service  VARCHAR(255),
    shardkey          VARCHAR(255),
    sm_id             INTEGER,
    date_created      VARCHAR(255),
    oof_shard         VARCHAR(255)
);

CREATE TABLE deliveries (
    id BIGSERIAL PRIMARY KEY,
    order_uid     VARCHAR(255) REFERENCES orders (order_uid) ON DELETE CASCADE,
    name        VARCHAR(255),
    phone       VARCHAR(255),
    zip         VARCHAR(255),
    city        VARCHAR(255),
    address     VARCHAR(255),
    region      VARCHAR(255),
    email       VARCHAR(255)
);

CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    order_uid     VARCHAR(255) REFERENCES orders (order_uid) ON DELETE CASCADE,
    transaction   VARCHAR(255),
    request_id    VARCHAR(255),
    currency      VARCHAR(255),
    provider      VARCHAR(255),
    amount        INTEGER,
    payment_dt    BIGINT,
    bank          VARCHAR(255),
    delivery_cost INTEGER,
    goods_total   INTEGER,
    custom_fee    INTEGER
);

CREATE TABLE items (
    id BIGSERIAL PRIMARY KEY,
    order_uid     VARCHAR(255) REFERENCES orders (order_uid) ON DELETE CASCADE,
    chrt_id      BIGINT,
    track_number VARCHAR(255),
    price        INTEGER,
    rid          VARCHAR(255),
    name         VARCHAR(255),
    sale         INTEGER,
    size         VARCHAR(255),
    total_price  INTEGER,
    nm_id        BIGINT,
    brand        VARCHAR(255),
    status       INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
