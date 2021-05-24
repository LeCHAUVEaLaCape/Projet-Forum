package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	. "./config"
	. "./cookies"

	_ "github.com/mattn/go-sqlite3"
)

var id int
var username, password, email, age string
var create_cookie, userFound = false, false

type singleUser struct {
	username string
	email    string
	age      string
}

var data = make(map[string]interface{})

func main() {
	data["user"] = ""
	// Open the database and create it if needed
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// Create users table in the database
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT, fewWords TEXT, age TEXT, address TEXT)")
	statement.Exec()

	fmt.Println("Please connect to http://localhost:8000")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.HandleFunc("/", index)
	http.HandleFunc("/SignUp", SignUp)
	http.HandleFunc("/logIn", logIn)
	http.HandleFunc("/logOut", logOut)
	http.HandleFunc("/welcome", welcome)
	http.HandleFunc("/allUsers", allUsers)
	http.HandleFunc("/user", user)
	err := http.ListenAndServe(":8000", nil) // Set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Generate the main page when first loading the site
func index(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_index := make(map[string]interface{})
	for k, v := range data {
		data_index[k] = v
	}
	data_index["cookieExist"] = false
	GetCookie(data_index, r)

	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "index", data_index)
}

// Generate the sign in page
func SignUp(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_SignUp := make(map[string]interface{})
	for k, v := range data {
		data_SignUp[k] = v
	}

	// Check if a cookie exist
	data_SignUp["cookieExist"] = false
	GetCookie(data_SignUp, r)

	// if a cookie already exist, redirect
	if data_SignUp["cookieExist"] == true {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	// intiate values of the map
	data_SignUp["accountCreated"] = false
	data_SignUp["username_used"] = false
	data_SignUp["email_used"] = false

	// Get the inputs
	input_username := r.FormValue("username")
	input_email := r.FormValue("email")
	input_password := r.FormValue("password")

	// Crypt the password
	hash := HashPassword(input_password) // ./config/HashPassword.go

	// Call the function to add users
	if input_email != "" || input_username != "" || input_password != "" {
		AddUser(input_username, input_email, hash, data_SignUp) // ./config/AddUser.go
	}

	t := template.New("SignUp-template")
	t = template.Must(t.ParseFiles("./html/SignUp.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "SignUp", data_SignUp)
}

func logIn(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_logIn := make(map[string]interface{})
	for k, v := range data {
		data_logIn[k] = v
	}
	data_logIn["cookieExist"] = false
	GetCookie(data_logIn, r)

	if data_logIn["cookieExist"] == true {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	// get user input to log in
	user_login := r.FormValue("user-login")
	password_login := r.FormValue("password-login")

	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	create_cookie = false
	// Parcourir la BDD
	rows, _ := database.Query("SELECT username, password FROM users")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&username, &password)
		// Si l'input username est trouvé
		if user_login == username {
			fmt.Println("username ok ", username)
			// Compare l'input password avec celui de la BDD
			if ComparePasswords(password, []byte(password_login)) {
				create_cookie = true
				data["user"] = username
				data["cookieExist"] = true
				break
			}
		}
	}

	if create_cookie {
		// Créé un cookie si user bien authentifié
		CreateCookie(w, r)
		data["cookieExist"] = true
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	} else {
		fmt.Println("connexion failed")
	}

	fmt.Println(data_logIn)
	t := template.New("logIn-template")
	t = template.Must(t.ParseFiles("./html/LogIn.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "LogIn", data_logIn)
}

// Generate the Welcome page (accessible only to logged in users)
func welcome(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_welcome := make(map[string]interface{})
	for k, v := range data {
		data_welcome[k] = v
	}

	GetCookie(data_welcome, r)

	t := template.New("welcome-template")
	t = template.Must(t.ParseFiles("./html/welcome.html", "./html/header&footer.html"))
	fmt.Println("data : ", data)
	fmt.Println("data_welcome : ", data_welcome)
	t.ExecuteTemplate(w, "welcome", data_welcome)
}

// Generate the user page
func user(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_user := make(map[string]interface{})
	for k, v := range data {
		data_user[k] = v
	}
	data_user["cookieExist"] = false
	GetCookie(data_user, r)

	user := r.FormValue("user")

	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// tx, _ := database.Begin()
	// Parcourir la BDD
	rows, _ := database.Query("SELECT username, email FROM users")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&username, &email)
		// Si l'input username est trouvé
		if user == username {
			userFound = true
			data_user["username"] = username
			data_user["email"] = email
			break
		}
	}

	// Check if the user logged is on his personnal page
	if data_user["username"] == data_user["user"] {
		data_user["sameUser"] = true
	} else {
		data_user["sameUser"] = false
	}
	fmt.Println(data_user)

	// get input age/address/FewWords of the user
	add_few_words := r.FormValue("addFewWords")
	add_age := r.FormValue("age")
	add_address := r.FormValue("address")
	// add the input to his data in the DB
	if add_few_words != "" || add_age != "" || add_address != "" {
		tx, err := database.Begin()
		if err != nil {
			fmt.Println(err)
		}
		stmt, err := database.Prepare("INSERT OR REPLACE INTO users(fewWords, age, address) values(?,?,?)")
		if err != nil {
			fmt.Println(err)
		}
		_, err = stmt.Exec(add_few_words, add_age, add_address)
		if err != nil {
			fmt.Println(err)
		}
		tx.Commit()
	}

	t := template.New("user-template")
	t = template.Must(t.ParseFiles("./html/user.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "user", data_user)
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	var aUser singleUser
	var all_users []singleUser
	data_allUsers := make(map[string]interface{})
	for k, v := range data {
		data_allUsers[k] = v
	}
	GetCookie(data_allUsers, r)

	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// Parcourir la BDD
	rows, _ := database.Query("SELECT username, email, age FROM users")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&username, &email, &age)
		aUser.username = username
		aUser.age = age
		aUser.email = email
		all_users = append(all_users, aUser)
	}
	data_allUsers["allUsers"] = all_users

	t := template.New("allUsers-template")
	t = template.Must(t.ParseFiles("./html/allUsers.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "allUsers", data_allUsers)
}

// Delete the cookie to disconnect the user
func logOut(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}
	DeleteCookie(w)
	http.Redirect(w, r, "/index", http.StatusSeeOther)
	delete(data, "user")
}
