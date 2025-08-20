package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/config"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/model"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/util"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Postgres struct {
	db *sql.DB
}

func New(cfg *config.Postgres) *Postgres {
	drivStr := "postgres"
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	db, err := sql.Open(drivStr, connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	err = goose.SetDialect(drivStr)
	if err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return &Postgres{db: db}
}

func (pg *Postgres) Close() {
	err := pg.db.Close()
	if err != nil {
		log.Printf("failed to close database: %v", err)
	}
}

func (pg *Postgres) SaveOrder(order *model.Order) error {
	tx, err := pg.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return util.ErrInternal
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	orderQuery := `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(orderQuery,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		log.Printf("failed to insert order: %v", err)
		return util.ErrInternal
	}

	deliveryQuery := `
		INSERT INTO deliveries (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(deliveryQuery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		log.Printf("failed to insert delivery: %v", err)
		return util.ErrInternal
	}

	paymentQuery := `
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(paymentQuery,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		log.Printf("failed to insert payment: %v", err)
		return util.ErrInternal
	}

	itemQuery := `
		INSERT INTO items (
			order_uid, chrt_id, track_number, price, rid, name,
			sale, size, total_price, nm_id, brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, item := range order.Items {
		_, err = tx.Exec(itemQuery,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price,
			item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			log.Printf("failed to insert item: %v", err)
			return util.ErrInternal
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("failed to commit transaction: %v", err)
		return util.ErrInternal
	}

	return nil
}

func (pg *Postgres) GetOrderByUID(orderUID string) (*model.Order, error) {
	var order model.Order

	orderQuery := `
		SELECT order_uid, track_number, entry, locale, internal_signature,
			   customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders 
		WHERE order_uid = $1`

	err := pg.db.QueryRow(orderQuery, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.ErrOrderNotFound
		}
		log.Printf("failed to get order: %v", err)
		return nil, util.ErrInternal
	}

	deliveryQuery := `
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries 
		WHERE order_uid = $1`

	err = pg.db.QueryRow(deliveryQuery, orderUID).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil {
		log.Printf("failed to get delivery: %v", err)
		return nil, util.ErrInternal
	}

	paymentQuery := `
		SELECT transaction, request_id, currency, provider, amount,
			   payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments 
		WHERE order_uid = $1`

	err = pg.db.QueryRow(paymentQuery, orderUID).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil {
		log.Printf("failed to get payment: %v", err)
		return nil, util.ErrInternal
	}

	itemsQuery := `
		SELECT chrt_id, track_number, price, rid, name, sale,
			   size, total_price, nm_id, brand, status
		FROM items 
		WHERE order_uid = $1`

	rows, err := pg.db.Query(itemsQuery, orderUID)
	if err != nil {
		log.Printf("failed to query items: %v", err)
		return nil, util.ErrInternal
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			log.Printf("failed to scan item: %v", err)
			return nil, util.ErrInternal
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		log.Printf("failed to iterate items: %v", err)
		return nil, util.ErrInternal
	}

	order.Items = items

	return &order, nil
}
