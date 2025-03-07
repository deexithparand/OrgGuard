package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Org struct {
	Action string `json:"action"`
}

// orgHandler processes the webhook request
func OrganizartionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Decode the `payload` form field
	rawPayload := r.FormValue("payload")
	decodedPayload, err := url.QueryUnescape(rawPayload)
	if err != nil {
		http.Error(w, "Failed to decode payload", http.StatusBadRequest)
		return
	}

	var orgResponse Org
	err = json.Unmarshal([]byte(decodedPayload), &orgResponse)
	if err != nil {
		log.Printf("error unmarshalling the payload : %v", err)
		return
	}

	fmt.Println("Action Performed : ", orgResponse.Action)

	response := map[string]interface{}{
		"action":  orgResponse.Action,
		"message": "found",
	}

	responseStr, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshalling the response : %v", err)
		return
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(string(responseStr))
}
