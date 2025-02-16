package model

import (
	"context"
	"time"
)

type IOrderRepository interface {
	FindAll(ctx context.Context, userID int64) ([]*Order, error)
	FindById(ctx context.Context, id string) (*Order, error)
	SaveOrder(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id string) error
}

type IOrderUsecase interface {
	FindAll(ctx context.Context, userID int64) ([]*Order, error)
	FindById(ctx context.Context, id string) (*Order, error)
	Create(ctx context.Context, in CreateOrderInput) error
	Delete(ctx context.Context, id string) error
}

type Order struct {
	ID          string      `json:"id"`
	UserID      int64       `json:"user_id"`
	TotalAmount float64     `json:"total_amount"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   *time.Time  `json:"-"`
	OrderItems  []OrderItem `json:"order_items"`
}

type OrderItem struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   string     `json:"order_id" gorm:"index"`
	ProductID int64      `json:"product_id"`
	Quantity  int        `json:"quantity"`
	Price     float64    `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type CreateOrderInput struct {
	UserID     int64             `json:"user_id" validate:"required"`
	OrderItems []CreateOrderItem `json:"order_items" validate:"required"`
}

type CreateOrderItem struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int   `json:"quantity" validate:"required"`
}
