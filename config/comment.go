package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func Adding_comment(add_comment string, post *[6]string, user string) {
	database_comment, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database_comment.Close()
	dt := time.Now()
	// commentaire à ajouter
	tx, err := database_comment.Begin()
	if err != nil {
		fmt.Println(err)
	}
	statement, err := tx.Prepare("INSERT INTO comments (idMainPost, content, author, date, like) VALUES (?, ?, ?, ?, 0)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = statement.Exec((*post)[0], add_comment, user, dt.Format("02-01-2006 15:04:05"))
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()

}
func Display_comments(data_post map[string]interface{}, post_id string) {
	// commentaires
	var comments [][5]string
	database_comment, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database_comment.Close()
	//range over database
	rows_comment, _ := database_comment.Query("SELECT content, author, date FROM comments WHERE idMainPost = ?", post_id)
	defer rows_comment.Close()
	for rows_comment.Next() {
		var tmp [5]string
		err := rows_comment.Scan(&tmp[0], &tmp[1], &tmp[2])
		if err != nil {
			log.Fatal(err)
		}
		// Ajoute le chemin de la photo qui a été choisit par l'utilisateur
		rows, err := database_comment.Query("SELECT photo FROM users WHERE username = ?", tmp[1])
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&tmp[4])
			if err != nil {
				log.Fatal(err)
			}
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		comments = append(comments, tmp)
	}
	err := rows_comment.Err()
	if err != nil {
		log.Fatal(err)
	}

	data_post["comments"] = comments

}
