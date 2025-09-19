package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("[1/5] Chargement des variables d'environnement...")
	_ = godotenv.Load("C:/Users/WIN!!/Documents/google/brain/.env")

	user := os.Getenv("SupabaseName")
	password := os.Getenv("SupabaseMdp")
	host := os.Getenv("SupabaseUrl")
	dbName := os.Getenv("SupabaseDB")
	port := os.Getenv("SupabasePort")

	spt := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)

	fmt.Println("[2/5] Connexion à la base de données...")
	bdd, err := sql.Open("postgres", spt)
	if err != nil {
		panic(err)
	}
	defer bdd.Close()

	query := `SELECT "URL", paragraphe FROM main`
	fmt.Println("[3/5] Exécution de la requête SELECT sur la table main...")
	rows, err := bdd.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// map[word] => set d'URLs pour éviter doublons
	index := make(map[string]map[string]bool)
	countRows := 0

	fmt.Println("[4/5] Construction de l'index inversé...")
	for rows.Next() {
		var url, paragraphe string
		err := rows.Scan(&url, &paragraphe)
		if err != nil {
			fmt.Printf("Erreur scan: %s\n", err)
			continue
		}

		countRows++
		if countRows%100 == 0 {
			fmt.Printf("  - %d lignes traitées...\n", countRows)
		}

		words := strings.Fields(paragraphe)
		for _, w := range words {
			w = strings.ToLower(strings.Trim(w, ".,!?;:\"()"))
			if w == "" {
				continue
			}
			if index[w] == nil {
				index[w] = make(map[string]bool)
			}
			index[w][url] = true
		}
	}

	fmt.Printf("[4/5] Terminé, %d lignes traitées. Nombre de mots uniques : %d\n", countRows, len(index))

	fmt.Println("[5/5] Insertion de l'index dans la table invertedindex...")
	insertCount := 0
	for word, urlsMap := range index {
		urlsSlice := make([]string, 0, len(urlsMap))
		for u := range urlsMap {
			urlsSlice = append(urlsSlice, u)
		}

		query := `
			INSERT INTO invertedindex (word, urls)
			VALUES ($1, $2)
			ON CONFLICT (word) 
			DO UPDATE SET urls = array(SELECT DISTINCT unnest(invertedindex.urls || EXCLUDED.urls))
		`

		_, err := bdd.Exec(query, word, pq.Array(urlsSlice))
		if err != nil {
			fmt.Printf("Erreur insertion du mot '%s': %s\n", word, err)
			continue
		}
		insertCount++
		if insertCount%100 == 0 {
			fmt.Printf("  - %d mots insérés...\n", insertCount)
		}
	}

	fmt.Printf("Index inversé inséré avec succès ! Total mots insérés : %d\n", insertCount)
}
