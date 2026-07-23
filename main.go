package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/demo", demoHandler)

	log.Println("server is running ...")
	err :=http.ListenAndServe(":8080", nil)
	if err == nil {
		log.Fatal(err)
	}
}

func demoHandler(w http.ResponseWriter,r *http.Request) {
	log.Printf("%+v",r)

	if r.Method == "GET" {
		log.Println("run this!!!")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	response := map[string]string{
		"message": "hello world",
	}

	data,err := json.Marshal(response)
	if err != nil {
		http.Error(w,"error in encoding", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}