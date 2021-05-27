package config

import "database/sql"

func CreateDB() {
	// Open the database and create it if needed
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// Create users table in the database
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT, fewWords TEXT, age TEXT, address TEXT, photo TEXT)")
	statement.Exec()
	//
	// Open the database_post and create it if needed
	database_post, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database_post.Close()

	// Create post table in the database_post
	statement_post, _ := database_post.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT, body TEXT, like INTEGER, author TEXT, date TEXT, category TEXT, likedBy TEXT)")
	statement_post.Exec()
	//
	database_comment, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database_comment.Close()

	// Create users table in the database
	statement_comment, _ := database_comment.Prepare("CREATE TABLE IF NOT EXISTS comments (id INTEGER PRIMARY KEY, idMainPost TEXT, content TEXT, like INTEGER, author TEXT, date TEXT)")
	statement_comment.Exec()
}
