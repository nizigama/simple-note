package controllers

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/nizigama/simple-note/models"
	auth "github.com/nizigama/simple-note/services/middleware"
)

type HomeController struct {
	ViewsTemplate *template.Template
	Year          int
}

func NewHomeController() HomeController {
	tpl := template.Must(template.ParseGlob("./templates/includes/*.html"))
	tpl = template.Must(tpl.ParseGlob("./templates/*.html"))

	return HomeController{
		ViewsTemplate: tpl,
		Year:          time.Now().Year(),
	}
}

func (hc HomeController) Index(w http.ResponseWriter, req *http.Request) {

	c, err := req.Cookie("sessionID")
	data := map[string]interface{}{
		"loggedIn": false,
		"users":    []models.User{},
		"year":     hc.Year,
	}

	if err == nil {
		for _, v := range auth.Sessions {
			if c.Value == strconv.Itoa(int(v.ID)) {
				data["loggedIn"] = true
				break
			}
		}
	}

	allUsers, err := models.ReadAllUsers()

	if err != nil {
		http.Error(w, "Error getting all app users", http.StatusInternalServerError)
		return
	}

	data["users"] = allUsers

	w.Header().Set("Content-Type", "text/html")
	hc.ViewsTemplate.ExecuteTemplate(w, "index.html", data)
}
