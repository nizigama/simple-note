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
)

type LoginController struct {
	ViewsTemplate *template.Template
}

func NewLoginController() LoginController {
	tpl := template.Must(template.ParseGlob("./templates/includes/*.html"))
	tpl = template.Must(tpl.ParseGlob("./templates/*.html"))

	return LoginController{
		ViewsTemplate: tpl,
	}
}

func (lc LoginController) Index(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lc.ViewsTemplate.ExecuteTemplate(w, "login.html", time.Now().Year())
}

func (lc LoginController) Login(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Error parsing your request", http.StatusUnprocessableEntity)
	}

	email, password := req.Form.Get("email"), req.Form.Get("pass")

	if isValid := helpers.ValidateEmail(email); isValid != nil {
		http.Error(w, "invalid email address", http.StatusUnprocessableEntity)
		return
	}
	if strings.Trim(password, " ") == "" {
		http.Error(w, "Password is required", http.StatusUnprocessableEntity)
		return
	}

	user, userID, err := models.ReadSingleUserByEmail(email)
	if err != nil {
		http.Error(w, "Wrong credentials, there is no such email in our records", http.StatusForbidden)
		return
	}

	err = helpers.ValidatePasswords([]byte(user.Password), []byte(password))

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

	w.Header().Set("Location", "/dashboard")
	w.WriteHeader(http.StatusSeeOther)
}

func (lc LoginController) Logout(w http.ResponseWriter, req *http.Request) {
	c, _ := req.Cookie("sessionID")

	var sessionID int

	for k, v := range auth.Sessions {
		if c.Value == strconv.Itoa(int(v.ID)) {
			sessionID = k
			break
		}
	}

	first := auth.Sessions[:sessionID]
	second := auth.Sessions[sessionID+1:]

	auth.Sessions = append(first, second...)

	c.MaxAge = -1
	http.SetCookie(w, c)

	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusSeeOther)
}
