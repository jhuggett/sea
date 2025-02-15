package db

import (
	"github.com/jhuggett/sea/data"
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
		panic("failed to connect data.base")
	}

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

func Migrate() {
	db := Conn()

	db.AutoMigrate(&data.Ship{})
	db.AutoMigrate(&data.WorldMap{})
	db.AutoMigrate(&data.Point{})
	db.AutoMigrate(&data.Continent{})
	db.AutoMigrate(&data.Port{})
	db.AutoMigrate(&data.Crew{})
	db.AutoMigrate(&data.Inventory{})
	db.AutoMigrate(&data.Item{})
	db.AutoMigrate(&data.Population{})
	db.AutoMigrate(&data.Industry{})
	db.AutoMigrate(&data.Person{})
	db.AutoMigrate(&data.EmploymentTerms{})
	db.AutoMigrate(&data.Contract{})
	db.AutoMigrate(&data.Fleet{})
	db.AutoMigrate(&data.Deed{})
	db.AutoMigrate(&data.Economy{})
	db.AutoMigrate(&data.Market{})

}

func SetupInMemDB() {
	persisted_db = nil

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect data.base")
	}

	persisted_db = db

	Migrate()
}

type Scope func(db *gorm.DB) *gorm.DB
