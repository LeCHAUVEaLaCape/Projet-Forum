package config

import (
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

var username, email string

// Fill the database with the input of the users
func AddUser(db *sql.DB, input_username string, input_email string, input_password string, info map[string]bool) {

	if input_email == "" || input_username == "" || input_password == "" {
		return
	}

	tx, _ := db.Begin()
	// range over the database and check if there is double username/email
	rows, _ := db.Query("SELECT username, email FROM users")
	for rows.Next() {
		rows.Scan(&username, &email)
		//stop the function if a double is found
		if username == input_username {
			info["username_used"] = true
			fmt.Println("username utilisé")
			return
		}
		if email == input_email {
			info["email_used"] = true
			fmt.Println("email utilisé")
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
	stmt, _ := tx.Prepare("insert into users (username, email, password, uuid) values (?, ?, ?, ?)")
	_, err = stmt.Exec(input_username, input_email, input_password, u1)
	if err != nil {
		panic(err)
	}
	tx.Commit()
	fmt.Println("account created")
	info["accountCreated"] = true
}
