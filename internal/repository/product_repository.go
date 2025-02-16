package repository

import (
	"context"
	"errors"

	"github.com/tubagusmf/ecommerce-user-product-service/internal/model"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) model.IProductRepository {
	return &ProductRepo{
		db: db,
	}
}

func (r *ProductRepo) FindAll(ctx context.Context, filter model.FindAllParam) ([]*model.Product, error) {
	var products []*model.Product

	query := r.db.WithContext(ctx).
		Table("products").
		Select("products.*, categories.name as category_name").
		Joins("LEFT JOIN categories ON categories.id = products.category_id").
		Where("products.deleted_at IS NULL")

	if filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Limit(int(filter.Limit)).Offset(int(offset))
	}

	err := query.Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepo) FindById(ctx context.Context, id int64) (*model.Product, error) {
	var product model.Product

	err := r.db.WithContext(ctx).
		Table("products").
		Select("products.*, categories.name as category_name").
		Joins("LEFT JOIN categories ON categories.id = products.category_id").
		Where("products.id = ? AND products.deleted_at IS NULL", id).
		First(&product).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepo) Create(ctx context.Context, product model.Product) error {
	err := r.db.WithContext(ctx).
		Omit("CategoryName").
		Create(&product).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) Update(ctx context.Context, product model.Product) error {
	err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ? AND deleted_at IS NULL", product.ID).
		Omit("CategoryName").
		Updates(product).Error

	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) GetPriceByID(ctx context.Context, productID int64, price *float64) error {
	return r.db.WithContext(ctx).Model(&model.Product{}).
		Select("price").
		Where("id = ?", productID).
		Scan(price).Error
}
