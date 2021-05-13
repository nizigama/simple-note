package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/nizigama/simple-note/models"
	auth "github.com/nizigama/simple-note/services/middleware"
)

type NoteController struct {
	ViewsTemplate *template.Template
	Year          int
	UserID        int
}

func NewNoteController() NoteController {
	tpl := template.Must(template.ParseGlob("./templates/includes/*.html"))
	tpl = template.Must(tpl.ParseGlob("./templates/*.html"))

	return NoteController{
		ViewsTemplate: tpl,
		Year:          time.Now().Year(),
	}
}

func (nc NoteController) Index(w http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{
		"year": nc.Year,
	}
	w.Header().Set("Content-Type", "text/html")
	nc.ViewsTemplate.ExecuteTemplate(w, "note.html", data)
}

func (nc NoteController) Create(w http.ResponseWriter, req *http.Request) {
	_, nc.UserID = auth.GetLoggedInUser(req)

	if err := req.ParseForm(); err != nil {
		http.Error(w, "Error parsing your request", http.StatusUnprocessableEntity)
	}

	title, note := req.Form.Get("title"), req.Form.Get("note")

	if strings.Trim(title, " ") == "" {
		http.Error(w, "Note title is required", http.StatusUnprocessableEntity)
		return
	}
	if strings.Trim(note, " ") == "" {
		http.Error(w, "note is required", http.StatusUnprocessableEntity)
		return
	}

	newNote := models.Note{
		Title:   title,
		Body:    note,
		OwnerID: nc.UserID,
	}

	_, err := newNote.Save()

	if err != nil {
		http.Error(w, "Failed to save the note, Contact support", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/dashboard")
	w.WriteHeader(http.StatusSeeOther)

}

func (nc NoteController) Edit(w http.ResponseWriter, req *http.Request) {
	noteID := req.FormValue("noteID")

	id, err := strconv.Atoi(noteID)
	_, nc.UserID = auth.GetLoggedInUser(req)

	if err != nil {
		http.Error(w, "Invalid noteID", http.StatusUnprocessableEntity)
		return
	}

	note, err := models.ReadNote(uint64(id))

	if err != nil {
		http.Error(w, "No note found with ID", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"note": note,
		"year": nc.Year,
	}

	w.Header().Set("Content-Type", "text/html")
	nc.ViewsTemplate.ExecuteTemplate(w, "updateNote.html", data)
}

func (nc NoteController) Update(w http.ResponseWriter, req *http.Request) {
	noteID := req.FormValue("noteID")

	id, err := strconv.Atoi(noteID)
	_, nc.UserID = auth.GetLoggedInUser(req)

	if err != nil {
		http.Error(w, "Invalid noteID", http.StatusUnprocessableEntity)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Error(w, "Error parsing your request", http.StatusUnprocessableEntity)
	}

	title, note := req.Form.Get("title"), req.Form.Get("note")

	if strings.Trim(title, " ") == "" {
		http.Error(w, "Note title is required", http.StatusUnprocessableEntity)
		return
	}
	if strings.Trim(note, " ") == "" {
		http.Error(w, "note is required", http.StatusUnprocessableEntity)
		return
	}

	updatedNote := models.Note{
		Title:   title,
		Body:    note,
		OwnerID: nc.UserID,
	}

	err = models.UpdateNote(updatedNote, id)

	if err != nil {
		http.Error(w, "Failed to update the note, Contact support", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/dashboard")
	w.WriteHeader(http.StatusSeeOther)

}

func (nc NoteController) Delete(w http.ResponseWriter, req *http.Request) {
	noteID := req.FormValue("noteID")

	id, err := strconv.Atoi(noteID)

	if err != nil {
		http.Error(w, "Invalid noteID", http.StatusUnprocessableEntity)
		return
	}

	_, err = models.ReadNote(uint64(id))

	if err != nil {
		http.Error(w, "No note found with ID", http.StatusNotFound)
		return
	}

	err = models.DeleteNote(id)

	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/dashboard")
	w.WriteHeader(http.StatusSeeOther)

}
