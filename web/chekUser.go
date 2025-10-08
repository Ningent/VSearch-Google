package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func chekUsers(data map[string]interface{}, w http.ResponseWriter) {
	ip := data["ip"].(string)
	uuid := data["uuid"].(string)
	theme := data["theme"].(string)

	fmt.Printf("test -> %s\n%s\n%s\n", ip, uuid, theme)

	_ = godotenv.Load()

	host := os.Getenv("SupabaseURL")
	port := os.Getenv("SupabasePort")
	user := os.Getenv("SupabaseName")
	password := os.Getenv("SupabaseMDP")
	dbName := os.Getenv("SupabaseDB")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbName,
	)

	bdd, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error in bdd open -> %s\n")
		return
	}

	defer bdd.Close()

	table := os.Getenv("tableName2")
	query := `SELECT * FROM %s WHERE uuid = $1`
	sendQuery, err := bdd.Query(query, table, uuid)

	var dbUuid string
	err = sendQuery.Scan(&dbUuid)

	if err != nil {
		fmt.Printf("Error to have dbUuid\n%s\n", err)
		return
	}

	if dbUuid == nil {
		//insert le nouveau users a la bdd
		//mes le theme dans le bdd
	} else {
		//envoie le theme aux front
	}

	packag := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packag)
}
