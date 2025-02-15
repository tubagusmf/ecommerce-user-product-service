package repository

import (
	"context"

	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) model.ICategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll(ctx context.Context, category model.Category) ([]*model.Category, error) {
	var categories []*model.Category
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) FindById(ctx context.Context, id int64) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(ctx context.Context, category model.Category) error {
	err := r.db.WithContext(ctx).Create(&category).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) Update(ctx context.Context, category model.Category) error {
	err := r.db.WithContext(ctx).Model(&model.Category{}).
		Where("id = ? AND deleted_at IS NULL", category.ID).
		Updates(category).Error

	if err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.Category{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
