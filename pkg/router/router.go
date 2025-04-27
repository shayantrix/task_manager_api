package router

import(
	//"net/http"
	//"github.com/gorilla/mux"
	"github.com/shayantrix/task_manager_api/pkg/middleware"
	"fmt"
	"log"
	"github.com/shayantrix/task_manager_api/pkg/controllers"
	"github.com/gin-gonic/gin"
)


//gin rout
var (
	Router = gin.Default()
	

		UnProtected = Router.Group("/unprotected")
		//UnProtected.POST("/register", controllers.Register)
		//UnProtected.POST("/users", controllers.)
	
		Protected = Router.Group("/protected")
)

func RunRouter(){
	Router.Use(gin.Recovery())
	UnProtected.POST("/register", controllers.Register)
	UnProtected.GET("/users", controllers.GetUsers)
	UnProtected.POST("/login", controllers.Login)

	Protected.Use(middleware.Authorization())
	{
		Protected.POST("/add", controllers.Add)
		Protected.POST("/delete", controllers.Delete)
		Protected.PATCH("/update", controllers.Update)
		Protected.POST("/mark", controllers.Mark)
		Protected.GET("/tasks/:status", controllers.TaskRetrieval)
	}	
	fmt.Println("Server running on localhost:8080")
   	if err := Router.Run(":8080"); err != nil {
        	log.Fatalf("Failed to start server: %v", err)
    	}
}

/*
var RoutingGroup = func(router *mux.Router){
	//router.HandleFunc("/health", controllers.Test).Methods("GET")
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.Handle("/protected/taskmanagement/add", middleware.Authorization(http.HandlerFunc(controllers.Add))).Methods("POST")
	router.Handle("/protected/taskmanagement/delete", middleware.Authorization(http.HandlerFunc(controllers.Delete))).Methods("POST")
	//router.Handle("/protected/taskmanagement/tasks/{id}", middleware.Authorization(http.HandlerFunc(controllers.GetAllTasks))).Methods("GET")
	router.Handle("/protected/taskmanagement/update", middleware.Authorization(http.HandlerFunc(controllers.Update))).Methods("PATCH")
	router.Handle("/protected/taskmanagement/mark", middleware.Authorization(http.HandlerFunc(controllers.Mark))).Methods("POST")
	router.Handle("/protected/taskmanagement/tasks/{status}", middleware.Authorization(http.HandlerFunc(controllers.TaskRetrieval))).Methods("GET")
}
*/
