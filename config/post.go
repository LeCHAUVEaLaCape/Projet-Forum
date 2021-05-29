package config

import (
	"database/sql"
	"strings"

	. "./err"
)

func AddNewPost(title string, body string, dt string, data_newPost map[string]interface{}, category []string) {
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// add the inputs to the database
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare("INSERT INTO posts (title, body, author, date, category, like, likedBy, nbComments) VALUES (?, ?, ?, ?, ?, 0, '', 0)")
	CheckError(err)
	_, err = stmt.Exec(title, body, data_newPost["user"], dt, strings.Join(category, ""))
	CheckError(err)

	tx.Commit()
}
func Display_post(post_id string, data_post map[string]interface{}, body string) [7]string {
	var post [7]string
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	//range over database
	rows, _ := database.Query("SELECT id, title, body, author, date, like FROM posts WHERE id = ?", post_id)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&post[0], &post[1], &body, &post[3], &post[4], &post[6])
		CheckError(err)
	}
	// Remplace les \n par des <br> pour sauter des lignes en html
	post[2] = strings.Replace(body, string('\r'), "", -1)
	post[2] = strings.Replace(body, string('\n'), "<br>", -1)

	// Ajoute le chemin de la photo qui a été choisit par l'utilisateur
	rows, err := database.Query("SELECT photo FROM users WHERE username = ?", post[3])
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&post[5])
		CheckError(err)
	}
	err = rows.Err()
	CheckError(err)
	data_post["main_post"] = post
	return post
}
