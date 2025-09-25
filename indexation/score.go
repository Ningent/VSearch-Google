package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"github.com/joho/godotenv"
	"strconv"
)

func getN(bdd sql.Conn) int8 {
	query := `SELECT COUNT(DISTINCT unnest(urls)) FROM invertedindex;`
	rows, err := bdd.Query("postgres", query)
	if err != nil {
		fmt.Printf("Erreur N -> %s\n")
		return
	}

	vor n int8
	for rows.Next() {
		err = rows.Scan(&n)
		if err != nil {
			fmt.Printf("Erreur N Scan => %s\n",err)
			return
		}
	}
	return n
}

func getDF (bdd sql.Conn) int8 {
	query := `SELECT id, word, array_length(urls, 1) AS docFreq FROM invertedindex;`
	rows , err := bdd.Query ("postgres",query)
	if err != nil {
		fmt.Printf("Erreur DF -> %s",err)
		return
	}

	var df int8

	for rows.Next() {
		err = rows.Scan(&n)
		if err != nil {
			fmt.Printf("Erreur N scan -> %s\n",err)
			return 
		}
	}

	return n
}


func main() {
	_ = godotenv.Load()

	user := os.Getenv("SupabaseName")
	mdp := os.Getenv("SupabaseMDP")
	DbName := os.Getenv("SupabaseDB")
	host := os.Getenv("SupabaseURL")
	port := os.Getenv("SupabasePort")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, mdp, DbName,
	)

	bdd, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Erreur bdd ouverture -> \t%s\n", err)
		return
	}

	query := `SELECT "tf","doclength","id" FROM invertedindex`

	rows, err = bdd.Query(query)
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
			fmt.Printf("Erreur attribution 'counte' 'doclen' -> \n %v\n")
			return
		}

		vraitf := tfcounte / doclen
		idf := math.Log(float64(getN(bdd) / float64(1 + getDF(bdd))))
		tfIdf = vraitf * idf

		strVraiTf = strconv.FormatFloat(vraitf,'f',6,64)
		strIdf = strconv.FormatFloat(idf,'f',6,64)
		strTfIdf = strconv.FormatFloat(tfIdf , 'f' , 6,64)

		beg, err := bdd.Begin()
		if err != nil {
			fmt.Printf("erreur beg -> %s\n", err)
			return
		}

		prep, err := beg.Prepare(`UPDATE invertedindex SET vraiTf = $1, idf = $2 , tfIdf = $4 WHERE id = $3`)
		if err != nil {
			fmt.Printf("Erreur preparation -> %s", err)
			return
		}

		defer prep.Close()

		for _, u := range invertedIndex {
			_, err = prep.Exec(id, strVraiTf, strIdf,strTfIdf)
			if err != nil {
				beg.Rollback()
				fmt.Printf("Erreur Insert -> %s\n")
				return
			}
		}

	}
}
