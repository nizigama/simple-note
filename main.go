package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nizigama/simple-note/controllers"
	"github.com/nizigama/simple-note/models"
	boltDB "github.com/nizigama/simple-note/services/database"
	auth "github.com/nizigama/simple-note/services/middleware"
)

// TODO: work on login feature
var (
	logger  *log.Logger
	logFile *os.File
)

func init() {

	logFile, err := os.OpenFile("app-log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)

	if err != nil {
		fmt.Println("Failed to create the log file")
		fmt.Println(err)
		os.Exit(1)
	}

	logger = log.New(logFile, "simple-note-app:", log.Ldate)

	boltDB.InitDBConnection(logger, models.UserTableName, models.NoteTableName)
}

func main() {
	defer boltDB.CloseDBConnection()
	defer logFile.Close()
	defer handlePanic()

	r := mux.NewRouter()
	hc := controllers.NewHomeController()
	lc := controllers.NewLoginController()
	rc := controllers.NewRegisterController()
	ac := controllers.NewAccountController()
	dc := controllers.NewDashboardController()
	nc := controllers.NewNoteController()

	r.HandleFunc("/", hc.Index).Methods("GET")
	r.Handle("/login", auth.Authorize(lc.Index)).Methods("GET")
	r.Handle("/login", auth.Authorize(lc.Login)).Methods("POST")
	r.Handle("/register", auth.Authorize(rc.Index)).Methods("GET")
	r.Handle("/register", auth.Authorize(rc.Create)).Methods("POST")
	r.Handle("/logout", auth.Authorize(lc.Logout)).Methods("GET")
	r.Handle("/profile", auth.Authorize(ac.Index)).Methods("GET")
	r.Handle("/profile", auth.Authorize(ac.Save)).Methods("POST")
	r.Handle("/profile-picture", auth.Authorize(ac.ProfilePicture)).Methods("GET")
	r.HandleFunc("/get-picture", ac.GetPicture).Methods("GET")
	r.Handle("/delete-account", auth.Authorize(ac.DeleteAccount)).Methods("GET")
	r.Handle("/dashboard", auth.Authorize(dc.Index)).Methods("GET")
	r.Handle("/new-note", auth.Authorize(nc.Index)).Methods("GET")
	r.Handle("/new-note", auth.Authorize(nc.Create)).Methods("POST")
	r.Handle("/update-note", auth.Authorize(nc.Edit)).Methods("GET")
	r.Handle("/update-note", auth.Authorize(nc.Update)).Methods("POST")
	r.Handle("/delete-note", auth.Authorize(nc.Delete)).Methods("GET")

	err := http.ListenAndServe(":3000", r)

	if err != nil {
		logger.Fatalln("Error starting the server, here is the error:", err)
	}
}

func handlePanic() {
	if message := recover(); message != nil {
		logger.Fatalln("Panic occurred:", message)
	}
}
