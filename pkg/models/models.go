package models 

import(
	"github.com/shayantrix/task_manager_api/pkg/database"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

var (
	db *gorm.DB
)

type(
	RegisterData struct{
		gorm.Model
		ID uuid.UUID `gorm:"primaryKey"json:"id"`
                Name string `json:"name"`
                Email string `json:"email"`
		Pass	string	`json:"-"`
		TasksData []Tasks	`gorm:"foreignKey:ParentRefer"`
	}

	Tasks struct{
		gorm.Model
		ID uuid.UUID    `gorm:"primaryKey"json:"id"`
		TaskString string `json:"task"`
        	Description string `json:"description"`
        	TaskStatus bool `json:"completed"`
		ParentRefer uuid.UUID	//foreign key: used in RegisterData
		User	RegisterData	`gorm:"foreignKey:ParentRefer"`
	}

	HashedPasswords struct{
		gorm.Model
		ID uint	`gorm:"primaryKey"json:"id"`
		Hashed	[]byte	`gorm:""json:"-"`
		ParentRefer uuid.UUID	//foreign key: used in RegisterData
		User	RegisterData	`gorm:"foreignKey:ParentRefer"`
	}
)

func Init(){
	db = config.Connection()
	db.AutoMigrate(&RegisterData{})
	db.AutoMigrate(&Tasks{})
	db.AutoMigrate(&HashedPasswords{})
}

func (r *RegisterData) AddUser() *RegisterData{
	db.Create(&r)
	return r
}

func (h *HashedPasswords) StoreHashPasswords() *HashedPasswords{
	db.Create(&h)
	return h
}

func (r *RegisterData) RetrieveUser(id uuid.UUID) *RegisterData{
	db.Take(&r, "id=?", id)
	return r
}

func GetAllUsers() []RegisterData {
	var result []RegisterData
	db.Table("users").Take(&result)
	return result
}
