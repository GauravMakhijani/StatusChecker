package statuschecker

import (
	"StatusChecker/db"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	Append(ctx context.Context, website db.WebsiteStatus) (err error)
	// Check(ctx context.Context, name string) (status bool, err error)
	GetSimilar(ctx context.Context, name string) (similar []db.WebsiteStatus, err error)
	CheckStatus(*time.Ticker)
	GetStatus(ctx context.Context, url string) (status string)
	GetAll(ctx context.Context) (ws []db.WebsiteStatus, err error)
}

type statusService struct {
	repo db.StatusStorer
}

func New(repo db.StatusStorer) Service {
	return &statusService{repo: repo}
}

func (s *statusService) GetSimilar(ctx context.Context, query string) (similar []db.WebsiteStatus, err error) {

	//Make a call to the database to get the status of the website
	similar, err = s.repo.GetWebsiteStatus(query)
	return
}

func (s *statusService) Append(ctx context.Context, website db.WebsiteStatus) (err error) {
	s.repo.InsertWebsite(website)
	//Make a call to the database to get the status of the website
	// s.repo.GetWebsiteStatus(url)
	// //return the status of the website

	// for {
	// 	select {
	// 	//wait for ticker to send a message
	// 	case <-ticker.C:

	// 		websites := []db.Website{}
	// 		err := db.DB.Select(&websites, "SELECT link,status FROM links")
	// 		if err != nil {
	// 			fmt.Println("Error:", err)
	// 			continue
	// 		}

	// 		//update status of each website
	// 		for _, website := range websites {

	// 			//send website to updateStatus function in a goroutine
	// 			go func(website db.Website) {
	// 				website.UpdateStatus(website.Link)
	// 			}(website)

	// 		}

	// 	}

	// }
	return
}

func (s *statusService) CheckStatus(ticker *time.Ticker) {
	//return the status of the website

	for {
		select {
		//wait for ticker to send a message
		case <-ticker.C:

			fmt.Println("Checking status of websites")
			websites, err := s.repo.GetAll()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			//update status of each website
			for _, website := range websites {

				//send website to updateStatus function in a goroutine
				go func(website db.WebsiteStatus) {
					status := s.GetStatus(context.Background(), website.Link)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}
					s.repo.UpdateWebsiteStatus(website.Link, status)
					fmt.Println("Status of", website.Link, "is updated to", status)
				}(website)

			}

		}
	}
}

func (s *statusService) GetStatus(ctx context.Context, url string) (status string) {
	// check if url starts with http
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	//get status
	fmt.Println("Checking status of", url)
	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != 200 {
		return "DOWN"
	}
	return "UP"
}

func (s *statusService) GetAll(ctx context.Context) (ws []db.WebsiteStatus, err error) {
	ws, err = s.repo.GetAll()
	return
}
