package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	. "./config"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	accountCreated bool
	email          bool
	username       bool
}

var data []Data

const LocalhostPort = ":8000"

func main() {
	// Open the database and create it if needed
	database, _ := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, email TEXT, uuid TEXT, password TEXT)")
	statement.Exec()

	fmt.Println("Please connect to http://localhost:8000")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.HandleFunc("/", index)
	http.HandleFunc("/signIn", signIn)
	err := http.ListenAndServe(LocalhostPort, nil) // Set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Generate the main page when first loading the site
func index(w http.ResponseWriter, r *http.Request) {
	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "index", data)
}

// Generate the sign in page
func signIn(w http.ResponseWriter, r *http.Request) {
	info_creation_account := map[string]bool{
		"accountCreated": false,
		"username_used":  false,
		"email_used":     false,
	}

	// Open the database
	database, _ := sql.Open("sqlite3", "./users.db")
	defer database.Close()

	// Get the inputs
	input_username := r.FormValue("username")
	input_email := r.FormValue("email")
	input_password := r.FormValue("password")

	// Crypt the password
	hash, _ := HashPassword(input_password) // ./config/HashPassword.go

	// Call the function to add users
	AddUser(database, input_username, input_email, hash, info_creation_account) // ./config/AddUser.go

	t := template.New("signIn-template")
	t = template.Must(t.ParseFiles("./html/signIn.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "signIn", info_creation_account)
}
