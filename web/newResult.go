package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type FileContent struct {
	Urls       string   `json:"urls"`
	Text       string   `json:"text"`
	Links      []string `json:"link"`
	Title      string   `json:"title"`
	SubTitle   []string `json:"subtitle"`
	Paragraphe []string `json:"paragraph"`
}

func newSearch(data map[string]interface{}, w http.ResponseWriter) {
	err := godotenv.Load("C:/Users/WIN!!/Documents/google/brain/.env")
	if err != nil {
		http.Error(w, "Erreur chargement .env", http.StatusInternalServerError)
		fmt.Printf("Erreur .env: %s\n", err)
		return
	}

	value, ok := data["search"].(string)
	if !ok || value == "" {
		http.Error(w, "Paramètre 'search' manquant ou invalide", http.StatusBadRequest)
		return
	}

	fmt.Printf("Nouvelle recherche: %s\n", value)

	scriptPath := "C:/Users/WIN!!/Documents/google/brain/inSearch/crawling.py"
	cmd := exec.Command("python", scriptPath, value)
	output, err := cmd.CombinedOutput()

	if err != nil {
		http.Error(w, "Erreur exécution Python", http.StatusInternalServerError)
		fmt.Printf("Erreur Python:\n%s\n%s\n", err, string(output))
		return
	}

	fmt.Printf("Sortie Python:\n%s\n", string(output))

	file, err := os.ReadFile("regodt.json")
	if err != nil {
		http.Error(w, "Erreur lecture fichier JSON", http.StatusInternalServerError)
		fmt.Printf("Erreur lecture regodt.json: %s\n", err)
		return
	}

	var fileContents []FileContent
	if err := json.Unmarshal(file, &fileContents); err != nil {
		http.Error(w, "Erreur parsing JSON", http.StatusInternalServerError)
		fmt.Printf("Erreur unmarshal: %s\n", err)
		return
	}

	connStr := os.Getenv("DATABASE_URL")
	tableName := os.Getenv("tableName")

	if connStr == "" {
		host := os.Getenv("SupabaseURL")
		port := os.Getenv("SupabasePort")
		user := os.Getenv("SupabaseName")
		password := os.Getenv("SupabaseMDP")
		dbName := os.Getenv("SupabaseDB")

		if host == "" || user == "" || password == "" || dbName == "" {
			http.Error(w, "Variables d'environnement manquantes", http.StatusInternalServerError)
			fmt.Printf("Erreur: Variables DB manquantes\n")
			fmt.Printf("host=%s, user=%s, password=%s, dbName=%s\n", host, user, password != "", dbName)
			return
		}

		if port == "" {
			port = "5432"
		}

		poolerHost := fmt.Sprintf("aws-0-eu-central-1.pooler.supabase.com")

		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			host, port, user, password, dbName,
		)

		fmt.Printf("Tentative connexion avec pooler: %s\n", poolerHost)
	}

	fmt.Printf("Tentative de connexion à la DB...\n")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Erreur connexion DB", http.StatusInternalServerError)
		fmt.Printf("Erreur connexion DB: %s\n", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		http.Error(w, "DB injoignable", http.StatusInternalServerError)
		fmt.Printf("Erreur ping DB: %s\n", err)
		return
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (urls, title, subtitle, paragraphe, lien, categoris)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, tableName)

	insertedCount := 0
	var insertedUrls []string

	for _, content := range fileContents {
		linksJSON, _ := json.Marshal(content.Links)
		subTitleJSON, _ := json.Marshal(content.SubTitle)
		paragrapheJSON, _ := json.Marshal(content.Paragraphe)

		_, err = db.Exec(query,
			content.Urls,
			content.Title,
			string(subTitleJSON),
			string(paragrapheJSON),
			string(linksJSON),
			value,
		)

		if err != nil {
			fmt.Printf("Erreur insert pour %s: %s\n", content.Urls, err)
			continue
		}

		insertedCount++
		insertedUrls = append(insertedUrls, content.Urls)
		fmt.Printf("insert: %s\n", content.Urls)
	}

	fmt.Printf("Insert done: %d/%d URLs insert\n", insertedCount, len(fileContents))

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":       "success",
		"message":      "Crawling done",
		"totalUrls":    len(fileContents),
		"inserted":     insertedCount,
		"insertedUrls": insertedUrls,
		"search_query": value,
	}
	json.NewEncoder(w).Encode(response)
}
