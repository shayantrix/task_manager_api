package controllers

import (
	"github.com/shayantrix/task_manager_api/pkg/tokens"
	"github.com/google/uuid"
	"log"
	"encoding/json"
	"io"
	"fmt"
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

type SecureAuth struct{
        Email   string  `json:"email"`
        Pass  []byte    `json:"-"`
}

var secureAuth []SecureAuth


type DataWithoutPass struct{
        Name string `json:"name"`
        Email string `json:"email"`
}

var registerData []RegisterData

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

	if registerData == nil {
		registerData = append(registerData, reg)
	}else{

		for _, v := range registerData{
			if v.Email == reg.Email{
				fmt.Printf("This email already exists")
			
			}else{
				registerData = append(registerData, reg)
	
			}
		}
	}
	fmt.Println(reg.ID)
	token, err := tokens.JWTGenerate(reg.ID)
        if err != nil {
                log.Fatal("Error in jwt token generation: %s", err)
        }
        json.NewEncoder(w).Encode(token)
	//registerData = append(registerData, reg)
}

func GetUsers(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var userResponse []DataWithoutPass
	for _, v := range registerData{
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

	for i, item := range registerData{		
		// Use jwt token to login
		userID := r.Context().Value("userID").(string)
		if item.ID.String() == userID{
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"userID": "%s"}`, userID)))
			json.NewEncoder(w).Encode(registerData[i])
			
		}

		if item.Email != reg.Email{
			log.Fatal("Email does not exists")
		}

		if err := auth.CheckHashedPassword(reg.Pass, string(HashedPasswords)); err != nil{
			log.Fatal("Password does not match! ", err)
		}else{
			json.NewEncoder(w).Encode(registerData[i])
		}
	}
	/*
	token, err := tokens.JWTGenerate(reg.ID)
	if err != nil {
		log.Fatal("Error in jwt token generation: %s", err)
	}
	json.NewEncoder(w).Encode(token)
	*/
	//fmt.Printf("token: %s", token)
	fmt.Printf("User %s Login secssussfully", reg.Name)
}
