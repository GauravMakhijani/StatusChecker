package main

import (
	"StatusChecker/db"
	"StatusChecker/statuschecker"
	"net/http"
	"time"

	logger "github.com/sirupsen/logrus"
)

// HandleWebsites returns a handler function that calls the service's
func InitRouter(service statuschecker.Service) {
	http.HandleFunc("/websites", statuschecker.HandleWebsites(service))
}

func main() {

	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})

	//initialize database
	db.Init()
	defer db.DB.Close()

	//initialize store
	dbStore := db.New(db.DB)

	//initialize service
	service := statuschecker.New(dbStore)

	//initialize router
	InitRouter(service)
	//start cron job
	CronJob(service)

	//start server
	logger.Info("Server starting on http://localhost:8080") //.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func CronJob(service statuschecker.Service) {
	ticker := time.NewTicker(time.Minute)
	//update status of websites every minute
	go service.CheckStatus(ticker)
}
