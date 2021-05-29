package config

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	. "./err"
)

func Like(change_nmb_like string, data_post map[string]interface{}, post_id string, w http.ResponseWriter, r *http.Request) string {
	// Verifie si l'utilisateur a liké
	var likedBy string
	data_post["already_liked"], likedBy = CheckIfLikedByUser(post_id, data_post)

	if change_nmb_like == "1" {
		likedBy = AddLike(post_id, data_post, likedBy)
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
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// créer un string des personnes qui ont liké
	var likedBy string
	rows, _ := database.Query("SELECT likedBy FROM posts WHERE id = ?", post_id)
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
func AddLike(post_id string, data_post map[string]interface{}, likedBy string) string {
	var like int
	user := data_post["user"].(string)

	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	tx, err := database.Begin()
	// Ajoute +1 like des POSTS
	rows, _ := database.Query("SELECT like FROM posts WHERE id = ?", post_id)
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
	query = "UPDATE posts SET likedBy = ? WHERE id = " + post_id
	stmt, err = tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec(likedBy)
	CheckError(err)
	tx.Commit()

	return likedBy
}

// remove a like on click
func RemoveLike(post_id string, data_post map[string]interface{}, likedBy string) string {
	var like int
	user := data_post["user"].(string)

	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	tx, err := database.Begin()
	// delete a like
	rows, _ := database.Query("SELECT like FROM posts WHERE id = ?", post_id)
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
