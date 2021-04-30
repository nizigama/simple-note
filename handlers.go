package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
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
	tpl = template.Must(template.ParseGlob("./templates/includes/*.html"))
	tpl = template.Must(tpl.ParseGlob("./templates/*.html"))
}

func index(w http.ResponseWriter, req *http.Request) {

	c, err := req.Cookie("sessionID")
	data := map[string]interface{}{
		"year":     time.Now().Year(),
		"loggedIn": false,
	}

	if err == nil {
		// var user users.User
		// var userID int

		for _, v := range auth.Sessions {
			if c.Value == strconv.Itoa(int(v.ID)) {
				// userID = v.UserID
				// user, _ = users.Read(uint64(v.UserID))
				data["loggedIn"] = true
				break
			}
		}
	}

	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "index.html", data)
}

func login(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			http.Error(w, "Error parsing your request", http.StatusUnprocessableEntity)
		}

		email, password := req.Form.Get("email"), req.Form.Get("pass")

		if isValid := validateEmail(email); isValid != nil {
			http.Error(w, "invalid email address", http.StatusUnprocessableEntity)
			return
		}
		if strings.Trim(password, " ") == "" {
			http.Error(w, "Password is required", http.StatusUnprocessableEntity)
			return
		}

		user, userID, err := users.ReadSingleByEmail(email)
		if err != nil {
			http.Error(w, "Wrong credentials, there is no such email in our records", http.StatusForbidden)
			return
		}

		err = auth.CheckCredentials([]byte(user.Password), []byte(password))

		if err != nil {
			http.Error(w, "Wrong credentials", http.StatusForbidden)
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

	}
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
			Password:  string(hashedPassword),
			Picture:   "avatar.png",
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

	year := time.Now().Year()
	user, userID := getLoggedInUser(req)

	// TODO: code to upload and save profile picture
	if req.Method == http.MethodPost {
		f, h, err := req.FormFile("picture")

		if err != nil {
			http.Error(w, "No image file found", http.StatusNotFound)
			return
		}
		xb, err := ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		xs := strings.Split(h.Filename, ".")
		ext := xs[len(xs)-1]
		tsp := time.Now().UnixNano()

		fileName := fmt.Sprintf("%v.%s", tsp, ext)
		file, err := os.Create("./assets/" + fileName)

		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		file.Write(xb)
		if user.Picture != "avatar.png" {
			os.Remove("./assets/" + user.Picture)
		}

		user.Picture = fileName

		if err := users.UpdateUser(user, userID); err != nil {
			http.Error(w, "Error updating database", http.StatusInternalServerError)
			return
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "profile.html", map[string]interface{}{
		"user": user,
		"r":    r.Uint64(),
		"year": year,
	})
}

func profilePicture(w http.ResponseWriter, req *http.Request) {
	user, _ := getLoggedInUser(req)

	http.ServeFile(w, req, "assets/"+user.Picture)
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

func getLoggedInUser(req *http.Request) (users.User, int) {

	c, _ := req.Cookie("sessionID")

	var user users.User
	var userID int

	for _, v := range auth.Sessions {
		if c.Value == strconv.Itoa(int(v.ID)) {
			userID = v.UserID
			user, _ = users.Read(uint64(v.UserID))
			break
		}
	}

	return user, userID
}
