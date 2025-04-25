package config_test

import (
	"testing"
	"github.com/shayantrix/task_manager_api/pkg/config"
)

func TestConnection(t *testing.T) *gorm.DB{
	host := "localhost"
	user := "postgres_test"
	password := "123456m."
	dbName := "TaskManagement"

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	d, err := gorm.Open(postgres.Open(dns), &gorm.Config{})

	for err != nil{
		log.Panic("Fatal error in connection to DB: ", err)
	}
	db = d
	return db
}
