package controllers

import (
	"github.com/shayantrix/task_manager_api/pkg/tokens"
	"github.com/shayantrix/task_manager_api/pkg/models"
	"github.com/google/uuid"
	"log"
	//"encoding/json"
	//"io"
	"fmt"
	//"context"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	//"github.com/shayantrix/task_manager_api/pkg/auth"
	//"github.com/shayantrix/task_manager_api/pkg/models"
	//"log"
	"github.com/gin-gonic/gin"
	//"github.com/gorilla/mux"
)

func Register(c *gin.Context){
	var reg models.RegisterData
	//Json.Unmarshal
	if err := c.ShouldBindJSON(&reg); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Input for user"})
		return
	}

	reg.ID = uuid.New()

	AllUsers := models.GetAllUsers()
        for _, item := range AllUsers{
                //fmt.Println(item)
                if item.Email == reg.Email{
			c.JSON(http.StatusBadRequest, gin.H{"Error": "This email already exist"})
                        return
                }       
        }

        reg.AddUser()
	// Error handling for adding user should be added
        var h models.HashedPasswords
        var hash_err error
        h.Hashed, hash_err  =  bcrypt.GenerateFromPassword([]byte(reg.Pass), bcrypt.DefaultCost)
        if hash_err != nil{
                log.Fatal("Hashing error: ", hash_err)
        }

        h.ParentRefer = reg.ID

        models.DB.Create(&h)

	c.JSON(http.StatusCreated, reg)

        fmt.Println(reg.ID)
}




/*
This part is just the same but implemented with http
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

*/

func GetUsers(c *gin.Context){
	AllNames := models.GetUsersNames()
        for _, item := range AllNames{
        	c.JSON(http.StatusOK, item.Name)
	}
}
/*
func GetUsers(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	
	AllNames := models.GetUsersNames()
	for _, item := range AllNames{
		json.NewEncoder(w).Encode(item.Name)
	}
}
*/

/*func MakeHandler(fn func(http.ResponseWriter, *http.Request, *RegisterData)) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
	}
}*/

func Login(c *gin.Context){
	var reg models.RegisterData
	
	if err := c.ShouldBindJSON(&reg); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid user input"})
	}

	found := false

	AllUsers := models.GetAllUsers()

        AllHashedPass := models.GetAllPass()
	// I have to change this to reduce the loop
        for _, item := range AllUsers{
                if item.Email == reg.Email{
                        found = true
                        for _, hash := range AllHashedPass{
                                if err := bcrypt.CompareHashAndPassword(hash.Hashed, []byte(reg.Pass)); err != nil{
                                        //http.Error(w, "Password does not match!", http.StatusNotFound)
					c.JSON(http.StatusBadRequest, gin.H{"Error":"Password does not match!"})
					fmt.Println(err)
                                }else{
                                        token, err := tokens.JWTGenerate(item.ID)
                                        if err != nil{
                                                log.Fatal("Error in JWT token generation: %s", err)
                                        }
                                        //json.NewEncoder(w).Encode(token)
                                        //json.NewEncoder(w).Encode(item)
					c.JSON(http.StatusOK, token)
					c.JSON(http.StatusOK, item)
                                        break
                                }
                        }
                        fmt.Printf("User %s Login seccessfully\n", item.Name)
                }
        }
        if !found{
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Email does not exists!"})
                //http.Error(w, "Email does not exists!", http.StatusNotFound)
        }
}
/*
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
*/

func Add(c *gin.Context){
        userIDInterface, exists := c.Get("id")
	if !exists {
        	c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
        	return
    	}

	userID, ok := userIDInterface.(uuid.UUID)
        if !ok {
                log.Fatal("Error in ID type")
        }

	var X models.Tasks
	
	if err := c.ShouldBindJSON(&X); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Error in receiving user's data"})
		return
	}
	
	X.ParentRefer = userID

        found := false
        AllTasks := models.GetAllTasks()

        for _, item := range AllTasks{
                if item.ParentRefer == userID {
                        X.AddTasks()
                        found = true
                        //c.JSON(http.StatusCreated, item)
			c.JSON(http.StatusCreated, gin.H{
        			"message": "Task added successfully",
       				 "task":    item, })
                        break

                }
        }
	if !found{
                X.AddTasks()
                c.JSON(http.StatusCreated, X)
        }
}
	

/*
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
*/

