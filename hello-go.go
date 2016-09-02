package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
)

var db *sqlx.DB

type Person struct {
	Name string
}

func webServerHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")

	if err != nil {
		showError(w, err)
		return
	}

	data := struct {
		Who string
	}{
		Who: "World",
	}

	err = t.Execute(w, data)

	if err != nil {
		showError(w, err)
		return
	}

	fmt.Fprintf(w, "")
}

func connectDb() error {
	var err error
	var schema = `
	CREATE TABLE IF NOT EXISTS people (
		name TEXT
	)
	`

	db, err = sqlx.Connect("postgres", "user=dev dbname=hello-go sslmode=disable")

	if err == nil {
		db.MustExec(schema)
	}

	return err
}

func saveName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	_, err := db.NamedExec("INSERT INTO people (name) VALUES (:name)", Person{r.PostForm.Get("name")})

	if err != nil {
		showError(w, err)
		return
	}

	fmt.Fprintf(w, "ok")
}

func showError(w http.ResponseWriter, err error) {
	fmt.Println(err)
	http.Error(w, "Internal Server Error", 500)
}

func main() {
	dbErr := connectDb()

	if dbErr != nil {
		fmt.Println(dbErr)
		return
	}

	http.HandleFunc("/", webServerHandler)
	http.HandleFunc("/save-my-name", saveName)
	http.ListenAndServe(":8182", nil)
}
