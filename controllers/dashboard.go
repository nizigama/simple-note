package controllers

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/nizigama/simple-note/models"
	auth "github.com/nizigama/simple-note/services/middleware"
)

type DashboardController struct {
	ViewsTemplate *template.Template
	Year          int
	UserID        int
}

func NewDashboardController() DashboardController {
	tpl := template.Must(template.ParseGlob("./templates/includes/*.html"))
	tpl = template.Must(tpl.ParseGlob("./templates/*.html"))

	return DashboardController{
		ViewsTemplate: tpl,
		Year:          time.Now().Year(),
	}
}

func (dc DashboardController) Index(w http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{
		"notes": []models.Note{},
		"year":  dc.Year,
	}

	_, dc.UserID = auth.GetLoggedInUser(req)

	allNotes, err := models.ReadAllUserNotes(dc.UserID)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error getting all app users", http.StatusInternalServerError)
		return
	}

	data["notes"] = allNotes

	dc.ViewsTemplate.ExecuteTemplate(w, "dashboard.html", data)
}
