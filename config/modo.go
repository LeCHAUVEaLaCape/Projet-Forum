package config

import (
	"database/sql"
	"net/http"

	erreur "./err"
)

type dataUsers struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
type dataForDashboard struct {
	Admins     []dataUsers
	Moderators []dataUsers
	Users      []dataUsers
}

// modo.go
func DisplayAdminModo(data *map[string]interface{}) {
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	rows, _ := database.Query("SELECT id, username, role FROM users")
	defer rows.Close()
	var tmp dataForDashboard
	for rows.Next() {
		var data_user dataUsers
		rows.Scan(&data_user.Id, &data_user.Username, &data_user.Role)
		if data_user.Role == "user" {
			tmp.Users = append(tmp.Users, data_user)
		} else if data_user.Role == "moderator" {
			tmp.Moderators = append(tmp.Moderators, data_user)
		} else if data_user.Role == "admin" {
			tmp.Admins = append(tmp.Moderators, data_user)
		}
	}
	(*data)["liste"] = tmp
}

// modo.go
func ResquestForModo(r *http.Request) {
	username := r.FormValue("user")
	if username == "" {
		return
	}
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// range over the database and check if there is double username
	rows, err := database.Query("SELECT username FROM pendingForModerator")
	erreur.CheckError(err)
	for rows.Next() {
		tmp := ""
		rows.Scan(&tmp)
		//stop the function if a double is found
		if username == tmp {
			return
		}
	}
	rows.Close()
	// add the inputs to the database
	tx, err := database.Begin()
	erreur.CheckError(err)
	stmt, err := tx.Prepare("INSERT INTO pendingForModerator (username) VALUES (?)")
	erreur.CheckError(err)
	_, err = stmt.Exec(username)
	erreur.CheckError(err)
	tx.Commit()
}

// modo.go
func DisplayPendingForModo(data *map[string]interface{}) {
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	var res []string
	// range over the database
	rows, err := database.Query("SELECT username FROM pendingForModerator")
	erreur.CheckError(err)
	for rows.Next() {
		tmp := ""
		rows.Scan(&tmp)
		res = append(res, tmp)
	}
	rows.Close()
	(*data)["liste_demande"] = res
}

//modo .go
func AccepterDemande(w http.ResponseWriter, r *http.Request) {
	user_accepter := r.FormValue("accepter")
	if user_accepter == "" {
		return
	}
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	tx, err := database.Begin()
	erreur.CheckError(err)
	stmt, err := tx.Prepare("UPDATE users SET role = ? WHERE username = ?")
	erreur.CheckError(err)
	_, err = stmt.Exec("moderator", user_accepter)
	erreur.CheckError(err)
	// Supprimer la ligne de l'utilisateur dans la BDD pendingForModerator
	stmt, err = tx.Prepare("DELETE FROM pendingForModerator WHERE username = ?")
	erreur.CheckError(err)
	_, err = stmt.Exec(user_accepter)
	erreur.CheckError(err)
	tx.Commit()
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

// modo.go
func RefuserDemande(w http.ResponseWriter, r *http.Request) {
	user_refuser := r.FormValue("refuser")
	if user_refuser == "" {
		return
	}
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	tx, err := database.Begin()
	erreur.CheckError(err)
	// Supprimer la ligne de l'utilisateur dans la BDD pendingForModerator
	stmt, err := tx.Prepare("DELETE FROM pendingForModerator WHERE username = ?")
	erreur.CheckError(err)
	_, err = stmt.Exec(user_refuser)
	erreur.CheckError(err)
	tx.Commit()
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}
