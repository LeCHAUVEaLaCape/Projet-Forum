package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	// "strconv"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var data []string

const LocalhostPort = ":8000"

func main() {
	// statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, pseudo TEXT, email TEXT)")
	// statement.Exec()
	// statement, _ = database.Prepare("INSERT INTO users (pseudo, email) VALUES (?, ?)")
	// statement.Exec("svenickx", "bg@gmail.com")
	// rows, _ := database.Query("SELECT id, pseudo, email FROM users")
	// var id int
	// var pseudo string
	// var email string
	// for rows.Next() {
	// 	rows.Scan(&id, &pseudo, &email)
	// 	fmt.Println(strconv.Itoa(id) + ": " + pseudo + " " + email)
	// }
	//generate the main page when first loading the site
	database, _ := sql.Open("sqlite3", "./users.db")
	
	tmp := func (w http.ResponseWriter, r *http.Request) {
		inputpseudo := r.FormValue("pseudo")
		inputemail := r.FormValue("email")
		// fmt.Println( inputpseudo , inputemail)
		addUser( database ,inputpseudo , inputemail)
		t := template.New("verif-template")
		t = template.Must(t.ParseFiles("./html/verif.html", "./html/header&footer.html"))
		t.ExecuteTemplate(w, "verif", data)
	}
	fmt.Println("Please connect to http://localhost:8000")

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.HandleFunc("/", index)
	http.HandleFunc("/signIn", signIn)
	http.HandleFunc("/verif", tmp)
	err := http.ListenAndServe(LocalhostPort, nil) // Set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//generate the main page when first loading the site
func index(w http.ResponseWriter, r *http.Request) {
	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "index", data)
}

//generate the main page when first loading the site
func signIn(w http.ResponseWriter, r *http.Request) {
	t := template.New("signIn-template")
	t = template.Must(t.ParseFiles("./html/signIn.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "signIn", data)
}
func addUser(db *sql.DB, pseudo string, email string) {
    tx, _ := db.Begin()
    stmt, _ := tx.Prepare("insert into users (pseudo,email) values (?,?)")
    _, err := stmt.Exec(pseudo, email)
    checkError(err)
    tx.Commit()
}


func checkError(err error) {
    if err != nil {
        panic(err)
    }

    // catch to error.

}

