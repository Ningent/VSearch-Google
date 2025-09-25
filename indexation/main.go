package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Document struct {
	URL        string
	Paragraphe string
}

type InvertedIndexEntry struct {
	Word      string
	URL       string
	TF        int
	DocLength int
}

func main() {
	dbUser := "postgres"
	dbPass := "dataBaseGooglePassworld1"
	dbName := "postgres"
	dbHost := "db.dqhyfbqwlkdfvuvtzdht.supabase.co"
	dbPort := "5432"

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Connexion DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Ping DB:", err)
	}
	fmt.Println("Connecté à Supabase/Postgres")

	rows, err := db.Query(`SELECT "URL", "paragraphe" FROM main`)
	if err != nil {
		log.Fatal("Erreur SELECT:", err)
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		var doc Document
		if err := rows.Scan(&doc.URL, &doc.Paragraphe); err != nil {
			log.Println("Erreur scan:", err)
			continue
		}
		documents = append(documents, doc)
	}
	fmt.Printf("%d documents récupérés\n", len(documents))

	for _, doc := range documents {
		words := strings.Fields(doc.Paragraphe)
		docLength := len(words)
		tfMap := make(map[string]int)

		for _, w := range words {
			w = strings.ToLower(strings.Trim(w, ".,!?;:\"()"))
			if w != "" {
				tfMap[w]++
			}
		}

		for word, tf := range tfMap {
			_, err := db.Exec(
				`INSERT INTO invertedindex ("word","urls","tf","doclength") VALUES ($1,$2,$3,$4)
					ON CONFLICT ("word","urls") DO UPDATE SET tf=$3, doclength=$4`,
				word, pq.Array([]string{doc.URL}), tf, docLength,
			)

			if err != nil {
				log.Println("Insert erreur:", err)
			} else {
				fmt.Printf("Inséré : %s → %s (tf=%d, len=%d)\n", word, doc.URL, tf, docLength)
			}
		}
	}

	fmt.Println("Index inversé terminé et inséré dans Supabase")
}
