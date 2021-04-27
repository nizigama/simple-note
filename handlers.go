package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	users "github.com/nizigama/simple-note/models"
	auth "github.com/nizigama/simple-note/services/middleware"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("./templates/*.html"))
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "index.html", time.Now().Year())
}

func login(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "login.html", time.Now().Year())
}

func register(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		if err := req.ParseForm(); err != nil {
			http.Error(w, "Error parsing your request", http.StatusUnprocessableEntity)
		}

		firstName, lastName, email, password, confirm := req.Form.Get("first"), req.Form.Get("last"), req.Form.Get("email"), req.Form.Get("pass"), req.Form.Get("confirm")

		if strings.Trim(firstName, " ") == "" {
			http.Error(w, "First name is required", http.StatusUnprocessableEntity)
			return
		}
		if strings.Trim(lastName, " ") == "" {
			http.Error(w, "Last name is required", http.StatusUnprocessableEntity)
			return
		}
		if strings.Trim(email, " ") == "" {
			http.Error(w, "Email is required", http.StatusUnprocessableEntity)
			return
		}
		if isValid := validateEmail(email); isValid != nil {
			http.Error(w, "invalid email address", http.StatusUnprocessableEntity)
			return
		}
		if strings.Trim(password, " ") == "" {
			http.Error(w, "Password is required", http.StatusUnprocessableEntity)
			return
		}
		if strings.Trim(confirm, " ") == "" {
			http.Error(w, "Password confirmation is required", http.StatusUnprocessableEntity)
			return
		}
		if password != confirm {
			http.Error(w, "Passwords don't match", http.StatusUnprocessableEntity)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

		if err != nil {
			http.Error(w, "Password hashing failed, Contact support", http.StatusInternalServerError)
			return
		}

		newUser := users.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  hashedPassword,
		}

		userID, err := newUser.Save()

		if err != nil {
			http.Error(w, "Signup failed, Contact support", http.StatusInternalServerError)
			return
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		sessionID := r.Uint64()

		auth.CreateSession(sessionID, int(userID))

		http.SetCookie(w, &http.Cookie{
			Name:  "sessionID",
			Value: strconv.Itoa(int(sessionID)),
		})

		w.Header().Set("Location", "/profile")
		w.WriteHeader(http.StatusSeeOther)

	} else {
		w.Header().Set("Content-Type", "text/html")
		tpl.ExecuteTemplate(w, "register.html", time.Now().Year())
	}
}

func profile(w http.ResponseWriter, req *http.Request) {
	// w.Header().Set("Content-Type", "text/html")
	// tpl.ExecuteTemplate(w, "profile.html", time.Now().Year())
}

func validateEmail(email string) error {

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email")
	}

	parts := strings.Split(email, "@")

	// if there is no content before or after the @ symbol
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return fmt.Errorf("invalid email")
	}

	afterAtSymbol := strings.Split(parts[1], ".")
	// if there is a dot after the @ symbol
	if len(afterAtSymbol) == 0 {
		return fmt.Errorf("invalid email")
	}

	// if there is content after the last dot(.)
	if len(afterAtSymbol[len(afterAtSymbol)-1]) == 0 {
		return fmt.Errorf("invalid email")
	}

	return nil
}
