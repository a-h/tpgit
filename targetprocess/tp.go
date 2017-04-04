package targetprocess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// API is a TargetProcess API endpoint.
type API struct {
	// URL of the API endpoint of TargetProcess.
	URL           string
	Authenticator func(*http.Request)
}

// Authenticator represents a way of authenticating against TargetProcess,
// e.g. PasswordAuth or TokenAuth.
type Authenticator func(*http.Request)

// NewAPI creates a new TargetProcess endpoint.
func NewAPI(url string, auth Authenticator) API {
	return API{
		URL:           strings.TrimSuffix(url, "/"),
		Authenticator: auth,
	}
}

// Comment represents a comment on a TargetProcess item.
type Comment struct {
	Description string  `json:"Description"`
	General     General `json:"General"`
}

// General represents the TargetProcess General entity.
type General struct {
	// The ID of the entity.
	ID int `json:"Id"`
}

// Comment comments on a TargetProcess entity.
func (api API) Comment(entityID int, message string) (err error) {
	comment := Comment{
		Description: message,
		General: General{
			ID: entityID,
		},
	}
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(comment)
	if err != nil {
		return fmt.Errorf("failed to encode comment with error: %v", err)
	}
	r, err := http.NewRequest("POST", api.URL+"/api/v1/comments", b)
	if err != nil {
		return err
	}
	api.Authenticator(r)
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read response of TargetProcess API with error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("TargetProcess API returned status %v, expected OK. Body was: %v", resp.Status, string(body))
	}
	return nil
}
