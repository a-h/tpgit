package targetprocess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
		URL:      url,
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
func (api API) Comment(entityID int, message string) (id int, err error) {
	comment := Comment{
		Description: message,
		General: General{
			ID: entityID,
		},
	}
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(comment)
	if err != nil {
		return 0, fmt.Errorf("failed to encode comment with error: %v", err)
	}
	r, err := http.NewRequest("POST", api.URL+"/api/v1/comments", b)
	if err != nil {
		return 0, err
	}
	r.Header.Add("Content-Type", "application/json")
	r.SetBasicAuth(api.Username, api.Password)

	client := &http.Client{}
	resp, err := client.Do(r)

	fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return 0, fmt.Errorf("failed to read response of TargetProcess API with error: %v", err)
	}
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("TargetProcess API returned status %v, expected OK.", resp.Status)
	}
	return 0, nil
}
