package targetprocess

import (
	"net/http"
	"testing"
)

func TestPasswordAuth(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com", nil)

	PasswordAuth("un", "pwd")(r)

	actualUsername, actualPassword, _ := r.BasicAuth()
	if actualUsername != "un" {
		t.Errorf("expected username 'un', but got '%v'.", actualUsername)
	}
	if actualPassword != "pwd" {
		t.Errorf("expected password 'pwd', but got '%v'.", actualPassword)
	}
}

func TestTokenAuth(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com", nil)

	TokenAuth("token_123")(r)

	actualToken := r.URL.Query().Get("access_token")

	if actualToken != "token_123" {
		t.Errorf("expected the access_token value to be set, but but got '%v'.", actualToken)
	}
}
