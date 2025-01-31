package repository

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type PgxRepository struct {
	db *pgxpool.Pool
}

var (
	once       sync.Once
	repository *PgxRepository
)

func NewPgRepository(databaseUrl string) (*PgxRepository, error) {
	var onceErr error
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		config, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			onceErr = fmt.Errorf("invalid database URL: %w", err)
			log.Error().Err(err).Msg("Failed to parse database configuration")
			return
		}

		config.MaxConns = 10
		config.MinConns = 2
		config.MaxConnLifetime = 30 * time.Minute
		config.MaxConnIdleTime = 5 * time.Second
		config.HealthCheckPeriod = 1 * time.Minute

		db, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			onceErr = fmt.Errorf("failed to create connection pool: %w", err)
			log.Error().Err(err).Msg("Database Connection Error")
			return
		}

		if err = db.Ping(ctx); err != nil {
			onceErr = fmt.Errorf("failed to ping database: %w", err)
			log.Error().Err(err).Msg("Database Ping Error")
			db.Close()
			return
		}

		repository = &PgxRepository{db: db}
		log.Info().Msg("Database connection pool successfully initialized")
	})

	return repository, onceErr
}

func (repo *PgxRepository) Close() {
	if repo.db != nil {
		repo.db.Close()
		log.Info().Msg("Database connection pool closed")
	}
}

func (repo *PgxRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	err := repo.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PgxRepository) CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, name, email, password,created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`
	err := repo.db.QueryRow(ctx, query, user.ID, user.Name, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PgxRepository) CreateOrder(ctx context.Context, order *Order) error {
	query := `INSERT INTO orders (
		user_id, order_consignment_id, store_id, merchant_order_id, recipient_name, 
		recipient_phone, recipient_address, recipient_city, recipient_zone, recipient_area, 
		delivery_type, item_type, special_instruction, item_quantity, item_weight, 
		amount_to_collect, item_description, total_fee, order_type_id, cod_fee, 
		promo_discount, discount, delivery_fee, order_status, archive, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, 
		$17, $18, $19, $20, $21, $22, $23, $24, $25, NOW(), NOW()
	) RETURNING id`

	err := repo.db.QueryRow(ctx, query,
		order.UserID, order.OrderConsignmentID, order.StoreID, order.MerchantOrderID,
		order.RecipientName, order.RecipientPhone, order.RecipientAddress, order.RecipientCity,
		order.RecipientZone, order.RecipientArea, order.DeliveryType, order.ItemType,
		order.SpecialInstruction, order.ItemQuantity, order.ItemWeight, order.AmountToCollect,
		order.ItemDescription, order.TotalFee, order.OrderTypeID, order.CODFee, order.PromoDiscount,
		order.Discount, order.DeliveryFee, order.OrderStatus, order.Archive,
	).Scan(&order.ID)

	if err != nil {
		return err
	}

	return nil

}

func (repo *PgxRepository) GetOrders(ctx context.Context, orderStatus string, archive bool, userID uuid.UUID, limit, page int) ([]Order, error) {
	offset := (page - 1) * limit

	query := `SELECT 
		order_consignment_id, created_at, item_description, merchant_order_id, 
		recipient_name, recipient_address, recipient_phone, amount_to_collect, 
		total_fee, special_instruction, order_type_id, cod_fee, promo_discount, 
		discount, delivery_fee,item_type,order_status
	FROM orders
	WHERE user_id = $1 AND order_status = $2 AND archive = $3
	ORDER BY created_at DESC
	LIMIT $4 OFFSET $5`

	rows, err := repo.db.Query(ctx, query, userID, orderStatus, archive, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order

	for rows.Next() {
		var order Order
		err2 := rows.Scan(
			&order.OrderConsignmentID, &order.CreatedAt, &order.ItemDescription, &order.MerchantOrderID,
			&order.RecipientName, &order.RecipientAddress, &order.RecipientPhone, &order.AmountToCollect,
			&order.TotalFee, &order.SpecialInstruction, &order.OrderTypeID, &order.CODFee,
			&order.PromoDiscount, &order.Discount, &order.DeliveryFee, &order.ItemType, &order.OrderStatus,
		)
		if err2 != nil {
			return nil, err2
		}
		orders = append(orders, order)
	}

	if err3 := rows.Err(); err3 != nil {
		return nil, err3
	}

	return orders, nil
}

func (repo *PgxRepository) GetOrderCount(ctx context.Context, userID uuid.UUID, orderStatus string, archive bool) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM orders 
		WHERE user_id = $1 AND order_status = $2 AND archive = $3`

	var count int
	err := repo.db.QueryRow(ctx, query, userID, orderStatus, archive).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *PgxRepository) CancelOrder(ctx context.Context, orderConsignmentID string) error {

	query := `UPDATE orders 
			  SET order_status = $1 
			  WHERE order_consignment_id = $2 AND order_status != $1`

	result, err := repo.db.Exec(ctx, query, "Cancelled", orderConsignmentID)
	if err != nil {
		return fmt.Errorf("failed to cancel order with consignment ID %s: %w", orderConsignmentID, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}
