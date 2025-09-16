package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
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

	bdd, err := sql.Open("postgres", spt)
	if err != nil {
		panic(err)
	}

	defer bdd.Close()

	query := `SELECT paragraphe FROM main`

	row, err := bdd.Query(query)

	if err != nil {
		panic(err)
	}

	fileName := "dicc.json"
	var file *os.File

	_, err = os.Stat(fileName)
	fileExist := os.IsNotExist(err)

	if fileExist {
		file, err = os.Create(fileName)
		if err != nil {
			fmt.Printf("Erreur lors de la creation du fichier : %s\n", err)
			return
		}
	}

	defer file.Close()

	for row.Next() {
		var paragraphe string

		err := row.Scan(&paragraphe)

		if err != nil {
			log.Fatal(err)
		}

		// for i,e := range paragraphe {

		// }

	}

	fmt.Printf("tu clc github mais az")
}
