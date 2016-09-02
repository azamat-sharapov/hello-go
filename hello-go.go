package main

import (
	"flag"
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
		name TEXT UNIQUE
	)
	`
	pguser := flag.String("pguser", "", "Enter postgres user login")
	pgpass := flag.String("pgpass", "", "Enter postgres user password")

	flag.Parse()

	if *pguser == "" || *pgpass == "" {
		return fmt.Errorf("Usage: \n  hello-go -pguser [postgres user] -pgpass [postgres user password]")
	}

	db, err = sqlx.Connect("postgres", fmt.Sprintf("user=%v password=%v dbname=hello-go sslmode=disable", *pguser, *pgpass))

	if err == nil {
		db.MustExec(schema)
	}

	return err
}

func saveName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	person := Person{r.PostForm.Get("name")}
	var message string
	var err error
	var rows *sqlx.Rows

	rows, err = db.NamedQuery("SELECT * from people where name = :name", person)

	if rows.Next() {
		message = fmt.Sprintf("Welcome back %s!", person.Name)
	} else {
		_, err = db.NamedExec("INSERT INTO people (name) VALUES (:name)", person)
		message = fmt.Sprintf("Hello %s! Now you are registered in our database.", person.Name)
	}

	if err != nil {
		showError(w, err)
		return
	}

	fmt.Fprintf(w, "{\"message\": \"%s\"}", message)
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

	fmt.Println("Server about to run on: http://localhost:8182")
	http.ListenAndServe(":8182", nil)
}
