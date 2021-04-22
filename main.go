package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func init() {
	logFile, err := os.Open("app-log")

	if err != nil {
		logFile, err = os.Create("app-log")
		if err != nil {
			fmt.Println("Failed to create the log file")
			fmt.Println(err)
			os.Exit(1)
		}
	}

	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetPrefix("simple-note-app:")
}

func main() {

	defer handlePanic()

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		log.Fatalln("Error starting the server, here is the error:", err)
	}
}

func handlePanic() {
	if message := recover(); message != nil {
		log.Fatalln("Application paniced, here is the message:", message)
	}
}
