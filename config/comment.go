package config

import (
	"database/sql"
	"net/http"
	"strings"
	"time"
)

func Adding_comment(w http.ResponseWriter, r *http.Request, add_comment string, post *[8]string, user string) {
	var notif Notif
	notif.idPost = r.FormValue("id-post")
	notif.liker = r.FormValue("liker") // liker = user commenting the post
	notif.action = r.FormValue("action")
	notif.userToSendNotif = r.FormValue("userToSendNotif")

	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)

	dt := time.Now()
	// commentaire à ajouter
	tx, err := database.Begin()
	CheckError(err)
	statement, err := tx.Prepare("INSERT INTO comments (idMainPost, content, author, date, like) VALUES (?, ?, ?, ?, 0)")
	CheckError(err)
	_, err = statement.Exec((*post)[0], add_comment, user, dt.Format("02-01-2006 15:04:05"))
	CheckError(err)

	// Augmente le nombre de like dans la BDD de post
	query := "UPDATE posts SET nbComments = nbComments + 1 WHERE id = " + (*post)[0]
	stmt, err := tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec()
	CheckError(err)
	tx.Commit()
	database.Close()

	UpdateNotif(notif) // config/like.go

}
func Display_comments(data_post map[string]interface{}, post_id string) {
	// commentaires
	var comments [][6]string
	var content string
	database_comment, err := sql.Open("sqlite3", "./db-sqlite.db")
	defer database_comment.Close()
	//range over database
	rows_comment, err := database_comment.Query("SELECT content, author, date, id FROM comments WHERE idMainPost = ?", post_id)
	CheckError(err)
	defer rows_comment.Close()
	for rows_comment.Next() {
		var tmp [6]string
		err := rows_comment.Scan(&content, &tmp[1], &tmp[2], &tmp[5])
		CheckError(err)
		// Remplace les \n par des <br> pour sauter des lignes en html
		tmp[0] = strings.Replace(content, string('\r'), "", -1)
		tmp[0] = strings.Replace(content, string('\n'), "<br>", -1)

		// Ajoute le chemin de la photo qui a été choisit par l'utilisateur
		rows, err := database_comment.Query("SELECT photo FROM users WHERE username = ?", tmp[1])
		CheckError(err)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&tmp[4])
			CheckError(err)
		}
		err = rows.Err()
		CheckError(err)
		comments = append(comments, tmp)
	}
	err = rows_comment.Err()
	CheckError(err)

	data_post["comments"] = comments

}
