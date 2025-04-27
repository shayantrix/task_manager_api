package controllers

import (
	"testing"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"net/http"
	"github.com/shayantrix/task_manager_api/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func TestRegister (t *testing.T){
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	// string, HandlerFunc
	r.POST("/register", gin.HandlerFunc(controllers.Register))
	
	//Fake registeration json

	type Req struct{
		Name string	`json:"name"`
		Email string	`json:"email"`
		Pass string	`json:"password"`
	}

	var reqbody = Req{Name: "testUser", Email: "test@mail.com", Pass: "123dasdjja"}

	requestBody, err := json.Marshal(reqbody)
	if err != nil {
		t.Fatalf("Error in handling json input(json.Marshal): %v", err)
	}

	//create a fake http request
	// req == http.Request
	req, err := http.NewRequest( http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Error in creating fake http request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	//Response Writer (Record the Response)
	w := httptest.NewRecorder()

	//Send the request
	//(Engin *Engin) ServeHTTP (w http.ResponseWriter, r *http.Request) 
	r.ServeHTTP(w, req)

	//Check if the status code is OK

	if w.Code != http.StatusCreated{
		t.Errorf("Expected the 201 status code, got %d", w.Code)
	}

}


		

