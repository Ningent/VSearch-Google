package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func searchGo(data map[string]interface{}, w http.ResponseWriter) {
	value := data["value"].(string)

	_ = godotenv.Load("C:/Users/WIN!!/Documents/google/brain/.env")

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
		fmt.Printf("Error in bdd open -> %s\n", err)
		return
	}

	defer bdd.Close()

	coll1 := os.Getenv("coll0T3")
	coll2 := os.Getenv("coll1T3")
	coll3 := os.Getenv("coll2")

	table := os.Getenv("TableName3")

	query := fmt.Sprintf(
		`SELECT "%s", "%s", "%s" FROM "%s" WHERE "%s" = $1`,
		coll1, coll2, coll3, table, coll1,
	)

	querySend, err := bdd.Query(query, value)

	if err != nil {
		fmt.Printf("Error (Searching.go) send query\n%s\n", err)
		return
	}

	for querySend.Next() {
		var world string
		var linkArr []string
		var tfIdf string

		err = querySend.Scan(&world, pq.Array(&linkArr), &tfIdf)

		if err != nil {
			fmt.Printf("Error (Searching.go) read data\n%s\n", err)
			return
		}

		packag := map[string]interface{}{
			"world": world,
			"link":  linkArr,
			"tfIdf": tfIdf,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(packag)
	}
}
