package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/nats-io/stan.go"
)

func Set_Data(cache Cache, conn *pgx.Conn, m *stan.Msg) {
	var d JsonData

	err := json.Unmarshal(m.Data, &d)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(" Got data from NATS-streaming!")
	}

	if _, ok := cache.c[d.OrderUID]; !ok {
		cache.c[d.OrderUID] = d
		_, err = conn.Exec(context.Background(), "INSERT INTO orders VALUES ($1, $2)", d.OrderUID, m.Data)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("data has been set to DB: %s\n", d.OrderUID)
		}
	}
}

func main() {
	cache, conn := Work_With_DB()
	sc, err := stan.Connect("test-cluster", "Sub", stan.NatsURL("nats://localhost:4222"))
	if err != nil && err != io.EOF {
		log.Fatalln(err)
	} else {
		log.Println("Connecting to Nats-streaming!")
	}
	defer sc.Close()

	_, err = sc.Subscribe("foo1", func(m *stan.Msg) {
		Set_Data(cache, conn, m)
	}, stan.StartWithLastReceived())
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", cache.Handler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
