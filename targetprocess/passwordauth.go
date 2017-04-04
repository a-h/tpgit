package targetprocess

import (
	"net/http"
)

func PasswordAuth(username, password string) func(r *http.Request) {
	return func(r *http.Request) {
		r.SetBasicAuth(username, password)
	}
}
