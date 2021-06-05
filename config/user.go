package config

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	. "./err"
)

var username, email string

// Fill the database with the input of the users
func AddUser(input_username string, input_email string, input_password string, info map[string]interface{}) {
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
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
			rows.Close()
			return
		} else if email == input_email {
			info["email_used"] = true
			fmt.Println("Email :", input_email+" déjà utilisé")
			rows.Close()
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
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	var role string
	// Prend le role de la personne connecté
	rows, err := database.Query("SELECT role FROM users WHERE username = ?", data_info["user"])
	CheckError(err)
	for rows.Next() {
		rows.Scan(&role)
	}
	rows.Close()
	data_info["role"] = role

	// Lorsque la page actuelle est le profil d'un utilisateur : Prend le role de cet utilisateur
	if on_user_page && user_page != "" {
		rows, err := database.Query("SELECT role FROM users WHERE username = ?", user_page)
		CheckError(err)
		for rows.Next() {
			rows.Scan(&role)
		}
		defer rows.Close()
		data_info["roleUserPage"] = role
	}
}

// Affiche les notifications de l'utilisateur s'il y en a
func CheckNotif(w http.ResponseWriter, r *http.Request, data_notif map[string]interface{}) {
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)

	var notification string
	var arrNotification []string
	var test []string
	var arr_notif [][]string
	rows, err := database.Query("SELECT notification FROM users WHERE username = ?", data_notif["user"])
	CheckError(err)
	for rows.Next() {
		err := rows.Scan(&notification)
		CheckError(err)
		// get the section splited with the ","
		arrNotification = strings.Split(notification, ",")
		// remove the last part of the section which is empty
		if len(arrNotification[len(arrNotification)-1]) < 1 {
			arrNotification = arrNotification[0 : len(arrNotification)-1]
		}
		// Now split with white spaces and add them to the final array
		for i := 0; i < len(arrNotification); i++ {
			test = strings.Split(arrNotification[i], " ")
			arr_notif = append(arr_notif, test)
		}
	}
	err = rows.Err()
	CheckError(err)

	data_notif["notif"] = arr_notif

	// Delete notif
	del_notif := r.FormValue("del-notif")
	userToDel := r.FormValue("user")
	if del_notif == "1" {
		tx, err := database.Begin()
		CheckError(err)
		stmt, err := tx.Prepare("UPDATE users SET notification = ? WHERE username = ?")
		CheckError(err)
		_, err = stmt.Exec("", userToDel)
		CheckError(err)
		tx.Commit()
		rows.Close()
		database.Close()
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
	}
	rows.Close()
	database.Close()
}
