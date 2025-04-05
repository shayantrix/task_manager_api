package controllers

import (
	"io"
	"net/http"
	//"log"
)

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
