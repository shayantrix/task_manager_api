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
	//"github.com/gorilla/mux"
)

type RegisterData struct{
        ID uuid.UUID    `json:"id"`
        Name string `json:"name"`
        Email string `json:"email"`
        Pass string `json:"password"`
}

type Tasks struct{
	ID uuid.UUID	`json:"id"`
	TaskString []string `json:"task"`
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

		for _, v := range Data{
			if v.Email == reg.Email{
				fmt.Printf("This email already exists")
			
			}else{
				Data = append(Data, reg)
	
			}
		}
	}
	fmt.Println(reg.ID)
	/*
	token, err := tokens.JWTGenerate(reg.ID)
        if err != nil {
                log.Fatal("Error in jwt token generation: %s", err)
        }
        json.NewEncoder(w).Encode(token)
	//Data = append(Data, reg)
	*/
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
				log.Fatal("Password does not match! ", err)
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
	fmt.Printf("User %v Login secssussfully", reg.Name)
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
	
	var X struct{
                //DeleteItem string `json:"delete"`
                AddItem  []string  `json:"add"`
        }
	
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
			TasksData[i].TaskString = append(TasksData[i].TaskString, X.AddItem...)
			found = true
			json.NewEncoder(w).Encode(TasksData[i])
			break
		}
	}
	
	if !found{
		NewTask := Tasks{
			ID: userID,
			TaskString: X.AddItem,
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
			for j, task := range item.TaskString{
				if task == X.DeleteItem {
					TasksData[i].TaskString = append(item.TaskString[:j], item.TaskString[j+1:]...)
					fmt.Fprintf(w, "%s is deleted from your task list", X.DeleteItem)
					json.NewEncoder(w).Encode(TasksData[i])
					return
				}
			}
			http.Error(w, "Task not found in your task list", http.StatusNotFound)
			return
		}
	}
	http.Error(w, "No task found for this user", http.StatusBadRequest)
}

