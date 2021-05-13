package controllers

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/nizigama/simple-note/models"
	"github.com/nizigama/simple-note/services/helpers"
	auth "github.com/nizigama/simple-note/services/middleware"
	"golang.org/x/crypto/bcrypt"
)

type RegisterController struct {
	ViewsTemplate *template.Template
}

func NewRegisterController() RegisterController {
	tpl := template.Must(template.ParseGlob("./templates/includes/*.html"))
	tpl = template.Must(tpl.ParseGlob("./templates/*.html"))

	return RegisterController{
		ViewsTemplate: tpl,
	}
}

func (rc RegisterController) Index(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	rc.ViewsTemplate.ExecuteTemplate(w, "register.html", time.Now().Year())

}

func (rc RegisterController) Create(w http.ResponseWriter, req *http.Request) {

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
	if isValid := helpers.ValidateEmail(email); isValid != nil {
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

	newUser := models.User{
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

	w.Header().Set("Location", "/dashboard")
	w.WriteHeader(http.StatusSeeOther)
}
