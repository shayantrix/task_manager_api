package controllers

import (
	"log"
	"encoding/json"
	"io"
	"net/http"
	//"log"
	//"github.com/gorilla/mux"
)

type RegisterData struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Pass string `json:"password"`
}

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
func Register(w http.ResponseWriter, r *http.Request, reg *RegisterData){
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	
	var registerData RegisterData

	//params := mux.Vars(r)
	if err := json.NewDecoder(r.Body).Decode(reg); err != nil{
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(registerData)
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, *RegisterData)) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
	}
}
