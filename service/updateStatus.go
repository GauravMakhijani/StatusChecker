package service

import (
	"StatusChecker/db"
	"fmt"
	"time"
)

func CheckStatus(ticker *time.Ticker) {
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

			//update status of each website
			for _, website := range websites {

				//send website to updateStatus function in a goroutine
				go func(website db.Website) {
					website.UpdateStatus(website.Link)
				}(website)

			}

		}

	}
}
