package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var data []string
var id int
var username, email string

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
	http.HandleFunc("/verif", verif)
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
	t := template.New("signIn-template")
	t = template.Must(t.ParseFiles("./html/signIn.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "signIn", data)
}

// Fill the database with the input of the users
func addUser(db *sql.DB, input_pseudo string, input_email string, input_password string) {

	// create UUID
	u1, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}

	tx, _ := db.Begin()

	// range over the database and check if there is double username/email
	rows, _ := db.Query("SELECT id, username, email, uuid FROM users")
	for rows.Next() {
		rows.Scan(&id, &username, &email)
		if username == input_pseudo || email == input_email {
			tx.Commit()
			//stop the function if a double is found
			return
		}
	}

	// add the inputs to the database and a UUID
	stmt, _ := tx.Prepare("insert into users (username, email, password, uuid) values (?, ?, ?, ?)")
	_, err = stmt.Exec(input_pseudo, input_email, input_password, u1)
	checkError(err)
	tx.Commit()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func verif(w http.ResponseWriter, r *http.Request) {
	// Open the database
	database, _ := sql.Open("sqlite3", "./users.db")
	defer database.Close()

	// Get the inputs
	input_pseudo := r.FormValue("username")
	input_email := r.FormValue("email")
	input_password := r.FormValue("password")

	// Crypt the password
	hash, _ := HashPassword(input_password)

	// call the function to add users
	addUser(database, input_pseudo, input_email, hash)

	t := template.New("verif-template")
	t = template.Must(t.ParseFiles("./html/verif.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "verif", data)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
