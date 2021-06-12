package config

import (
	"database/sql"
)

func CreateDB() {
	// Open the database and create it if needed
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	// Create users table in the database
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT, fewWords TEXT, age TEXT, address TEXT, photo TEXT, notification TEXT, role TEXT)")
	CheckError(err)
	statement.Exec()

	// Create post table in the database_post
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT, body TEXT, like INTEGER, dislike INTEGER,author TEXT, date TEXT, category TEXT, likedBy TEXT, dislikedBy TEXT, nbComments INTEGER, image TEXT)")
	CheckError(err)
	statement.Exec()

	// Create comments table in the database
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS comments (id INTEGER PRIMARY KEY, idMainPost TEXT, content TEXT, like INTEGER, author TEXT, date TEXT)")
	CheckError(err)
	statement.Exec()

	// Create pendingPosts table in the database
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS pendingPosts (id INTEGER PRIMARY KEY, title TEXT, body TEXT, like INTEGER, author TEXT, date TEXT, category TEXT, likedBy TEXT, dislikedBy TEXT, nbComments INTEGER, image TEXT)")
	CheckError(err)
	statement.Exec()
	// Create pending table in the database
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS pendingForModerator (username TEXT)")
	CheckError(err)
	statement.Exec()
	// Create report table in the database
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS report (nameUser TEXT, reasonReport TEXT, answerReport TEXT, reported TEXT, nameReporter TEXT)")
	CheckError(err)
	statement.Exec()
	// Create categories table in the database
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY, categorie TEXT, color TEXT)")
	CheckError(err)
	statement.Exec()
}
