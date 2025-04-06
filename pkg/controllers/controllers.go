package controllers

import (
	"log"
	"encoding/json"
	"io"
	"fmt"
	"net/http"
	//"github.com/shayantrix/task_manager_api/pkg/models"
	//"log"
	//"github.com/gorilla/mux"
)

type RegisterData struct{
        Name string `json:"name"`
        Email string `json:"email"`
        Pass string `json:"password"`
}

type DataWithoutPass struct{
	Name string `json:"name"`
	Email string `json:"email"`
}

var registerData []RegisterData	

func Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// type Status struct{
	// 	statusCode string `json:"status"`
	// }

	// s := Status{statusCode: "ok"}

	// json.NewEncoder(w).Encode(s)
	io.WriteString(w, "Hello World")
}

//Register handler -> Register(w, r)
func Register(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var reg RegisterData
	
	/*if err := json.NewDecoder(r.Body).Decode(&reg); err != nil{
		log.Fatal(err)
	}*/
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &reg); err != nil{
		log.Fatal("Error in decoding json file", err)
	}

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
