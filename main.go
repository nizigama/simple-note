package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	http.HandleFunc("/", index)
	http.Handle("/login", auth.Authorize(login))
	http.Handle("/register", auth.Authorize(register))
	http.Handle("/profile", auth.Authorize(profile))
	http.Handle("/profile-picture", auth.Authorize(profilePicture))
	http.HandleFunc("/get-picture", getPicture)
	http.Handle("/dashboard", auth.Authorize(dashboard))
	http.Handle("/new-note", auth.Authorize(newNote))
	http.Handle("/update-note", auth.Authorize(updateNote))

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
