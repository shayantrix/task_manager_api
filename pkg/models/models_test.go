package models_test

import (
	"testing"
	"github.com/shayantrix/task_manager_api/pkg/models"
	"gorm.io/gorm"
        "github.com/google/uuid"
)

var (
        DB *gorm.DB
)

func TestInit(t *testing.T){
	models.Init()
}
func TestAddUser(t *testing.T){
	id, _ := uuid.Parse("4a767943-426c-4af2-9a81-5d8a766f0877")

	X := &models.RegisterData{ID: id, Name:"test", Email:"test@email.com", Pass:"test"}
		
	result := X.AddUser()

	if result.ID != X.ID {
		t.Errorf("AddUser returned ID %s, expected %s", result.ID, X.ID)
	}

}

