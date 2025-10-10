package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Requete struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Action")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == "POST" {
		var requete Requete

		err := json.NewDecoder(r.Body).Decode(&requete)
		if err != nil {
			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		if requete.Action == "chekUsers" {
			chekUsers(requete.Data, w)
		} else if requete.Action == "changeTheme" {
			changeTheme(requete.Data, w)
		}
	}
}

func main() {
	fmt.Printf("welcome to backend\n")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/back/main", handler)
	http.ListenAndServe(":"+port, nil)
}
