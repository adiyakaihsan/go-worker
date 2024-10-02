package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

var channel = make(chan message)

type message struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	var wg sync.WaitGroup

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	router := httprouter.New()

	server := &http.Server{Addr: ":8081", Handler: router}

	fmt.Println("Ready to receive job!")

	router.POST("/v1/api/jobs/submit", jobReceiver)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			return
		}
	}()
	go func() {
		for msg := range channel {
			wg.Add(1)
			go func() {
				defer wg.Done()
				jobProcessor(msg)
			}()

		}
	}()

	sig := <-sigChan
	log.Printf("Caught Signal %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return
	}

	close(channel)

	wg.Wait()

	log.Println("All jobs processed, shutting down")
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
	rand.Seed(time.Now().UnixNano())

	min := 1
	max := 5
	delay := rand.Intn(max-min+1) + min
	// delay := 5
	time.Sleep(time.Duration(delay) * time.Second)

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
