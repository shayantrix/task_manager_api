package router

import(
	"net/http"
	"github.com/gorilla/mux"
	"github.com/shayantrix/task_manager_api/pkg/middleware"
	"github.com/shayantrix/task_manager_api/pkg/controllers"
)

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

