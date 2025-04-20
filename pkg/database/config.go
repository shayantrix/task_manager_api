package config

import(
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
)

const(
	host = "localhost"
	port = "5432"
	user = "postgres_test"
	password = "Mir@ge308"
	dbName = "TaskManagement"
)

func Connection() *gorm.DB{
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{})
	if err != nil{
		fmt.Errorf("Error in Database Connection: ", err)
	}
	return db
}
