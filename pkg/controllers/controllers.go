package controllers

import (
	"github.com/shayantrix/task_manager_api/pkg/tokens"
	"github.com/google/uuid"
	"log"
	"encoding/json"
	"io"
	"fmt"
	//"context"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/shayantrix/task_manager_api/pkg/auth"
	//"github.com/shayantrix/task_manager_api/pkg/models"
	//"log"
	"github.com/gorilla/mux"
)

type RegisterData struct{
        ID uuid.UUID    `json:"id"`
        Name string `json:"name"`
        Email string `json:"email"`
        Pass string `json:"password"`
}

type Tasks struct{
	ID uuid.UUID	`json:"id"`
	//TaskString []string `json:"task"`
	TasksDatabase []TasksMark	`json:"tasks"`
}
//store tasks data
var TasksData []Tasks

type SecureAuth struct{
        Email   string  `json:"email"`
        Pass  []byte    `json:"-"`
}

var secureAuth []SecureAuth


type DataWithoutPass struct{
        Name string `json:"name"`
        Email string `json:"email"`
}

type TasksMark struct{
	TaskString string `json:"task"`
	Description string `json:"description"`
	TaskStatus bool	`json:"completed"`
}

//var CompletedTasks []TasksMark


var Data []RegisterData

var HashedPasswords []byte

//Register handler -> Register(w, r)
func Register(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var reg RegisterData
	
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &reg); err != nil{
		log.Fatal("Error in decoding json file", err)
	}
	
	var hash_err error
	HashedPasswords, hash_err  =  bcrypt.GenerateFromPassword([]byte(reg.Pass), bcrypt.DefaultCost)
	if hash_err != nil{
		log.Fatal("Hashing error: ", hash_err)
	}
	
	reg.ID = uuid.New()

	if Data == nil {
		Data = append(Data, reg)
	}else{
		found := false

		for _, v := range Data{
			if v.Email == reg.Email{
				http.Error(w, "This email already exists", http.StatusBadRequest)
				found = true
			}
		}
		if !found{
                	Data = append(Data, reg)
                }
	}
	fmt.Println(reg.ID)
}

func GetUsers(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var userResponse []DataWithoutPass
	for _, v := range Data{
		userResponse = append(userResponse, DataWithoutPass{
			Name: v.Name,
			Email: v.Email,
		})
	}
	json.NewEncoder(w).Encode(userResponse)
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
	
	var reg RegisterData
	body, _ := io.ReadAll(r.Body)

	if err := json.Unmarshal(body, &reg); err != nil{
		log.Fatal("Error in Decoding json: ", err)
	}

	found := false

	for i, item := range Data{		
		if item.Email == reg.Email{
			found = true
			if err := auth.CheckHashedPassword(reg.Pass, string(HashedPasswords)); err != nil{
				http.Error(w, "Password does not match!", http.StatusNotFound)
				fmt.Println(err)
			}else{
				token, err := tokens.JWTGenerate(item.ID)
       	 			if err != nil {
                			log.Fatal("Error in jwt token generation: %s", err)
       				}
        			json.NewEncoder(w).Encode(token)
				json.NewEncoder(w).Encode(Data[i])
			}
		}
	//fmt.Printf("token: %s", token)
	fmt.Printf("User %v Login secssussfully\n", reg.Name)
	}
	if !found{
		http.Error(w, "Email does not exists", http.StatusNotFound)
	}
}

func Test(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("id")

	response := map[string]interface{}{
		"message": "Xou have accessed a protected route",
		"user_id": userID,
	}
	json.NewEncoder(w).Encode(response)
}

func Add(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	userIDInterface := r.Context().Value("id")
	
	var X TasksMark

	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &X); err != nil{
		log.Fatal("Error in recieving user's data: %s", err)
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		log.Fatal("Error in ID type")
	}

	found := false
	for i, item := range TasksData{
		if item.ID == userID {

			TasksData[i].TasksDatabase = append(TasksData[i].TasksDatabase, X)
			found = true
			json.NewEncoder(w).Encode(TasksData[i].TasksDatabase)
			break
			
		}
	}


	
	if !found{
		NewTask := Tasks{
			ID: userID,
			TasksDatabase: []TasksMark{X},
		}
		TasksData = append(TasksData, NewTask)
		json.NewEncoder(w).Encode(NewTask)
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

	//for i, item := range CompletedTasks{
	//	if userID == item.ID{


	http.Error(w, "No task found for this user", http.StatusBadRequest)
}

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
		http.Error(w, "Get the fuck out", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
	}

}
