package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// create a user and save the number of races he completed.

func raceCompletedhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	updateRaceCompleted("anshul")
	fmt.Fprintf(w, "Hello from golang webserver")
}

func raceCnt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	races := getRaceCompleted("anshul")
	fmt.Fprintf(w, "%d", races)
}

func main() {
	createMongoClient()
	// createUser("anshul")
	// insertUserProfile("anshul")
	// updateRaceCompleted("anshul")

	fmt.Println("Hello World")
	http.HandleFunc("/raceCompleted", raceCompletedhandler)
	http.HandleFunc("/races", raceCnt)

	http.ListenAndServe(":8081", nil)

	// Setup signal handling to catch SIGTERM
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-stopChan
		fmt.Printf("Received signal %v. Shutting down...\n", sig)

		// Disconnect MongoDB client before exiting
		DisconnectMongoClient()

		os.Exit(0)
	}()

}