func Delete(c *gin.Context){
	userIDInterface, exists := c.Get("id")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
                return
        }

        userID, ok := userIDInterface.(uuid.UUID)
        if !ok {
                log.Fatal("Error in ID type")
        }
	
	var X struct{
                DeleteItem string `json:"delete"`
        }
	
	if err := c.ShouldBindJSON(&X); err != nil{
                c.JSON(http.StatusBadRequest, gin.H{"Error": "Error in receiving user's data"})
                return
        }
	// DB interactions
	taskToDelete := models.GetTaskByName(X.DeleteItem)
        taskToDelete.DeleteTaskByName(X.DeleteItem)
        t := models.GetTaskByRefID(userID)
	
	c.JSON(http.StatusOK, t)
}

/*
func Delete(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	userIDInterface := r.Context().Value("id")
	userID, ok := userIDInterface.(uuid.UUID)
        if !ok {
                log.Fatal("Error in ID type")
        }
	
	var X struct{
		DeleteItem string `json:"delete"`
	}

	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &X); err != nil{
		log.Fatal("Error in receiving user's data: %s", err)
	}
	
	taskToDelete := models.GetTaskByName(X.DeleteItem)
	taskToDelete.DeleteTaskByName(X.DeleteItem)
	t := models.GetTaskByRefID(userID)
	
	json.NewEncoder(w).Encode(t)
	
}

*/

func Update(c *gin.Context){
	userIDInterface, exists := c.Get("id")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
                return
        }

        userID, ok := userIDInterface.(uuid.UUID)
        if !ok {
                log.Fatal("Error in ID type")
        }
	
	var X struct{
                OldItem string `json:"old"`
                NewItem string `json:"new"`
        }

	if err := c.ShouldBindJSON(&X); err != nil{
                c.JSON(http.StatusBadRequest, gin.H{"Error": "Error in receiving user's data"})
                return
        }

	//DB interactions
	if X.NewItem == "" || X.OldItem == ""{
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input"})
                return
        }

        taskToUpdate := models.UpdateTask(userID, X.OldItem, X.NewItem)
        if taskToUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Error in Updating the task"})
		return
	}

        t := models.GetTaskByRefID(userID)

	c.JSON(http.StatusOK, t)
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
	
	taskToUpdate := models.UpdateTask(userID, X.OldItem, X.NewItem)
	if taskToUpdate.Error != nil {
		http.Error(w, "Error in updating the task!", http.StatusBadRequest)
		return
	}

	t := models.GetTaskByRefID(userID)

	json.NewEncoder(w).Encode(t)
}

*/

func Mark(c *gin.Context){
	userIDInterface, exists := c.Get("id")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
                return
        }

        userID, ok := userIDInterface.(uuid.UUID)
        if !ok {
                log.Fatal("Error in ID type")
        }
	
	var X models.Tasks

	if err := c.ShouldBindJSON(&X); err != nil{
                c.JSON(http.StatusBadRequest, gin.H{"Error": "Error in receiving user's data"})
                return
        }

	//DB interactions
	taskMarked := models.UpdateTaskMark(userID, X.TaskString, X.Description)
        if taskMarked.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error in marking the task"})
                return
        }

        t := models.GetMarkedTasks(userID)
	
	c.JSON(http.StatusOK, t)
}
/*
func Mark(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	
	userIDInterface := r.Context().Value("id")
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		log.Fatal("Error in userID type")
	}

	var X models.Tasks

	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &X); err != nil{
		log.Fatal("Error in JSON input file: ", err)
	}
	//X.Description X.Status
	
	//UpdateTask(ID uuid.UUID, name string, description string)
	taskMarked := models.UpdateTaskMark(userID, X.TaskString, X.Description)
	if taskMarked.Error != nil{
		http.Error(w, "Error in marking the task!", http.StatusBadRequest)
		return
	}
	
	t := models.GetMarkedTasks(userID)

	json.NewEncoder(w).Encode(t)

}
*/

func TaskRetrieval(c *gin.Context){
	userIDInterface, exists := c.Get("id")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
                return
        }

        userID, ok := userIDInterface.(uuid.UUID)
        if !ok {
                log.Fatal("Error in ID type")
        }

	params := c.Param("status")

        switch params{
        case "completed":
                t := models.GetMarkedTasks(userID)
                c.JSON(http.StatusOK, t)
        case "incomplete":
                incompleteTasks := models.GetIncompleteTasks(userID)
		c.JSON(http.StatusOK, incompleteTasks)
        default:
		c.JSON(http.StatusBadRequest, gin.H{"Error":"The parameters are not true"})
        }
}
/*
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
		t := models.GetMarkedTasks(userID)
		json.NewEncoder(w).Encode(t)
		w.WriteHeader(http.StatusOK)
	case "incomplete":
		incompleteTasks := models.GetIncompleteTasks(userID)
		json.NewEncoder(w).Encode(incompleteTasks)
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
	}

}*/
