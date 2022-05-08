package main

import (
	"html/template"
	"log"
	"net/http"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func writeList(writer http.ResponseWriter, list []string) {
	html, err := template.ParseFiles("list.html")
	check(err)
	err = html.Execute(writer, list)
	check(err)
}

func fruitHandler(writer http.ResponseWriter, request *http.Request) {
	writeList(writer, []string{"apples", "oranges", "pears"})
}

func meatHandler(writer http.ResponseWriter, request *http.Request) {
	writeList(writer, []string{"chicken", "beef", "lamb"})
}

func main() {
	http.HandleFunc("/fruit", fruitHandler)
	http.HandleFunc("/meat", meatHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
