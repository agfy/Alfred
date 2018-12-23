package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"log"
	"os"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s", host, port, user, dbname, sslmode)

func getTelegramIds() ([]int, error) {
	log.Println(dbInfo)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT telegramid FROM users")
	if err != nil {
		log.Fatal(err)
	}

	var id int
	ids := make([]int, 0)
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return ids, err
}