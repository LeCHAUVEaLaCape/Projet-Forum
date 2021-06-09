package config

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
)

type Notif struct {
	idPost          string
	liker           string
	action          string
	userToSendNotif string
}

func Like(change_nmb_like string, data_post map[string]interface{}, post_id string, w http.ResponseWriter, r *http.Request) string {
	// Verifie si l'utilisateur a liké
	var likedBy string
	data_post["already_liked"], likedBy = CheckIfLikedByUser(post_id, data_post)

	var notif Notif
	notif.idPost = r.FormValue("id-post")
	notif.liker = r.FormValue("liker")
	notif.action = r.FormValue("action")
	notif.userToSendNotif = r.FormValue("userToSendNotif")

	if change_nmb_like == "1" {
		likedBy = AddLike(post_id, data_post, likedBy, notif)
		http.Redirect(w, r, "/post?id="+post_id, http.StatusSeeOther)
	} else if change_nmb_like == "0" {
		likedBy = RemoveLike(post_id, data_post, likedBy)
		http.Redirect(w, r, "/post?id="+post_id, http.StatusSeeOther)
	}

	return likedBy
}

// verifie Si l'utilisateur connecter à déjà liker le post
func CheckIfLikedByUser(post_id string, data_post map[string]interface{}) (bool, string) {
	var user string
	if data_post["user"] != nil {
		user = data_post["user"].(string)
	}

	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	// créer un string des personnes qui ont liké
	var likedBy string
	rows, err := database.Query("SELECT likedBy FROM posts WHERE id = ?", post_id)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&likedBy)
		CheckError(err)
	}

	if likedBy != "" {
		// Split le string en array
		all_likedBy := strings.Split(likedBy, " ")
		// parcour l'array de ceux qui ont liké pour éviter les doublons
		for i := 0; i < len(all_likedBy); i++ {
			if all_likedBy[i] == user {
				return false, likedBy
			}
		}
	}

	return true, likedBy
}

// ajoute un like au post quand l'utilisateur clique sur le bouton
func AddLike(post_id string, data_post map[string]interface{}, likedBy string, notif Notif) string {
	var like int
	user := data_post["user"].(string)

	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)

	tx, err := database.Begin()
	CheckError(err)
	// Ajoute +1 like des POSTS
	rows, err := database.Query("SELECT like FROM posts WHERE id = ?", post_id)
	CheckError(err)
	for rows.Next() {
		err := rows.Scan(&like)
		CheckError(err)
	}
	like += 1

	// Update the nb of like
	query := "UPDATE posts SET like = " + strconv.Itoa(like) + " WHERE id = " + post_id
	stmt, err := tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec()
	CheckError(err)

	if likedBy == "" {
		likedBy = user
	} else {
		likedBy += " " + user
	}

	// Update the users who liked
	query = "UPDATE posts SET likedBy =  ? WHERE id = " + post_id
	stmt, err = tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec(likedBy)
	CheckError(err)
	tx.Commit()
	rows.Close()
	database.Close()

	UpdateNotif(notif)
	return likedBy
}

// remove a like on click
func RemoveLike(post_id string, data_post map[string]interface{}, likedBy string) string {
	var like int
	user := data_post["user"].(string)

	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	tx, err := database.Begin()
	CheckError(err)
	// delete a like
	rows, err := database.Query("SELECT like FROM posts WHERE id = ?", post_id)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&like)
		CheckError(err)
	}
	like -= 1

	// Update the nb of like
	query := "UPDATE posts SET like = " + strconv.Itoa(like) + " WHERE id = " + post_id
	stmt, err := tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec()
	CheckError(err)

	all_likedBy := strings.Split(likedBy, " ")

	for i := 0; i < len(all_likedBy); i++ {
		if user == all_likedBy[i] {
			all_likedBy = append(all_likedBy[:i], all_likedBy[i+1:]...)
			break
		}
	}
	likedBy = strings.Join(all_likedBy, " ")

	// Update the users who liked
	query = "UPDATE posts SET likedBy = ? WHERE id = " + post_id
	stmt, err = tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec(likedBy)
	CheckError(err)
	tx.Commit()

	return likedBy
}

