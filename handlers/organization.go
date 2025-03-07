package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Org struct {
	Organization Organization `json:"organization"`
	Action       string       `json:"action"`
	Sender       Sender       `json:"sender"`
	Membership   Membership   `json:"membership"`
	Invitation   Invitation   `json:"invitation"`
}

type Invitation struct {
	User    string  `json:"login"`
	Inviter Inviter `json:"inviter"`
}

type Inviter struct {
	User string `json:"login"`
}

type Organization struct {
	Name string `json:"login"`
	Url  string `json:"url"`
}

type Sender struct {
	Username string `json:"login"`
	Type     string `json:"type"`
}

type Membership struct {
	Url   string         `json:"url"`
	State string         `json:"state"`
	User  MembershipUser `json:"user"`
}

type MembershipUser struct {
	Username string `json:"login"`
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
		"organization":            orgResponse.Organization.Name,
		"organization URL":        orgResponse.Organization.Url,
		"action":                  orgResponse.Action,
		"sender":                  orgResponse.Sender.Username,
		"victim":                  orgResponse.Membership.User.Username,
		"victim membership URL":   orgResponse.Membership.Url,
		"victim membership state": orgResponse.Membership.State,

		// for invitation action
		"invite sender":   orgResponse.Invitation.Inviter.User,
		"invite receiver": orgResponse.Invitation.User,
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
