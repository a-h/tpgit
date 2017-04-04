package targetprocess

import (
	"net/http"
)

// TokenAuth adds the access_token to the querystring.
func TokenAuth(token string) func(r *http.Request) {
	return func(r *http.Request) {
		q := r.URL.Query()
		q.Set("access_token", token)
		r.URL.RawQuery = q.Encode()
	}
}
