package models 

import(
	"github.com/shayantrix/task_manager_api/pkg/config"
	"fmt"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

var (
	DB *gorm.DB
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
	DB = config.Connection()
	if DB == nil{
		panic("DB IS not Initialized")
	}
	err := DB.AutoMigrate(&RegisterData{}, &Tasks{}, &HashedPasswords{})
   	if err != nil {
        	panic("Migration failed: " + err.Error())
    	}
}

func (r *RegisterData) AddUser() *RegisterData{
	result := DB.Create(&r)
	if result.Error != nil{
		fmt.Println("Error in Registering")
	}
	return r
}

func (t *Tasks) AddTasks() *Tasks{
	DB.Select(t.TaskString, t.Description, t.TaskStatus).Create(&t)
	return t
}

func (h *HashedPasswords) StoreHashPasswords() *HashedPasswords{
	DB.Create(&h)
	return h
}

func GetAllPass() []HashedPasswords{
	var result []HashedPasswords
	DB.Distinct("hashed").Find(&result)
	return result
}

func (r *RegisterData) RetrieveUser(id uuid.UUID) *RegisterData{
	DB.Take(&r, "id=?", id)
	return r
}

func GetAllUsers() []RegisterData {
	var result []RegisterData
	DB.Find(&result)
	return result
}

func GetUsersNames() []RegisterData{
	var result []RegisterData
	DB.Distinct("name").Find(&result)
	return result
}

func GetAllTasks() []Tasks{
	var result []Tasks
	DB.Find(&result)
	return result
}

