package config

import (
	"database/sql"
	"fmt"
)

func AddNewPost(title string, body string, typeOfPost string, dt string, data_newPost map[string]interface{}) {
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// add the inputs to the database
	tx, err := database.Begin()
	if err != nil {
		fmt.Println(err)
	}
	stmt, err := tx.Prepare("INSERT INTO posts (title, body, type, author, date) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(title, body, typeOfPost, data_newPost["user"], dt)
	if err != nil {
		fmt.Println(err)
	} else {
		tx.Commit()
	}
}
