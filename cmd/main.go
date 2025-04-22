package main

import (
	//"context"
	//"os/signal"
	//"syscall"
	//"time"
	"os"
	"log"
	//"fmt"
	"net/http"
	"github.com/gorilla/mux"
	//"github.com/shayantrix/task_manager_api/pkg/config"
	"github.com/shayantrix/task_manager_api/pkg/models"
	"github.com/shayantrix/task_manager_api/pkg/router"
)

func main(){
	//db := config.Connection()
	models.Init()
	
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

	if err := server.ListenAndServe(); err != nil{
                        log.Fatal(err)
        }

	/*
	go func(){
		if err := server.ListenAndServe(); err != nil{
			log.Fatal(err)
		}
	}()
	
	// Server Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit //wait here until interrupt

	//start shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil{
		log.Fatal("Forced shutdown: ", err)
	}
	log.Println("Server stopped gracefully")
	*/
}


