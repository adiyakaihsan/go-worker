package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var channel = make(chan message)

type message struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	fmt.Println("Ready to receive job!")
	router := httprouter.New()

	router.POST("/v1/api/jobs/submit", jobReceiver)
	go func() {
		http.ListenAndServe(":8081", router)
	}()
	for msg := range channel {
		jobProcessor(msg)
	}
}

func jobReceiver(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var msg message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		log.Printf("%s: %v", "Unable to read request body", err)
		return
	}

	channel <- msg

	w.Write([]byte("Job is sent"))
}

func jobProcessor(msg message) {
	log.Printf("====================================")
	log.Printf("Title: %s", msg.Title)
	log.Printf("Body: %s", msg.Body)
	log.Printf("====================================")
}

/*
Sample request:
curl localhost:8081/v1/api/jobs/submit -d '{"title":"P1", "body":"trigger alert"}'

*/