func UpdateNotif(notif Notif) {
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	var notification string
	// Get the notification from the user to add the new notif
	rows, err := database.Query("SELECT notification FROM users WHERE username = ?", notif.userToSendNotif)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&notification)
		CheckError(err)
	}

	tx, err := database.Begin()
	CheckError(err)
	// Update the notif section
	if notif.liker != notif.userToSendNotif {
		query := "UPDATE users SET notification = ? WHERE username = ?"
		stmt, err := tx.Prepare(query)
		CheckError(err)
		_, err = stmt.Exec(notification+notif.idPost+" "+notif.liker+" "+notif.action+",", notif.userToSendNotif)
		CheckError(err)
	}

	tx.Commit()
}
func DisLike(change_nmb_dislike string, data_post map[string]interface{}, post_id string, w http.ResponseWriter, r *http.Request) string {
	// Verifie si l'utilisateur a disliké
	var dislikedBy string
	data_post["already_disliked"], dislikedBy = CheckIfdisLikedByUser(post_id, data_post)

	var notif Notif
	notif.idPost = r.FormValue("id-post")
	notif.liker = r.FormValue("disliker")
	notif.action = r.FormValue("action")
	notif.userToSendNotif = r.FormValue("userToSendNotif")

	if change_nmb_dislike == "1" {
		dislikedBy = AddDisLike(post_id, data_post, dislikedBy, notif)
		http.Redirect(w, r, "/post?id="+post_id, http.StatusSeeOther)
	} else if change_nmb_dislike == "0" {
		dislikedBy = RemoveDisLike(post_id, data_post, dislikedBy)
		http.Redirect(w, r, "/post?id="+post_id, http.StatusSeeOther)
	}

	return dislikedBy
}

// verifie Si l'utilisateur connecter à déjà disliker le post
func CheckIfdisLikedByUser(post_id string, data_post map[string]interface{}) (bool, string) {
	var user string
	if data_post["user"] != nil {
		user = data_post["user"].(string)
	}

	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	// créer un string des personnes qui ont liké
	var dislikedBy string
	rows, err := database.Query("SELECT dislikedBy FROM posts WHERE id = ?", post_id)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&dislikedBy)
		CheckError(err)
	}

	if dislikedBy != "" {
		// Split le string en array
		all_dislikedBy := strings.Split(dislikedBy, " ")
		// parcour l'array de ceux qui ont disliké pour éviter les doublons
		for i := 0; i < len(all_dislikedBy); i++ {
			if all_dislikedBy[i] == user {
				return false, dislikedBy
			}
		}
	}

	return true, dislikedBy
}

// ajoute un dislike au post quand l'utilisateur clique sur le bouton
func AddDisLike(post_id string, data_post map[string]interface{}, dislikedBy string, notif Notif) string {
	var dislike int
	user := data_post["user"].(string)

	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)

	tx, err := database.Begin()
	CheckError(err)
	// Ajoute +1 disslike des POSTS
	rows, err := database.Query("SELECT dislike FROM posts WHERE id = ?", post_id)
	CheckError(err)
	for rows.Next() {
		err := rows.Scan(&dislike)
		CheckError(err)
	}
	dislike += 1

	// Update the nb of dislike
	query := "UPDATE posts SET dislike = " + strconv.Itoa(dislike) + " WHERE id = " + post_id
	stmt, err := tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec()
	CheckError(err)

	if dislikedBy == "" {
		dislikedBy = user
	} else {
		dislikedBy += " " + user
	}

	// Update the users who disliked
	query = "UPDATE posts SET dislikedBy =  ? WHERE id = " + post_id
	stmt, err = tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec(dislikedBy)
	CheckError(err)
	tx.Commit()
	rows.Close()
	database.Close()

	UpdateNotif(notif)
	return dislikedBy
}

// remove a dislike on click
func RemoveDisLike(post_id string, data_post map[string]interface{}, dislikedBy string) string {
	var dislike int
	user := data_post["user"].(string)

	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	tx, err := database.Begin()
	CheckError(err)
	// delete a dislike
	rows, err := database.Query("SELECT dislike FROM posts WHERE id = ?", post_id)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&dislike)
		CheckError(err)
	}
	dislike -= 1

	// Update the nb of like
	query := "UPDATE posts SET dislike = " + strconv.Itoa(dislike) + " WHERE id = " + post_id
	stmt, err := tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec()
	CheckError(err)

	all_dislikedBy := strings.Split(dislikedBy, " ")

	for i := 0; i < len(all_dislikedBy); i++ {
		if user == all_dislikedBy[i] {
			all_dislikedBy = append(all_dislikedBy[:i], all_dislikedBy[i+1:]...)
			break
		}
	}
	dislikedBy = strings.Join(all_dislikedBy, " ")

	// Update the users who liked
	query = "UPDATE posts SET dislikedBy = ? WHERE id = " + post_id
	stmt, err = tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec(dislikedBy)
	CheckError(err)
	tx.Commit()

	return dislikedBy
}
