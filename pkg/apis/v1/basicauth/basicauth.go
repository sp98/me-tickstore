package basicauth

import (
	"fmt"
	"net/http"
)

// New returns a piece of middleware that will allow access only
// if the provided credentials match within the given service
// otherwise it will return a 401 and not call the next handler.
func New(app string, credentials map[string][]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				unauthorized(w, app)
				return
			}

			validPasswords, userFound := credentials[username]
			if !userFound {
				unauthorized(w, app)
				return
			}

			for _, validPassword := range validPasswords {
				if password == validPassword {
					next.ServeHTTP(w, r)
					return
				}
			}

			unauthorized(w, app)
		})
	}
}

func unauthorized(w http.ResponseWriter, app string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic MarketEdge="%s"`, app))
	w.WriteHeader(http.StatusUnauthorized)
}
