package config

import (
	"database/sql"
	"net/http"
	"strconv"

	. "./err"
)

func GetInfoUser(w http.ResponseWriter, r *http.Request, data_user map[string]interface{}) {
	user := r.FormValue("user")

	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	var id, fewWords, address, age, photo string
	rows, _ := database.Query("SELECT id, username, email, fewWords, address, age, photo FROM users")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id, &username, &email, &fewWords, &address, &age, &photo)
		// Si l'input username est trouvé
		if user == username {
			data_user["username"] = username
			data_user["email"] = email
			data_user["fewWords"] = fewWords
			data_user["address"] = address
			data_user["age"] = age
			data_user["photo"] = photo
			break
		}
	}
	rows.Close()
}

func UpdateInfoUsers(w http.ResponseWriter, r *http.Request, id int) {
	// get input age/address/FewWords of the user
	add_few_words := r.FormValue("addFewWords")
	add_age := r.FormValue("age")
	add_address := r.FormValue("address")
	change_photo := r.FormValue("photo")
	var state string
	// add the input to his data in the DB
	if add_few_words != "" {
		state = "fewWords"
		UpdateInfoUser(add_few_words, state, id)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}
	if add_address != "" {
		state = "address"
		UpdateInfoUser(add_address, state, id)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}
	if add_age != "" {
		state = "age"
		UpdateInfoUser(add_age, state, id)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}
	if change_photo != "" {
		state = "photo"
		UpdateInfoUser(change_photo, state, id)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}
}

func UpdateInfoUser(input string, state string, id int) {
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	tx, err := database.Begin()
	CheckError(err)
	query := "UPDATE users SET " + state + " = ? WHERE id = " + strconv.Itoa(id)
	stmt, err := tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec(input)
	CheckError(err)
	tx.Commit()
}

func ModifPostAndComment(modif_post string, id_mainPost string, modif_comment string, id_comment string) {
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// Modify the body of the main post
	tx, err := database.Begin()
	CheckError(err)
	var query string
	if modif_post != "" {
		query = "UPDATE posts SET body = ? WHERE id = " + id_mainPost
	} else if modif_comment != "" {
		query = "UPDATE comments SET content = ? WHERE id = " + id_comment
	}
	stmt, err := tx.Prepare(query)
	CheckError(err)
	if modif_post != "" {
		_, err = stmt.Exec(modif_post)
	} else if modif_comment != "" {
		_, err = stmt.Exec(modif_comment)
	}
	CheckError(err)
	tx.Commit()
}
