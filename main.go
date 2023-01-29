package main

import (
	"StatusChecker/db"
	"StatusChecker/service"
	"fmt"
	"net/http"
	"time"
)

func main() {

	//initialize database
	db.Init()
	defer db.DB.Close()

	//handle requests
	http.HandleFunc("/websites", service.CheckWebsites)

	//start ticker
	ticker := time.NewTicker(time.Second * 10)
	//update status of websites every minute
	go service.CheckStatus(ticker)

	//start server
	fmt.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
