package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Person struct {
	First string `json:"firstname"`
}

func main() {
	p1 := Person{
		First: "Jenny",
	}

	p2 := Person{
		First: "James",
	}

	xp := []Person{p1, p2}

	bs, err := json.Marshal(xp)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(bs))

	xp2 := []Person{}
	if err := json.Unmarshal(bs, &xp2); err != nil {
		panic(err)
	}

	fmt.Println(xp2)

	http.HandleFunc("/encode", foo)
	http.HandleFunc("/decode", bar)
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {

}

func bar(w http.ResponseWriter, r *http.Request) {

}
