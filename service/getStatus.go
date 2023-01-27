package service

import (
	"net/http"
	"strings"
)

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
