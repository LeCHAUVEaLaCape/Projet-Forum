package config

import (
	"database/sql"
	"net/http"
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
}

func UpdateInfoUsers(w http.ResponseWriter, r *http.Request) {
	// get input age/address/FewWords of the user
	add_few_words := r.FormValue("addFewWords")
	add_age := r.FormValue("age")
	add_address := r.FormValue("address")
	change_photo := r.FormValue("photo")
	user := r.FormValue("user")

	// add the input to his data in the DB
	if add_few_words != "" {
		UpdateInfoUser(add_few_words, "fewWords", user)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}
	if add_address != "" {
		UpdateInfoUser(add_address, "address", user)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}
	if add_age != "" {
		UpdateInfoUser(add_age, "age", user)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}
	if change_photo != "" {
		UpdateInfoUser(change_photo, "photo", user)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}
}

func UpdateInfoUser(input string, state string, user string) {
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")

	tx, err := database.Begin()
	CheckError(err)
	query := "UPDATE users SET " + state + " = ? WHERE username = ?"
	stmt, err := tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec(input, user)
	CheckError(err)
	tx.Commit()
	database.Close()
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

// Supprime le compte et remplace le auteurs des posts et commentaires par "user deleted"
func DelAccount(delete_account string) {
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	tx, err := database.Begin()
	CheckError(err)
	// Supprimer la ligne de l'utilisateur dans la BDD users
	stmt, err := tx.Prepare("DELETE FROM pendingForModerator WHERE username = ?")
	CheckError(err)
	_, err = stmt.Exec(delete_account)
	CheckError(err)
	// changer l'auteur dans la BDD users
	stmt, err = tx.Prepare("UPDATE posts SET author = ? WHERE author = ?")
	CheckError(err)
	_, err = stmt.Exec("user deleted", delete_account)
	CheckError(err)
	// changer l'auteur dans la BDD comments
	stmt, err = tx.Prepare("UPDATE comments SET author = ? WHERE author = ?")
	CheckError(err)
	_, err = stmt.Exec("user deleted", delete_account)
	CheckError(err)
	// changer l'auteur dans la BDD pendingPosts
	stmt, err = tx.Prepare("UPDATE pendingPosts SET author = ? WHERE author = ?")
	CheckError(err)
	_, err = stmt.Exec("user deleted", delete_account)
	CheckError(err)
	// Supprimer la ligne de l'utilisateur dans la BDD users
	stmt, err = tx.Prepare("DELETE FROM users WHERE username = ?")
	CheckError(err)
	_, err = stmt.Exec(delete_account)
	CheckError(err)

	tx.Commit()
}

// supprime un commentaire
func DelComment(w http.ResponseWriter, r *http.Request) {
	var nbComment int
	delete_comment := r.FormValue("delComment")
	idMainPost := r.FormValue("id-mainPost")
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	tx, err := database.Begin()
	CheckError(err)

	// remove 1 to the nb of comments when delete post
	rows, err := database.Query("SELECT nbComments FROM posts WHERE id = " + idMainPost)
	CheckError(err)
	for rows.Next() {
		rows.Scan(&nbComment)
	}
	nbComment -= 1
	rows.Close()
	stmt, err := tx.Prepare("UPDATE posts SET nbComments = ? WHERE id = ?")
	CheckError(err)
	_, err = stmt.Exec(nbComment, idMainPost)
	CheckError(err)

	// DELETE the comments of the main post
	stmt, err = tx.Prepare("DELETE FROM comments WHERE id = ?")
	CheckError(err)
	_, err = stmt.Exec(delete_comment)
	CheckError(err)
	tx.Commit()
	http.Redirect(w, r, "/post?id="+idMainPost, http.StatusSeeOther)
}

func DelPost(w http.ResponseWriter, r *http.Request) {
	delete_post := r.FormValue("delPost")
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	// delete the image
	Delete_image("mainPost", delete_post)

	// DELETE the comments of the main post
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare("DELETE FROM comments WHERE idMainPost = ?")
	CheckError(err)
	_, err = stmt.Exec(delete_post)
	CheckError(err)

	// DELETE the main POST
	stmt, err = tx.Prepare("DELETE FROM posts WHERE id = ?")
	CheckError(err)
	_, err = stmt.Exec(delete_post)
	CheckError(err)
	tx.Commit()

	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

// Supprime ou déplace le post s'il a été accepté ou non
func PostAcceptedOrNot(post_accepted string, id_pendingPost string) {
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	tx, err := database.Begin()
	// Accepter le post et le transferer jusqu'a la BDD posts
	if post_accepted == "true" {
		var transfer_post [10]string
		rows, err := database.Query("SELECT title, body, like, author, date, category, likedBy, nbComments, dislikedBy, image FROM pendingPosts WHERE id = " + id_pendingPost)
		CheckError(err)
		for rows.Next() {
			rows.Scan(&transfer_post[0], &transfer_post[1], &transfer_post[2], &transfer_post[3], &transfer_post[4], &transfer_post[5], &transfer_post[6], &transfer_post[7], &transfer_post[8], &transfer_post[9])
		}
		rows.Close()
		// add the inputs to the database
		stmt, err := tx.Prepare("INSERT INTO posts (title, body, like, author, date, category, likedBy, nbComments, dislikedBy, image, dislike) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)")
		CheckError(err)
		_, err = stmt.Exec(transfer_post[0], transfer_post[1], transfer_post[2], transfer_post[3], transfer_post[4], transfer_post[5], transfer_post[6], transfer_post[7], &transfer_post[8], &transfer_post[9])
		CheckError(err)
	} else {
		state := "pendingPosts"
		Delete_image(state, id_pendingPost)
	}
	// DELETE the comments of the main post
	CheckError(err)
	stmt, err := tx.Prepare("DELETE FROM pendingPosts WHERE id = ?")
	CheckError(err)
	_, err = stmt.Exec(id_pendingPost)
	CheckError(err)
	tx.Commit()
}

// Modifie le role d'un utilisateur quand un Admin le veut
func ChangeRole(w http.ResponseWriter, r *http.Request) {
	newRole := r.FormValue("change-role")
	user_to_modify_role := r.FormValue("username")

	if newRole != "" && user_to_modify_role != "" {
		// Open the database
		database, err := sql.Open("sqlite3", "./db-sqlite.db")
		CheckError(err)

		tx, err := database.Begin()
		CheckError(err)
		query := "UPDATE users SET role = ? WHERE username = ?"
		stmt, err := tx.Prepare(query)
		CheckError(err)
		_, err = stmt.Exec(newRole, user_to_modify_role)
		CheckError(err)
		tx.Commit()
		database.Close()
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}
