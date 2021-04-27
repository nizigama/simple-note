package auth

import (
	"net/http"
	"strconv"

	users "github.com/nizigama/simple-note/models"
)

type Session struct {
	ID     uint64
	UserID int
}

var Sessions []Session

func Authorize(f func(http.ResponseWriter, *http.Request)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie("sessionID")
		if err != nil {
			// no cookie with sessionID found
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusSeeOther)
			return
		}
		for _, v := range Sessions {
			if c.Value == strconv.Itoa(int(v.ID)) {
				_, err := users.Read(uint64(v.UserID))
				if err != nil {
					// no user found
					w.Header().Set("Location", "/login")
					w.WriteHeader(http.StatusSeeOther)
					return
				}
				if req.URL.Path == "/login" || req.URL.Path == "/register" {
					// login and register are not available for logged in user
					w.Header().Set("Location", "/")
					w.WriteHeader(http.StatusSeeOther)
					return
				}
				// proceed with requested path handler
				f(w, req)
				return
			}
		}

		// no sessionID found
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusSeeOther)
	})
}

func CreateSession(sessionID uint64, userID int) {
	newSession := Session{
		ID:     sessionID,
		UserID: userID,
	}

	Sessions = append(Sessions, newSession)
}
