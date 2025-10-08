package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getN(bdd *sql.DB) float64 {
	var n float64
	query := `SELECT COUNT(DISTINCT url) FROM invertedindex, unnest(urls) AS url;`
	err := bdd.QueryRow(query).Scan(&n)
	if err != nil {
		fmt.Printf("Erreur N -> %s\n", err)
		return 0
	}
	return n
}

func getDF(bdd *sql.DB) float64 {
	var df float64
	query := `SELECT array_length(urls, 1) FROM invertedindex;`
	err := bdd.QueryRow(query).Scan(&df)
	if err != nil {
		fmt.Printf("Erreur DF -> %s\n", err)
		return 0
	}
	return df
}

func main() {
	err := godotenv.Load("C:/Users/WIN!!/Documents/google/brain/.env")
	if err != nil {
		fmt.Printf("Impossible de charger '.env' -> %s\n", err)
		return
	}

	user := os.Getenv("SupabaseName")
	mdp := os.Getenv("SupabaseMDP")
	DbName := os.Getenv("SupabaseDB")
	host := os.Getenv("SupabaseURL")
	port := os.Getenv("SupabasePort")

	fmt.Printf("HOST -> %s\n", host)

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, mdp, DbName,
	)

	bdd, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Erreur bdd ouverture -> %s\n", err)
		return
	}
	defer bdd.Close()

	query := `SELECT "tf","doclength","id" FROM invertedindex`
	rows, err := bdd.Query(query)
	if err != nil {
		fmt.Printf("Erreur query -> %s\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tfcounte, doclen float64
		var id int32

		err = rows.Scan(&tfcounte, &doclen, &id)
		if err != nil {
			fmt.Printf("Erreur attribution -> %v\n", err)
			return
		}

		vraitf := tfcounte / doclen
		idf := math.Log(getN(bdd) / (1 + getDF(bdd)))
		tfIdf := vraitf * idf

		strVraiTf := strconv.FormatFloat(vraitf, 'f', 6, 64)
		strIdf := strconv.FormatFloat(idf, 'f', 6, 64)
		strTfIdf := strconv.FormatFloat(tfIdf, 'f', 6, 64)

		tx, err := bdd.Begin()
		if err != nil {
			fmt.Printf("Erreur begin -> %s\n", err)
			return
		}

		prep, err := tx.Prepare(`UPDATE invertedindex SET "vraiTf"=$1, "idf"=$2, "tfIdf"=$3 WHERE id=$4`)
		if err != nil {
			fmt.Printf("Erreur preparation -> %s\n", err)
			return
		}

		_, err = prep.Exec(strVraiTf, strIdf, strTfIdf, id)
		if err != nil {
			tx.Rollback()
			fmt.Printf("Erreur update -> %s\n", err)
			return
		}

		tx.Commit()
		prep.Close()
	}
}
