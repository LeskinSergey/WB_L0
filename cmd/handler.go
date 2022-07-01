package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func (cache *Cache) Handler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		tmpl, err := template.ParseFiles("./html/interface.html")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	case "POST":
		str := req.PostFormValue("order_uid")
		val, ok := cache.c[str]
		if !ok {
			log.Println("encorrect UID")
		} else {
			b, err := json.MarshalIndent(val, "", "\t")
			if err != nil {
				log.Println(err)
			}
			log.Printf("Data send: %s\n", str)
			fmt.Fprint(w, string(b))
		}
	}
}
