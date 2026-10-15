package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func createCategoryDB(db *gorm.DB, category *Category) error {
	if category == nil || strings.TrimSpace(category.Name) == "" {
		return errors.New("category name is required")
	}
	category.Name = strings.TrimSpace(category.Name)
	return db.Create(category).Error
}

func deleteCategoryByIDDB(db *gorm.DB, id uint) error {
	var category Category
	if err := db.First(&category, id).Error; err != nil {
		return err
	}

	var count int64
	if err := db.Unscoped().Model(&Product{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete category %d: it still has %d product(s)", id, count)
	}
	return db.Delete(&category).Error
}
