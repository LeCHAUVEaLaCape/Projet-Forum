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
	rows, err := database.Query("SELECT username, email FROM users")
	CheckError(err)
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
	rows.Close()

	// Compte le nombre de ligne dans la BDD users
	var count, query string
	rows, err = database.Query("SELECT COUNT(*) FROM users")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&count)
	}
	// Si aucun utilisateur n'existe, le 1er créé sera automatiquement admin
	if count == "0" {
		query = "INSERT INTO users (username, email, password, fewWords, age, address, photo, notification, role) VALUES (?, ?, ?, '', '', '', '../assets/images/default.png', '', 'admin')"
	} else {
		query = "INSERT INTO users (username, email, password, fewWords, age, address, photo, notification, role) VALUES (?, ?, ?, '', '', '', '../assets/images/default.png', '', 'user')"
	}

	// add the inputs to the database
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec(input_username, input_email, input_password)
	CheckError(err)
	tx.Commit()
	info["accountCreated"] = true
}

// Envoie le role de l'utilisateur connecté à la page
func GetRole(data_info map[string]interface{}, on_user_page bool, user_page string) {
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	var role string
	// Prend le role de la personne connecté
	rows, _ := database.Query("SELECT role FROM users WHERE username = ?", data_info["user"])
	for rows.Next() {
		rows.Scan(&role)
	}
	rows.Close()
	data_info["role"] = role

	// Lorsque la page actuelle est le profil d'un utilisateur : Prend le role de cet utilisateur
	if on_user_page && user_page != "" {
		rows, _ := database.Query("SELECT role FROM users WHERE username = ?", user_page)
		for rows.Next() {
			rows.Scan(&role)
		}
		rows.Close()
		data_info["roleUserPage"] = role
	}
}
