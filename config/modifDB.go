package config

import (
	"database/sql"
	"strconv"

	. "./err"
)

func UpdateInfoUser(database *sql.DB, input string, state string, id int) {
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
