package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
)

var (
	db *gorm.DB
)

const (
	host = "localhost"
	port = "5432"
	user = "postgres_test"
	password = "Mir@ge308"
	dbName = "TaskManagement"
)

func Connection() *gorm.DB{
	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	d, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{})
	if err != nil{
		panic(fmt.Sprintf("Error in Database Connection: %v", err)) // Better visibility
	}
	db = d
	
	fmt.Println("DB is connecteddddddddddddddddddddd!")
	return db
}
