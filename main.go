package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	r.HandleFunc("/", index).Methods("GET")
	r.Handle("/login", auth.Authorize(login))
	r.Handle("/register", auth.Authorize(register))
	r.Handle("/profile", auth.Authorize(profile))
	r.Handle("/profile-picture", auth.Authorize(profilePicture)).Methods("GET")
	r.HandleFunc("/get-picture", getPicture).Methods("GET")
	r.Handle("/dashboard", auth.Authorize(dashboard)).Methods("GET")
	r.Handle("/new-note", auth.Authorize(newNote))
	r.Handle("/update-note", auth.Authorize(updateNote))
	r.Handle("/delete-note", auth.Authorize(deleteNote)).Methods("GET")
	r.Handle("/logout", auth.Authorize(logout)).Methods("GET")
	r.Handle("/delete-account", auth.Authorize(deleteAccount)).Methods("GET")

	http.Handle("/", r)

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		logger.Fatalln("Error starting the server, here is the error:", err)
	}
}

func handlePanic() {
	if message := recover(); message != nil {
		logger.Fatalln("Panic occurred:", message)
	}
}
