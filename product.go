package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func createProductDB(db *gorm.DB, product *Product) error {
	if product == nil {
		return errors.New("product is required")
	}
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("product name is required")
	}
	if product.Price < 0 || product.Quantity < 0 {
		return errors.New("price and quantity must not be negative")
	}
	if err := db.First(&Category{}, product.CategoryID).Error; err != nil {
		return fmt.Errorf("category not found: %w", err)
	}
	product.Name = strings.TrimSpace(product.Name)
	return db.Create(product).Error
}

func getAllProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").Find(&products).Error
	return products, err
}

func getProductByIDDB(db *gorm.DB, id uint) (*Product, error) {
	var product Product
	err := db.Preload("Category").First(&product, id).Error
	return &product, err
}

func getProductByNameDB(db *gorm.DB, name string) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").Where("name LIKE ?", "%"+name+"%").Find(&products).Error
	return products, err
}

func getProductsByCategoryDB(db *gorm.DB, categoryID uint) ([]Product, error) {
	var products []Product
	err := db.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

func getProductsByPriceRangeDB(db *gorm.DB, minPrice, maxPrice float64) ([]Product, error) {
	if minPrice < 0 || maxPrice < minPrice {
		return nil, errors.New("invalid price range")
	}
	var products []Product
	err := db.Preload("Category").Where("price BETWEEN ? AND ?", minPrice, maxPrice).Find(&products).Error
	return products, err
}

func getProductsLowStockDB(db *gorm.DB, threshold int) ([]Product, error) {
	if threshold < 0 {
		return nil, errors.New("stock threshold must not be negative")
	}
	var products []Product
	err := db.Preload("Category").Where("quantity < ?", threshold).Find(&products).Error
	return products, err
}

func getProductsByCategoryAndPriceDB(db *gorm.DB, categoryID uint, minPrice, maxPrice float64) ([]Product, error) {
	if minPrice < 0 || maxPrice < minPrice {
		return nil, errors.New("invalid price range")
	}
	var products []Product
	err := db.Preload("Category").Where("category_id = ? AND price BETWEEN ? AND ?", categoryID, minPrice, maxPrice).Find(&products).Error
	return products, err
}

func getProductsSortedByPriceDB(db *gorm.DB, direction string) ([]Product, error) {
	direction = strings.ToUpper(strings.TrimSpace(direction))
	if direction != "ASC" && direction != "DESC" {
		return nil, errors.New("sort direction must be ASC or DESC")
	}
	var products []Product
	err := db.Preload("Category").Order("price " + direction).Find(&products).Error
	return products, err
}

func getTopExpensiveProductsDB(db *gorm.DB, limit int) ([]Product, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than zero")
	}
	var products []Product
	err := db.Preload("Category").Order("price DESC").Limit(limit).Find(&products).Error
	return products, err
}

func updateProductByIDDB(db *gorm.DB, id uint, update *Product) error {
	if update == nil || strings.TrimSpace(update.Name) == "" {
		return errors.New("product name is required")
	}
	if update.Price < 0 || update.Quantity < 0 {
		return errors.New("price and quantity must not be negative")
	}
	if err := db.First(&Category{}, update.CategoryID).Error; err != nil {
		return fmt.Errorf("category not found: %w", err)
	}
	result := db.Model(&Product{}).Where("id = ?", id).Updates(map[string]any{
		"name": update.Name, "price": update.Price, "quantity": update.Quantity, "category_id": update.CategoryID,
	})
	return rowsAffectedError(result, "product", id)
}

func updateProductPriceDB(db *gorm.DB, id uint, price float64) error {
	if price < 0 {
		return errors.New("price must not be negative")
	}
	result := db.Model(&Product{}).Where("id = ?", id).Update("price", price)
	return rowsAffectedError(result, "product", id)
}

func updateProductStockDB(db *gorm.DB, id uint, quantityChange int) error {
	result := db.Model(&Product{}).
		Where("id = ? AND quantity + ? >= 0", id, quantityChange).
		UpdateColumn("quantity", gorm.Expr("quantity + ?", quantityChange))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("product not found or resulting stock would be negative")
	}
	return nil
}

func softDeleteProductDB(db *gorm.DB, id uint) error {
	result := db.Delete(&Product{}, id)
	return rowsAffectedError(result, "product", id)
}

func hardDeleteProductDB(db *gorm.DB, id uint) error {
	var count int64
	if err := db.Unscoped().Model(&Order{}).Where("product_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot permanently delete product %d: it is used by %d order(s)", id, count)
	}
	result := db.Unscoped().Delete(&Product{}, id)
	return rowsAffectedError(result, "product", id)
}

func getDeletedProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Unscoped().Preload("Category").Where("products.deleted_at IS NOT NULL").Find(&products).Error
	return products, err
}

func restoreProductDB(db *gorm.DB, id uint) error {
	result := db.Unscoped().Model(&Product{}).Where("id = ? AND deleted_at IS NOT NULL", id).Update("deleted_at", nil)
	return rowsAffectedError(result, "deleted product", id)
}

func rowsAffectedError(result *gorm.DB, entity string, id uint) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%s with id %d not found", entity, id)
	}
	return nil
}
