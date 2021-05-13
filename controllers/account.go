package controllers

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

	"github.com/nizigama/simple-note/models"
	auth "github.com/nizigama/simple-note/services/middleware"
)

type AccountController struct {
	ViewsTemplate *template.Template
	Year          int
	UserAccount   models.User
	userID        int
}

func NewAccountController() AccountController {
	tpl := template.Must(template.ParseGlob("./templates/includes/*.html"))
	tpl = template.Must(tpl.ParseGlob("./templates/*.html"))

	return AccountController{
		ViewsTemplate: tpl,
		Year:          time.Now().Year(),
	}
}

func (ac AccountController) Index(w http.ResponseWriter, req *http.Request) {
	ac.UserAccount, ac.userID = auth.GetLoggedInUser(req)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	w.Header().Set("Content-Type", "text/html")
	ac.ViewsTemplate.ExecuteTemplate(w, "profile.html", map[string]interface{}{
		"user": ac.UserAccount,
		"r":    r.Uint64(),
		"year": ac.Year,
	})
}

func (ac AccountController) Save(w http.ResponseWriter, req *http.Request) {

	ac.UserAccount, ac.userID = auth.GetLoggedInUser(req)

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
	if ac.UserAccount.Picture != "avatar.png" {
		os.Remove("./assets/" + ac.UserAccount.Picture)
	}

	ac.UserAccount.Picture = fileName

	if err := models.UpdateUser(ac.UserAccount, ac.userID); err != nil {
		http.Error(w, "Error updating database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/profile")
	w.WriteHeader(http.StatusSeeOther)

}

func (ac AccountController) ProfilePicture(w http.ResponseWriter, req *http.Request) {
	ac.UserAccount, _ = auth.GetLoggedInUser(req)

	http.ServeFile(w, req, "assets/"+ac.UserAccount.Picture)
}

func (ac AccountController) GetPicture(w http.ResponseWriter, req *http.Request) {

	picture := req.FormValue("name")

	_, err := os.Open("assets/" + picture)

	if err != nil {
		http.ServeFile(w, req, "assets/avatar.png")
		return
	}

	http.ServeFile(w, req, "assets/"+picture)
}

func (ac AccountController) DeleteAccount(w http.ResponseWriter, req *http.Request) {
	c, _ := req.Cookie("sessionID")

	var sessionID int
	var userID int

	for k, v := range auth.Sessions {
		if c.Value == strconv.Itoa(int(v.ID)) {
			sessionID = k
			userID = v.UserID
			break
		}
	}

	allNotes, err := models.ReadAllUserNotes(userID)

	if err != nil {
		http.Error(w, "Error getting user's notes", http.StatusInternalServerError)
		return
	}

	for _, note := range allNotes {
		err = models.DeleteNote(note.ID)

		if err != nil {
			http.Error(w, "Failed to delete note", http.StatusInternalServerError)
			return
		}
	}

	err = models.DeleteUser(userID)

	if err != nil {
		http.Error(w, "Failed to delete the user", http.StatusInternalServerError)
		return
	}

	first := auth.Sessions[:sessionID]
	second := auth.Sessions[sessionID+1:]

	auth.Sessions = append(first, second...)

	c.MaxAge = -1
	http.SetCookie(w, c)

	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusSeeOther)
}
