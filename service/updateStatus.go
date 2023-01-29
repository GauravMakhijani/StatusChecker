package service

import (
	"StatusChecker/db"
	"fmt"
	"time"
)

func update(website db.Website) {
	fmt.Println("Checking status of", website.Link)

	status := GetStatus(website.Link)

	fmt.Println("Updating status of", website.Link, "to", status)

	_, err := db.DB.Exec("UPDATE links SET status = $2 where link = $1", website.Link, status)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func UpdateStatus(ticker *time.Ticker) {
	for {
		select {
		//wait for ticker to send a message
		case <-ticker.C:

			websites := []db.Website{}
			err := db.DB.Select(&websites, "SELECT link,status FROM links")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			for _, website := range websites {
				fmt.Println("sent to update", website.Link)
				go update(website)

			}
		default:
			fmt.Println("Waiting for ticker")
		}

	}
}
