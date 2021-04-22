package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger

func main() {

	logFile, err := os.OpenFile("app-log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)

	if err != nil {
		fmt.Println("Failed to create the log file")
		fmt.Println(err)
		os.Exit(1)
	}

	defer logFile.Close()

	logger = log.New(logFile, "simple-note-app:", log.Ldate)

	defer handlePanic()

	http.HandleFunc("/", index)

	err = http.ListenAndServe(":3000", nil)

	if err != nil {
		logger.Fatalln("Error starting the server, here is the error:", err)
	}
}

func handlePanic() {
	if message := recover(); message != nil {
		logger.Fatalln("Panic occurred:", message)
	}
}
