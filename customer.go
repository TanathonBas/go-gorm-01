package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func createCustomerDB(db *gorm.DB, customer *Customer) error {
	if customer == nil {
		return errors.New("customer is required")
	}
	if strings.TrimSpace(customer.Name) == "" || strings.TrimSpace(customer.Email) == "" || strings.TrimSpace(customer.Phone) == "" {
		return errors.New("customer name, email and phone are required")
	}
	customer.Name = strings.TrimSpace(customer.Name)
	customer.Email = strings.ToLower(strings.TrimSpace(customer.Email))
	customer.Phone = strings.TrimSpace(customer.Phone)
	return db.Create(customer).Error
}

func getCustomerByEmailDB(db *gorm.DB, email string) (*Customer, error) {
	var customer Customer
	err := db.Where("email = ?", strings.ToLower(strings.TrimSpace(email))).First(&customer).Error
	return &customer, err
}

func updateCustomerByIDDB(db *gorm.DB, id uint, update *Customer) error {
	if update == nil {
		return errors.New("customer is required")
	}
	if strings.TrimSpace(update.Name) == "" || strings.TrimSpace(update.Email) == "" || strings.TrimSpace(update.Phone) == "" {
		return errors.New("customer name, email and phone are required")
	}
	result := db.Model(&Customer{}).Where("id = ?", id).Updates(map[string]any{
		"name": strings.TrimSpace(update.Name), "email": strings.ToLower(strings.TrimSpace(update.Email)), "phone": strings.TrimSpace(update.Phone),
	})
	return rowsAffectedError(result, "customer", id)
}

func deleteCustomerByIDDB(db *gorm.DB, id uint) error {
	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		return err
	}
	var count int64
	if err := db.Unscoped().Model(&Order{}).Where("customer_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete customer %d: customer has %d order(s)", id, count)
	}
	return db.Delete(&customer).Error
}
