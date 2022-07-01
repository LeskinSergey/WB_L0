package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nats-io/stan.go"
)

func main() {
	st, err := stan.Connect("test-cluster", "public", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Println(err)
	}
	var FlName string
	for {
		fmt.Println("insert file name")
		fmt.Scanln(&FlName)
		data, err := ioutil.ReadFile(FlName)
		if err != nil {
			log.Println(err)
		}
		st.Publish("foo1", data)
		fmt.Println("data send")
	}
}
