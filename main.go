package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// TODO:
	// 1. Add sleep simulation di jobProcessor, untuk simulate long processing
	// 2. How to implement: Graceful shutdown
	// 2.1. Close HTTP dulu, make sure no new requests coming in
	// 2.2. CTRL+C dihold sampe channel kosong + last job completed

	// TODO:
	// Explore sync.WaitGroup
	// Perdalam behavior Goroutine
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
	txt := fmt.Sprintf("Title: %s\nBody: %s", msg.Title, msg.Body)
	saveStringToFile(txt)

}

func saveStringToFile(data string) error {
	fileName := fmt.Sprintf("output_%s.txt", time.Now().Format("20060102_150405"))

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	fmt.Println("String successfully saved to", fileName)
	return nil
}

/*
Sample request:
curl localhost:8081/v1/api/jobs/submit -d '{"title":"P1", "body":"trigger alert"}'

*/
