package config

import (
	"database/sql"
	"fmt"

	. "./err"
)

var username, email string

// Fill the database with the input of the users
func AddUser(input_username string, input_email string, input_password string, info map[string]interface{}) {
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// range over the database and check if there is double username/email
	rows, _ := database.Query("SELECT username, email FROM users")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&username, &email)
		//stop the function if a double is found
		if username == input_username {
			info["username_used"] = true
			fmt.Println("Username :", input_username+" déjà utilisé")
			return
		} else if email == input_email {
			info["email_used"] = true
			fmt.Println("Email :", input_email+" déjà utilisé")
			return
		}
	}

	// add the inputs to the database
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare("INSERT INTO users (username, email, password, fewWords, age, address, photo) VALUES (?, ?, ?, '', '', '', '../assets/images/default.png')")
	CheckError(err)
	_, err = stmt.Exec(input_username, input_email, input_password)
	CheckError(err)
	tx.Commit()
	info["accountCreated"] = true

}
