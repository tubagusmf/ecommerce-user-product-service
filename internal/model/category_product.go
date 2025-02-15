package model

import (
	"context"
	"time"
)

type ICategoryRepository interface {
	FindAll(ctx context.Context, category Category) ([]*Category, error)
	FindById(ctx context.Context, id int64) (*Category, error)
	Create(ctx context.Context, category Category) error
	Update(ctx context.Context, category Category) error
	Delete(ctx context.Context, id int64) error
}

type ICategoryUsecase interface {
	FindAll(ctx context.Context, category Category) ([]*Category, error)
	FindById(ctx context.Context, id int64) (*Category, error)
	Create(ctx context.Context, in CreateCategoryInput) error
	Update(ctx context.Context, id int64, in UpdateCategoryInput) error
	Delete(ctx context.Context, id int64) error
}

type Category struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type CreateCategoryInput struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategoryInput struct {
	Name string `json:"name" validate:"required"`
}
