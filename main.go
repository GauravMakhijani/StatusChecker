package main

import (
	"StatusChecker/db"
	"StatusChecker/statuschecker"
	"fmt"
	"net/http"
	"time"
)

func main() {

	//initialize database
	db.Init()
	defer db.DB.Close()

	//initialize store
	dbStore := db.New(db.DB)

	//initialize service
	service := statuschecker.New(dbStore)

	//initialize router
	statuschecker.InitRouter(service)

	//start cron job
	CronJob(service)

	//start server
	fmt.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func CronJob(service statuschecker.Service) {
	ticker := time.NewTicker(time.Minute)
	//update status of websites every minute
	go service.CheckStatus(ticker)
}
