package models

import(
	"github.com/google/uuid"
	"fmt"
)

type RegisterData struct{
        ID uuid.UUID    `json:"id"`
        Name string `json:"name"`
        Email string `json:"email"`
        Pass string `json:"password"`
}

type DataWithoutPass struct{
        Name string `json:"name"`
        Email string `json:"email"`
}

var registerData []RegisterData

func hello(){
	fmt.Println("Hello")
}
