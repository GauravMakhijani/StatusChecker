package statuschecker

import (
	"StatusChecker/db"
	"encoding/json"

	"net/http"

	"github.com/sirupsen/logrus"
)

func CreateWebsiteHandler(w http.ResponseWriter, r *http.Request, service Service) {
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
		err := service.Add(r.Context(), db.WebsiteStatus{Link: url, Status: status})
		logrus.Info("when error occured : ", err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	//send response
	msg := "Websites added"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
	logrus.Info("Post request successful")
}

func GetWebsiteHandler(w http.ResponseWriter, r *http.Request, service Service) {

	//get query from url
	websiteName := r.URL.Query().Get("name")

	//check if query is empty
	if websiteName != "" {

		// Search for websites in the database based on the query

		websites, err := service.GetSimilar(r.Context(), websiteName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(websites)
		return

	}

	// Retrieve all websites from the database

	websites, err := service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(websites)

	logrus.Info("Get request successful")
}

// HandleWebsites handles all requests to the /websites endpoint
func HandleWebsites(service Service) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost { //only accept post requests

			CreateWebsiteHandler(w, r, service)

		} else if r.Method == http.MethodGet { //only accept get requests

			GetWebsiteHandler(w, r, service)

		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}

}
