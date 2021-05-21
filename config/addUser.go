package config

import (
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

var username, email string

// Fill the database with the input of the users
func AddUser(db *sql.DB, input_username string, input_email string, input_password string, info map[string]bool) {
	tx, _ := db.Begin()
	// range over the database and check if there is double username/email
	rows, _ := db.Query("SELECT username, email FROM users")
	for rows.Next() {
		rows.Scan(&username, &email)
		// fmt.Println("input : ", input_email, input_username)
		// fmt.Println("data : ", email, username)
		// fmt.Println("-----------")
		//stop the function if a double is found
		if username == input_username {
			tx.Commit()
			info["username_used"] = true
			fmt.Println("Username :", input_username+" déjà utilisé")
			return
		}
		if email == input_email {
			tx.Commit()
			info["email_used"] = true
			fmt.Println("Email :", input_email+" déjà utilisé")
			return
		}
	}

	// create UUID
	u1, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}

	// add the inputs to the database and a UUID
	stmt, _ := tx.Prepare("INSERT INTO USERS (username, email, password, uuid) VALUES (?, ?, ?, ?)")
	_, err = stmt.Exec(input_username, input_email, input_password, u1)
	if err != nil {
		panic(err)
	}
	tx.Commit()
	fmt.Println("Account created : ", input_username+", "+input_email)
	info["accountCreated"] = true
}
