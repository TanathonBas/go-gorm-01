package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DBUser = "root"
	DBPass = ""
	DBHost = "127.0.0.1"
	DBPort = 3306
	DBName = "store_db"
)

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal("connect database: ", err)
	}
	fmt.Println("Database connection successful")

	if err := db.AutoMigrate(&Category{}, &Product{}, &Customer{}, &Order{}); err != nil {
		log.Fatal("migrate database: ", err)
	}
	fmt.Println("Database migration successful")

	products, err := getAllProductsDB(db)
	if err != nil {
		log.Println("get all products:", err)
	} else {
		fmt.Printf("All products: %d item(s)\n", len(products))
		for _, product := range products {
			fmt.Printf("ID: %d | Name: %s | Price: %.2f | Stock: %d | Category: %s\n",
				product.ID, product.Name, product.Price, product.Quantity, product.Category.Name)
		}
	}

}

func connectDB() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	serverDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", DBUser, DBPass, DBHost, DBPort)
	serverDB, err := gorm.Open(mysql.Open(serverDSN), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, err
	}
	if err := serverDB.Exec("CREATE DATABASE IF NOT EXISTS `" + DBName + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error; err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&clientFoundRows=true", DBUser, DBPass, DBHost, DBPort, DBName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
}
