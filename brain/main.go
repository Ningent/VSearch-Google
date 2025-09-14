package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	user := os.Getenv("SupabaseName")
	password := os.Getenv("SupabaseMDP")
	dbName := os.Getenv("SupabaseDB")
	host := os.Getenv("SupabaseURL")
	port := os.Getenv("SupabasePort")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	read := csv.NewReader(file)
	read.LazyQuotes = true
	read.FieldsPerRecord = -1
	rows, err := read.ReadAll()
	if err != nil {
		panic(err)
	}

	if len(rows) <= 1 {
		fmt.Println("CSV vide ou seulement header")
		return
	}

	table := os.Getenv("TableName")
	cols := `"URL", title, "subTitle", paragraphe, lien, categoris`

	var placeholders []string
	var values []interface{}
	paramCounter := 1

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 6 {
			continue
		}

		ph := []string{}
		for range row[:6] {
			ph = append(ph, fmt.Sprintf("$%d", paramCounter))
			paramCounter++
		}
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(ph, ",")))

		values = append(
			values,
			cleanUTF8(row[0]),
			cleanUTF8(row[1]),
			cleanUTF8(row[2]),
			cleanUTF8(row[3]),
			cleanUTF8(row[4]),
			cleanUTF8(row[5]),
		)

		// values = append(values, row[0], row[1], row[2], row[3], row[4], row[5])
	}

	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES %s`,
		table, cols, strings.Join(placeholders, ","),
	)

	_, err = db.Exec(query, values...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Insertion bulk réussie : %d lignes insérées\n", len(values)/6)

	ClearFile("data.csv")
	WriteByte("binary.bin", 1)

	binRead, err := ReadByte("binary.bin")

	if err != nil {
		fmt.Printf("Erreur lecture 'binRead':\n%v\n", err)
		return
	}

	if binRead == 0 {
		err := WriteByte("binary.bin", 0)

		if err != nil {
			fmt.Printf("Erreur read 'bidread'\n%v\n", err)
			return
		}
	}

	os.Exit(1)
}

func ClearFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		fmt.Printf("Error supresion data 'datacsv' from main.go \n%v\n", err)
		return nil
	}

	fmt.Printf("Suprimer data from 'datacsv' from go succes\n")

	defer file.Close()
	return nil
}

func WriteByte(file string, value byte) error {
	return os.WriteFile(file, []byte{value}, 0644)
}

func ReadByte(file string) (byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Erreur function ReadByte from go\n%v\n", err)
		return 0, err
	}

	if len(data) == 0 {
		fmt.Printf("Erreur fichier vide from go (ReadByte)\n")
		return 0, fmt.Errorf("Erreur fichier vide from go (ReadByte)\n")
	}

	return data[0], nil
}

func cleanUTF8(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == 0 || r == '\uFFFD' {
			b.WriteRune(' ')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
