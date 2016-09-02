package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")

	if err != nil {
		showError(w, err)
		return
	}

	err = t.Execute(w, nil)

	if err != nil {
		showError(w, err)
		return
	}

	fmt.Fprintf(w, "")
}

func showError(w http.ResponseWriter, err error) {
	fmt.Println(err)
	http.Error(w, "Internal Server Error", 500)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8182", nil)
}
