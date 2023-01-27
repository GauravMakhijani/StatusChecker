package service

import (
	"StatusChecker/db"
	"encoding/json"
	"fmt"
	"net/http"
)

func CheckWebsites(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost { //only accept post requests

		input := make(map[string][]string) //map to store json data

		//decode json data
		decode := json.NewDecoder(r.Body)
		if err := decode.Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//loop through websites and get status
		for _, v := range input["websites"] {
			status := GetStatus(v)

			_, err := db.DB.Exec("INSERT INTO links (link, status) VALUES ($1, $2)", v, status)
			if err != nil {
				fmt.Println("Error inserting data into database", err)
				return
			}

		}

		//send response
		msg := "Websites added"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)

		fmt.Println("Post request successful")

	} else if r.Method == http.MethodGet { //only accept get requests

		//get query from url
		query := r.URL.Query().Get("name")

		websites := []db.Website{}

		//check if query is empty
		if query != "" {

			// Search for websites in the database based on the query
			err := db.DB.Select(&websites, "SELECT link,status FROM links WHERE link LIKE $1", "%"+query+"%")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(websites)

		} else {
			// Retrieve all websites from the database
			err := db.DB.Select(&websites, "SELECT link,status FROM links")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(websites)

		}

		fmt.Println("Get request successful")

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
