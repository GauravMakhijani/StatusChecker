package db

import (
	"fmt"
	"net/http"
	"strings"
)

type StatusChecker interface {
	CheckStatus(url string) (int, error)
}

type Website struct {
	Link   string `db:"link"`
	Status string `db:"status"`
}

func (ws Website) UpdateStatus(url string) {
	fmt.Println("Checking status of", ws.Link)

	status := GetStatus(ws.Link)

	fmt.Println("Updating status of", ws.Link, "to", status)

	_, err := DB.Exec("UPDATE links SET status = $2 where link = $1", ws.Link, status)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func GetStatus(url string) string {

	//check if url starts with http
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	//get status
	resp, err := http.Get(url)

	if err != nil || resp.StatusCode != 200 {
		return "DOWN"
	}
	return "UP"
}
