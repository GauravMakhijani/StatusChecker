package statuschecker

import (
	"StatusChecker/db"
	"encoding/json"
	"fmt"
	"net/http"
)

func handlePostRequest(w http.ResponseWriter, r *http.Request, service Service) {
	input := make(map[string][]string) //map to store json data

	//decode json data
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//loop through websites and get status
	for _, url := range input["websites"] {
		status := service.GetStatus(r.Context(), url)

		//insert website and status into database
		err := service.Append(r.Context(), db.WebsiteStatus{Link: url, Status: status})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			continue
		}
	}

	//send response
	msg := "Websites added"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)

	fmt.Println("Post request successful")
}

func handleGetRequest(w http.ResponseWriter, r *http.Request, service Service) {

	//get query from url
	query := r.URL.Query().Get("name")

	//check if query is empty
	if query != "" {

		// Search for websites in the database based on the query

		websites, err := service.GetSimilar(r.Context(), query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(websites)

	} else {
		// Retrieve all websites from the database

		websites, err := service.GetAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(websites)

	}

	fmt.Println("Get request successful")
}

// HandleWebsites handles all requests to the /websites endpoint
func HandleWebsites(service Service) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost { //only accept post requests

			handlePostRequest(w, r, service)

		} else if r.Method == http.MethodGet { //only accept get requests

			handleGetRequest(w, r, service)

		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}

}
