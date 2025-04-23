package controllers

import (
	"github.com/shayantrix/task_manager_api/pkg/tokens"
	"github.com/shayantrix/task_manager_api/pkg/models"
	"github.com/google/uuid"
	"log"
	"encoding/json"
	"io"
	"fmt"
	//"context"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	//"github.com/shayantrix/task_manager_api/pkg/auth"
	//"github.com/shayantrix/task_manager_api/pkg/models"
	//"log"
	//"github.com/gorilla/mux"
)

func Register(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var reg models.RegisterData
	
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &reg); err != nil{
		log.Fatal("Error in decoding json file", err)
	}
	
	reg.ID = uuid.New()
	
	AllUsers := models.GetAllUsers()
	for _, item := range AllUsers{
		//fmt.Println(item)
		if item.Email == reg.Email{
			http.Error(w, "This email already exists", http.StatusBadRequest)
			return
		} 	
	}

	reg.AddUser()
	
	var h models.HashedPasswords
	var hash_err error
	h.Hashed, hash_err  =  bcrypt.GenerateFromPassword([]byte(reg.Pass), bcrypt.DefaultCost)
	if hash_err != nil{
		log.Fatal("Hashing error: ", hash_err)
	}

	h.ParentRefer = reg.ID

	models.DB.Create(&h)

	fmt.Println(reg.ID)
}

func GetUsers(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	
	AllNames := models.GetUsersNames()
	for _, item := range AllNames{
		json.NewEncoder(w).Encode(item.Name)
	}
}


/*func MakeHandler(fn func(http.ResponseWriter, *http.Request, *RegisterData)) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
	}
}*/

func Login(w http.ResponseWriter, r *http.Request){
	// User should put email and password
	// We will check whether password matches the hashed one that we have in authentication
	//If email does not exist Wont move further.
	w.Header().Set("Content-Type", "application/json")
	
	var reg models.RegisterData

	body, _ := io.ReadAll(r.Body)

	if err := json.Unmarshal(body, &reg); err != nil{
		log.Fatal("Error in Decoding json: ", err)
	}

	found := false

	AllUsers := models.GetAllUsers()

	AllHashedPass := models.GetAllPass()

	for _, item := range AllUsers{
		if item.Email == reg.Email{
			found = true
			for _, hash := range AllHashedPass{
				if err := bcrypt.CompareHashAndPassword(hash.Hashed, []byte(reg.Pass)); err != nil{
					http.Error(w, "Password does not match!", http.StatusNotFound)
					fmt.Println(err)
				}else{
					token, err := tokens.JWTGenerate(item.ID)
					if err != nil{
						log.Fatal("Error in JWT token generation: %s", err)
					}
					json.NewEncoder(w).Encode(token)
					json.NewEncoder(w).Encode(item)
					break
				}
			}
			fmt.Printf("User %s Login seccessfully\n", item.Name)
		}
	}
	if !found{
		http.Error(w, "Email does not exists!", http.StatusNotFound)
	}
}


func Add(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	userIDInterface := r.Context().Value("id")
	
	var X models.Tasks

	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &X); err != nil{
		log.Fatal("Error in recieving user's data: %s", err)
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		log.Fatal("Error in ID type")
	}

	X.ParentRefer = userID

	found := false
	AllTasks := models.GetAllTasks()

	for _, item := range AllTasks{
		if item.ParentRefer == userID {
			X.AddTasks()
			found = true
			json.NewEncoder(w).Encode(item)
			break
			
		}
	}

	if !found{
		X.AddTasks()
		json.NewEncoder(w).Encode(X)
	}
			
}

