package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"orgguard/db"

	_ "github.com/lib/pq"
)

// Organization event payload structure
type Org struct {
	Action       string       `json:"action"`
	Organization Organization `json:"organization"`
	Sender       Sender       `json:"sender"`
	Invitation   Invitation   `json:"invitation,omitempty"`
	Membership   Membership   `json:"membership,omitempty"`
}

// Organization details
type Organization struct {
	Name string `json:"login"`
}

// Sender details
type Sender struct {
	Username string `json:"login"`
}

// Invitation details (for `member_invited`)
type Invitation struct {
	User string `json:"login"`
}

// Membership details (for `member_added`)
type Membership struct {
	User  MembershipUser `json:"user"`
	State string         `json:"state"`
	Url   string         `json:"html_url"`
}

// Membership user details
type MembershipUser struct {
	Username string `json:"login"`
}

// Validate target user with OPA
func validateWithOPA(targetUser string) (bool, error) {
	opaURL := "http://localhost:8181/v1/data/orgguard/allow"

	payload := map[string]interface{}{
		"input": map[string]string{
			"target_user": targetUser,
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	resp, err := http.Post(opaURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return false, fmt.Errorf("OPA request failed: %v", err)
	}
	defer resp.Body.Close()

	var opaResp map[string]bool
	err = json.NewDecoder(resp.Body).Decode(&opaResp)
	if err != nil {
		return false, fmt.Errorf("Failed to parse OPA response: %v", err)
	}

	return opaResp["result"], nil
}

// Handle incoming GitHub organization events
func OrganizationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	rawPayload := r.FormValue("payload")
	decodedPayload, err := url.QueryUnescape(rawPayload)
	if err != nil {
		http.Error(w, "Failed to decode payload", http.StatusBadRequest)
		return
	}

	var orgResponse Org
	err = json.Unmarshal([]byte(decodedPayload), &orgResponse)
	if err != nil {
		log.Printf("Error unmarshalling payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	var targetUser string
	if orgResponse.Action == "member_invited" {
		targetUser = orgResponse.Invitation.User
	} else if orgResponse.Action == "member_added" {
		targetUser = orgResponse.Membership.User.Username
	}

	// allowed, err := validateWithOPA(targetUser)
	// if err != nil {
	// 	log.Printf("OPA validation error: %v", err)
	// 	http.Error(w, "OPA validation failed", http.StatusInternalServerError)
	// 	return
	// }

	// opaResult := "passed"
	// if !allowed {
	// 	opaResult = "failed"
	// }

	allowed := true

	opaResult := "true"

	err = db.SaveLog(
		orgResponse.Action,
		orgResponse.Organization.Name,
		orgResponse.Sender.Username,
		targetUser,
		orgResponse.Membership.State,
		orgResponse.Membership.Url,
		opaResult,
	)

	if err != nil {
		log.Printf("Failed to log event: %v", err)
	}

	if !allowed {
		http.Error(w, "Policy violation: Usernames must start with 'jmd'", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Event processed successfully"}`))
}
