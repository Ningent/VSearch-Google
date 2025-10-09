package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func chekUsers(data map[string]interface{}, w http.ResponseWriter) {
	ip := data["ip"].(string)
	uuid := data["uuid"].(string)
	theme := data["theme"].(string)

	fmt.Printf("test -> %s\n%s\n%s\n", ip, uuid, theme)

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

	table := os.Getenv("tableName2")
	query := fmt.Sprintf(`SELECT uuid FROM %s WHERE uuid = $1`, table)

	rows, err := bdd.Query(query, uuid)
	if err != nil {
		fmt.Printf("Erreur lors du SELECT : %v\n", err)
		return
	}
	defer rows.Close()

	var dbUuid string
	found := false
	if rows.Next() {
		err = rows.Scan(&dbUuid)
		if err != nil {
			fmt.Printf("Erreur de lecture : %v\n", err)
			return
		}
		found = true
		fmt.Println("UUID trouvé :", dbUuid)
	} else {
		found = false
		fmt.Println("Aucun utilisateur trouvé, on peut en créer un nouveau.")
	}

	if !found {
		fmt.Printf("%s\n", found)
		col1 := os.Getenv("coll2T1")
		col2 := os.Getenv("coll2T2")
		col3 := os.Getenv("coll2T3")

		query = fmt.Sprintf(
			"INSERT INTO %s (%s, %s, %s) VALUES ($1, $2, $3)",
			table, col1, col2, col3,
		)
		_, err = bdd.Exec(query, ip, uuid, theme)
		if err != nil {
			fmt.Printf("Erreur insert data : %s\n", err)
			return
		}

		fmt.Print("Insertion done\n")
	} else {
		fmt.Printf("%s\n", found)
		query = fmt.Sprintf("SELECT theme FROM %s WHERE uuid = $1", table)
		row := bdd.QueryRow(query, uuid)

		var bddTheme string
		err = row.Scan(&bddTheme)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Aucun theme trouvé pour cet UUID")
			} else {
				fmt.Printf("Erreur lecture theme : %v\n", err)
			}
			return
		}

		fmt.Printf("bddTheme -> %s\n", bddTheme)

	}

	packag := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packag)
}
