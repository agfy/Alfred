package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s", host, port, user, dbname, sslmode)

func getUsers() (map[int]string, error) {
	log.Println(dbInfo)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT telegramid, nickname FROM users")
	if err != nil {
		log.Fatal(err)
	}

	var id int
	var nickName string
	users := make(map[int]string)
	for rows.Next() {
		err := rows.Scan(&id, &nickName)
		if err != nil {
			log.Fatal(err)
		}
		users[id] = nickName
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return users, err
}
