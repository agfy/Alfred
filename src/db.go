package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
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
		if !include(*arr, raw) && raw != "" {
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

	rows, err := db.Query("SELECT class FROM goods WHERE shop=$1 AND foodtype=$2", shop, foodType)
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

	rows, err := db.Query("SELECT id, name, price FROM goods WHERE shop=$1 AND foodtype=$2 AND class=$3 "+
		"AND volume=$4", shop, foodType, class, volume)
	if err != nil {
		log.Fatal(err)
	}

	gds := make([]good, 0)
	var gd good
	for rows.Next() {
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

func createOrder(ownerId, amount, goodId int) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO orders (owner_telegram_id, buyer_telegram_id, goods_id, amount, create_time) "+
		"VALUES ($1, 0, $2, $3, now())", strconv.Itoa(ownerId), strconv.Itoa(goodId), strconv.Itoa(amount))
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func getOrders(userId int) ([]order, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, owner_telegram_id, goods_id, amount FROM orders WHERE buyer_telegram_id = 0 AND " +
		"create_time > now() - INTERVAL '12 hours'"
	if userId != 0 {
		query += "AND owner_telegram_id = " + strconv.Itoa(userId)
	}
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	orders := make([]order, 0)
	var ordr order
	for rows.Next() {
		err := rows.Scan(&ordr.Id, &ordr.OwnerTelegramId, &ordr.GoodId, &ordr.Amount)
		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, ordr)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return orders, err
}

func getGood(id int) (good, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return good{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, volume, price FROM goods WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	var gd good
	if !rows.Next() {
		return good{}, nil
	}
	err = rows.Scan(&gd.Name, &gd.Volume, &gd.Price)
	if err != nil {
		log.Fatal(err)
	}

	return gd, err
}

func getOrder(id int) (order, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return order{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT owner_telegram_id, goods_id, amount FROM orders WHERE id = " + strconv.Itoa(id))
	if err != nil {
		log.Fatal(err)
	}

	var ordr order
	if !rows.Next() {
		return order{}, nil
	}
	err = rows.Scan(&ordr.OwnerTelegramId, &ordr.GoodId, &ordr.Amount)
	if err != nil {
		log.Fatal(err)
	}

	return ordr, err
}

func markOrdersBought(orders *[]int, buyerId int) []int {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	affectedIds := make([]int, 0)

	for _, orderId := range *orders {
		result, err := db.Exec("UPDATE orders SET buyer_telegram_id = $1 WHERE id = $2 AND buyer_telegram_id = 0", strconv.Itoa(buyerId), strconv.Itoa(orderId))
		if err != nil {
			log.Fatal(err)
		}
		rawsAffected, _ := result.RowsAffected()
		if rawsAffected > 0 {
			affectedIds = append(affectedIds, orderId)
		}
	}

	return affectedIds
}

func deleteOrderfunc(id int) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Query("DELETE FROM orders WHERE id = " + strconv.Itoa(id))
	if err != nil {
		log.Fatal(err)
	}
	return err
}
