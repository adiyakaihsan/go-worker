package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	fmt.Printf("Hallo %v", "okokok")
	router := httprouter.New()

	router.POST("/v1/api/jobs/submit", jobReceiver)
	http.ListenAndServe(":8081", router)
}

func jobReceiver(w http.ResponseWriter, r * http.Request, _ httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}
	jobProcessor(string(body))
	w.Write([]byte(body))
}

func jobProcessor(message string) {
	fmt.Printf("Processing Job %v", message)
}