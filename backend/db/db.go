package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var persisted_db *gorm.DB

func Conn() *gorm.DB {
	if persisted_db != nil {
		return persisted_db
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// // Create
	// db.Create(&Product{Code: "D42", Price: 100})

	// // Read
	// var product Product
	// db.First(&product, 1)                 // find product with integer primary key
	// db.First(&product, "code = ?", "D42") // find product with code D42

	// slog.Info("Product: ", product)

	persisted_db = db
	return persisted_db
}

func Close() {
	if persisted_db != nil {
		db, err := persisted_db.DB()
		if err != nil {
			panic("failed to get db")
		}
		db.Close()
	}
}
