package statuschecker

import "net/http"

// HandleWebsites returns a handler function that calls the service's
func InitRouter(service Service) {
	http.HandleFunc("/websites", HandleWebsites(service))
}
