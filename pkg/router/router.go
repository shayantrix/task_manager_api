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
	//protected := middleware.Authentication()
	//router.HandleFunc("/protected", protected).Methods("POST")
	router.Handle("/protected", middleware.Authentication(http.HandlerFunc(controllers.Test))).Methods("GET")
	//router.Handle("/protected", middleware.Authentication(controllers.Test))
}
