package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

func Work_With_DB() (Cache, *pgx.Conn) {
	var cache Cache
	cache.c = make(map[string]JsonData)
	conf := fmt.Sprintf("postgres://%s:%s@localhost:5432/postgres?sslmode=disable", "postgres", "postgres")
	conn, err := pgx.Connect(context.Background(), conf)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("connected to DB")
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(context.Background(), "SELECT data_js FROM orders")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var dst []byte
		var data JsonData
		err = rows.Scan(&dst)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(dst, &data)
		if err != nil {
			log.Println(err)
		}
		cache.c[data.OrderUID] = data
	}
	rows.Close()
	return cache, conn
}
