package controllers

import (
	"log"
	"encoding/json"
	"io"
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
	//params := mux.Vars(r)
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil{
		log.Fatal(err)
	}

	registerData = append(registerData, reg)

	json.NewEncoder(w).Encode(registerData)

}

func GetUsers(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	/*params := mux.Vars(r)
	name, err := strconv.ParseInt(params["name"], 10, 64)
	if err != nil {
		log.Fatal(err)
	}*/
	json.NewEncoder(w).Encode(registerData)
}


/*func MakeHandler(fn func(http.ResponseWriter, *http.Request, *RegisterData)) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
	}
}*/
