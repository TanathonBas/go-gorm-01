package main

import "gorm.io/gorm"

const (
	OrderStatusPending   = "Pending"
	OrderStatusCompleted = "Completed"
	OrderStatusCancelled = "Cancelled"
)

type Category struct {
	gorm.Model
	Name        string    `gorm:"size:100;not null;uniqueIndex"`
	Description string    `gorm:"size:255"`
	Products    []Product `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type Product struct {
	gorm.Model
	Name       string   `gorm:"size:150;not null;index"`
	Price      float64  `gorm:"type:decimal(12,2);not null"`
	Quantity   int      `gorm:"not null"`
	CategoryID uint     `gorm:"not null;index"`
	Category   Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Orders     []Order  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type Customer struct {
	gorm.Model
	Name   string  `gorm:"size:150;not null"`
	Email  string  `gorm:"size:150;not null;uniqueIndex"`
	Phone  string  `gorm:"size:30;not null"`
	Orders []Order `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type Order struct {
	gorm.Model
	CustomerID uint     `gorm:"not null;index"`
	Customer   Customer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ProductID  uint     `gorm:"not null;index"`
	Product    Product  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Quantity   int      `gorm:"not null"`
	TotalPrice float64  `gorm:"type:decimal(12,2);not null"`
	Status     string   `gorm:"size:20;not null;default:Pending;index"`
}
