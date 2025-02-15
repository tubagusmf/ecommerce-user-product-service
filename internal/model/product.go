package model

import (
	"context"
	"time"
)

type IProductRepository interface {
	FindAll(ctx context.Context, filter FindAllParam) ([]*Product, error)
	FindById(ctx context.Context, id int64) (*Product, error)
	Create(ctx context.Context, product Product) error
	Update(ctx context.Context, product Product) error
	Delete(ctx context.Context, id int64) error
}

type IProductUsecase interface {
	FindAll(ctx context.Context, filter FindAllParam) ([]*Product, error)
	FindById(ctx context.Context, id int64) (*Product, error)
	Create(ctx context.Context, in CreateProductInput) error
	Update(ctx context.Context, id int64, in UpdateProductInput) error
	Delete(ctx context.Context, id int64) error
}

type Product struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Price        int64      `json:"price"`
	Stock        int64      `json:"stock"`
	CategoryID   int64      `json:"category_id"`
	CategoryName string     `json:"category_name"`
	ImageUrl     string     `json:"image_url"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"-"`
}

type FindAllParam struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type CreateProductInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required"`
	Stock       int64  `json:"stock" validate:"required"`
	CategoryID  int64  `json:"category_id" validate:"required"`
	ImageUrl    string `json:"image_url" validate:"required"`
}

type UpdateProductInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required"`
	Stock       int64  `json:"stock" validate:"required"`
	CategoryID  int64  `json:"category_id" validate:"required"`
	ImageUrl    string `json:"image_url" validate:"required"`
}
