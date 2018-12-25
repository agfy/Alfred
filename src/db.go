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

func fillValues(arr *[]string, rawName string) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT " + rawName + " FROM goods")
	if err != nil {
		log.Fatal(err)
	}

	var raw string
	for rows.Next() {
		err := rows.Scan(&raw)
		if err != nil {
			log.Fatal(err)
		}
		if !include(*arr, raw) {
			*arr = append(*arr, raw)
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return err
}

func getClasses(shop, foodType string) ([]string, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT class FROM goods WHERE shop='" + shop + "' AND foodtype='" + foodType + "'")
	if err != nil {
		log.Fatal(err)
	}

	classes := make([]string, 0)
	var class string
	for rows.Next() {
		err := rows.Scan(&class)
		if err != nil {
			log.Fatal(err)
		}
		if !include(classes, class) {
			classes = append(classes, class)
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return classes, err
}

func getGoods(shop, foodType, class, volume string) ([]good, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, name, price FROM goods WHERE shop='" + shop +
		"' AND foodtype='" + foodType + "' AND class='" + class + "' AND volume='" + volume + "'"
	println(query)
	rows, err := db.Query("SELECT id, name, price FROM goods WHERE shop='" + shop +
		"' AND foodtype='" + foodType + "' AND class='" + class + "' AND volume='" + volume + "'")
	if err != nil {
		log.Fatal(err)
	}

	gds := make([]good, 0)
	var gd good
	for rows.Next() {
		println(gd.Name)
		err := rows.Scan(&gd.Id, &gd.Name, &gd.Price)
		if err != nil {
			log.Fatal(err)
		}
		gds = append(gds, gd)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return gds, err
}
