package main

import (
	"os"
	"log"
	//"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/shayantrix/task_manager_api/pkg/router"
)

func main(){
	r := mux.NewRouter()
	router.RoutingGroup(r)
	http.Handle("/health", r)
	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}
	log.Printf("Server start on port: %s", port)
	
	//Serve the server
	server := &http.Server{
		Addr: ":"+ port,
		Handler: r,
	}
	//go func(){
		if err := server.ListenAndServe(); err != nil{
			log.Fatal(err)
		}
	//}()
	/*
	// Server Shutdown
	quit := make(chan os.Signal, 1)
	*/
}


