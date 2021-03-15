package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresProduct struct {
	gorm.Model
	Code  string
	Price uint
	Code2 string
}

func postgresTest() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres password=12345 dbname=dbname port=5432 sslmode=disable TimeZone=Asia/Shanghai", // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true,                                                                                                         // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&PostgresProduct{})
	if err != nil {
		panic(err)
	}
	// Create
	db.Create(&PostgresProduct{Code: "D42", Price: 100})

	// Read
	var product PostgresProduct
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	db.Delete(&product)

}
