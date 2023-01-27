package main

import (
	"StatusChecker/db"
	"fmt"
	"net/http"
)

func main() {

	//initialize database
	db.Init()

	//handle requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})

	//start server
	fmt.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
