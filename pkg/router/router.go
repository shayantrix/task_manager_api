package router

import(
	"github.com/gorilla/mux"
	"github.com/shayantrix/task_manager_api/pkg/controllers"
)

var RoutingGroup = func(router *mux.Router){
	router.HandleFunc("/health", controllers.Test).Methods("GET")
	router.HandleFunc("/register", controllers.MakeHandler(controllers.Register)).Methods("POST")
}

