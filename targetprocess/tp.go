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
	URL      string
	Username string
	Password string
}

// NewAPI creates a new TargetProcess endpoint.
func NewAPI(url string, username string, password string) API {
	return API{
		URL:      strings.TrimSuffix(url, "/"),
		Username: username,
		Password: password,
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
	r.Header.Add("Content-Type", "application/json")
	r.SetBasicAuth(api.Username, api.Password)

	client := &http.Client{}
	resp, err := client.Do(r)

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read response of TargetProcess API with error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("TargetProcess API returned status %v, expected OK. Body was: %v", resp.Status, string(body))
	}
	return nil
}