func Delete(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	userIDInterface := r.Context().Value("id")
	userID, ok := userIDInterface.(uuid.UUID)
        if !ok {
                log.Fatal("Error in ID type")
        }
	
	var X struct{
		DeleteItem string `json:"delete"`
		//AddItem	string	`json:"add"`
	}

	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &X); err != nil{
		log.Fatal("Error in receiving user's data: %s", err)
	}
	
	taskToDelete := models.GetTaskByName(X.DeleteItem)
	taskToDelete.DeleteTaskByName(X.DeleteItem)
	//if err != nil{
	//	io.Write(w, err)
	//	http.Error(w, "Problem in Deleting the task!", http.StatusBadRequest)
	//}
	t := models.GetTaskByRefID(userID)
	
	json.NewEncoder(w).Encode(t)
	
	/*
	for i, item := range TasksData{
		if userID == item.ID{
			for j, task := range item.TasksDatabase{
				if task.TaskString == X.DeleteItem {
					TasksData[i].TasksDatabase = append(item.TasksDatabase[:j], item.TasksDatabase[j+1:]...)
					fmt.Fprintf(w, "%s is deleted from your task list", X.DeleteItem)
					json.NewEncoder(w).Encode(TasksData[i])
					return
				}
			}
			http.Error(w, "Task not found in your task list", http.StatusNotFound)
			return
		}
	}
	*/
	//for i, item := range CompletedTasks{
	//	if userID == item.ID{


	//http.Error(w, "No task found for this user", http.StatusBadRequest)
}
/*
func Update(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	userIDInterface := r.Context().Value("id")
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		log.Fatal("Error in ID type")
	}

	var X struct{
		OldItem string `json:"old"`
		NewItem string `json:"new"`
	}

	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &X); err != nil{
		log.Fatal("Error in receiving user's JSON data: ", err)
	}

	//we iterate and delete the OldItem and add NewItem in Tasks

	if X.NewItem == "" || X.OldItem == ""{
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	for i, item := range TasksData{
		if userID == item.ID{
			for j, task := range item.TasksDatabase{
				if task.TaskString == X.OldItem {
					//TasksData[i].TaskString = append(item.TaskString[:j], item.TaskString[j+1:]...)
					TasksData[i].TasksDatabase[j].TaskString = X.NewItem
					TasksData[i].TasksDatabase[j].Description = "Changed the task"
					TasksData[i].TasksDatabase[j].TaskStatus = false 
					fmt.Fprintf(w, "'%s' is Changed to '%s'", X.OldItem, X.NewItem)
					json.NewEncoder(w).Encode(TasksData[i])
					return
				}else{
					http.Error(w, "Task does not exist", http.StatusNotFound)
				}
			}
		}else{
			http.Error(w, "There is no task for this user", http.StatusNotFound)
		}
	}
}

func Mark(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	
	userIDInterface := r.Context().Value("id")
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		log.Fatal("Error in userID type")
	}

	var X TasksMark

	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &X); err != nil{
		log.Fatal("Error in JSON input file: ", err)
	}
	//X.Description X.Status

	for i, item := range TasksData{
		if userID == item.ID{
			for j, task := range item.TasksDatabase{
				if X.TaskStatus == true && X.TaskString == task.TaskString{
					TasksData[i].TasksDatabase[j].TaskString = X.TaskString
					TasksData[i].TasksDatabase[j].Description = X.Description
					TasksData[i].TasksDatabase[j].TaskStatus = true
				}

			}
			json.NewEncoder(w).Encode(TasksData[i])
		}else{
			http.Error(w, "There is no task for this user", http.StatusBadRequest)
			return
		}
	}

}

func TaskRetrieval(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	userIDInterface := r.Context().Value("id")
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok{
		log.Fatal("Error in userID type")
	}
	params := mux.Vars(r)
	switch X := params["status"]; X{
	case "completed":
		for _, item := range TasksData{
			if userID == item.ID{
				for j, task := range item.TasksDatabase{
					if task.TaskStatus == true{
						json.NewEncoder(w).Encode(item.TasksDatabase[j])
					}
				}
			}else{
				http.Error(w, "There is no task for this user", http.StatusBadRequest)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	case "incomplete":
		for _, item := range TasksData{
                        if userID == item.ID{
                                for j, task := range item.TasksDatabase{
                                        if task.TaskStatus == false{
                                                json.NewEncoder(w).Encode(item.TasksDatabase[j])
                                        }
                                }
                        }else{
                                http.Error(w, "There is no task for this user", http.StatusBadRequest)
                                return
                        }
                }
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
	}

}*/
