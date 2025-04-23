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
		ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"json:"id"`
                Name string `json:"name"`
                Email string `json:"email"`
		Pass	string	`json:"-"`
		TasksData []Tasks	`gorm:"foreignKey:ParentRefer"`
	}

	Tasks struct{
		ID uint    `gorm:"primaryKey;autoIncrement"json:"id"`
		TaskString string `json:"task"`
        	Description string `json:"description"`
        	TaskStatus bool `json:"completed"`
		ParentRefer uuid.UUID	//foreign key: used in RegisterData
		User	RegisterData	`gorm:"foreignKey:ParentRefer;references:ID"`
	}

	HashedPasswords struct{
		ID uint	`gorm:"primaryKey"json:"id"`
		Hashed	[]byte	`json:"-"`
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
	DB.Create(&t)
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

func GetTaskByName(x string) *Tasks{
	var t *Tasks
	DB.First(&t, "task_string", x)
	return t
}

func (t *Tasks) DeleteTaskByName(x string){
	DB.Where("task_string", x).Delete(&t)
}

func GetTaskByRefID(ID uuid.UUID) *Tasks{
	var t *Tasks
	DB.First(&t, "parent_refer = ?", ID)
	return t
}

func UpdateTask(ID uuid.UUID,delet , add string) *gorm.DB{
	var t *Tasks
	result := DB.Model(&t).Where("parent_refer = ?", ID).Where("task_string", delet).Update("task_string", add)
	return result
}
