package repository

import (
	"context"
	"github.com/gofrs/uuid"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	CreateOrder(ctx context.Context, order *Order) error
	GetOrders(ctx context.Context, orderStatus string, archive bool, userID uuid.UUID, limit, page int) ([]Order, error)
	GetOrderCount(ctx context.Context, userID uuid.UUID) (int, error)
	CancelOrder(ctx context.Context, orderConsignmentID string) error
	Close()
}
