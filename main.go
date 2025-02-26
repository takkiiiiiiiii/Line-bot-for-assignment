package main

import (
	"log"
	"net/http"
)


func main() {
	http.HandleFunc("/kadai", ReplyKadai)
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Print(err)
	}
}