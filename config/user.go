package config

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var id int
var username, email, photo, likedBy string

// Fill the database with the input of the users
func AddUser(input_username string, input_email string, input_password string, info map[string]interface{}) {
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	// range over the database and check if there is double username/email
	rows, err := database.Query("SELECT username, email FROM users")
	CheckError(err)
	for rows.Next() {
		rows.Scan(&username, &email)
		//stop the function if a double is found
		if username == input_username {
			info["username_used"] = true
			fmt.Println("Username :", input_username+" déjà utilisé")
			rows.Close()
			return
		} else if email == input_email {
			info["email_used"] = true
			fmt.Println("Email :", input_email+" déjà utilisé")
			rows.Close()
			return
		}
	}
	rows.Close()

	// Compte le nombre de ligne dans la BDD users
	var count, query string
	rows, err = database.Query("SELECT COUNT(*) FROM users")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&count)
	}
	// Si aucun utilisateur n'existe, le 1er créé sera automatiquement admin
	if count == "0" {
		query = "INSERT INTO users (username, email, password, fewWords, age, address, photo, notification, role) VALUES (?, ?, ?, '', '', '', '../assets/images/default.png', '', 'admin')"
	} else {
		query = "INSERT INTO users (username, email, password, fewWords, age, address, photo, notification, role) VALUES (?, ?, ?, '', '', '', '../assets/images/default.png', '', 'user')"
	}

	// add the inputs to the database
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare(query)
	CheckError(err)
	_, err = stmt.Exec(input_username, input_email, input_password)
	CheckError(err)
	tx.Commit()
	info["accountCreated"] = true
}

// Envoie le role de l'utilisateur connecté à la page
func GetRole(data_info map[string]interface{}, on_user_page bool, user_page string) {
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	var role string
	// Prend le role de la personne connecté
	rows, err := database.Query("SELECT role FROM users WHERE username = ?", data_info["user"])
	CheckError(err)
	for rows.Next() {
		rows.Scan(&role)
	}
	rows.Close()
	data_info["role"] = role

	// Lorsque la page actuelle est le profil d'un utilisateur : Prend le role de cet utilisateur
	if on_user_page && user_page != "" {
		rows, err := database.Query("SELECT role FROM users WHERE username = ?", user_page)
		CheckError(err)
		for rows.Next() {
			rows.Scan(&role)
		}
		defer rows.Close()
		data_info["roleUserPage"] = role
	}
}

// Affiche les notifications de l'utilisateur s'il y en a
func CheckNotif(w http.ResponseWriter, r *http.Request, data_notif map[string]interface{}) {
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)

	var notification string
	var arrNotification []string
	var test []string
	var arr_notif [][]string
	rows, err := database.Query("SELECT notification FROM users WHERE username = ?", data_notif["user"])
	CheckError(err)
	for rows.Next() {
		err := rows.Scan(&notification)
		CheckError(err)
		// get the section splited with the ","
		arrNotification = strings.Split(notification, ",")
		// remove the last part of the section which is empty
		if len(arrNotification[len(arrNotification)-1]) < 1 {
			arrNotification = arrNotification[0 : len(arrNotification)-1]
		}
		// Now split with white spaces and add them to the final array
		for i := 0; i < len(arrNotification); i++ {
			test = strings.Split(arrNotification[i], " ")
			arr_notif = append(arr_notif, test)
		}
	}
	err = rows.Err()
	CheckError(err)

	data_notif["notif"] = arr_notif

	// Delete notif
	del_notif := r.FormValue("del-notif")
	userToDel := r.FormValue("user")
	if del_notif == "1" {
		tx, err := database.Begin()
		CheckError(err)
		stmt, err := tx.Prepare("UPDATE users SET notification = ? WHERE username = ?")
		CheckError(err)
		_, err = stmt.Exec("", userToDel)
		CheckError(err)
		tx.Commit()
		rows.Close()
		database.Close()
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
	rows.Close()
	database.Close()
}

// get all posts liked/created/commented by someone
func Feed(data_Info map[string]interface{}) {
	//Show the liked post by the user
	LikedPosts(data_Info, "feedlike")

	//Show the post created by the user
	Createdposts(data_Info, "userpage")

	//Show comment posted
	// commentaires
	var comments [][11]string
	var content string
	var tmp [11]string
	database_comment, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database_comment.Close()
	//range over database
	rows_comment, err := database_comment.Query("SELECT content, idMainPost, date, id FROM comments WHERE author = ?", data_Info["username"])
	CheckError(err)
	defer rows_comment.Close()
	for rows_comment.Next() {

		err := rows_comment.Scan(&content, &tmp[1], &tmp[2], &tmp[5])
		CheckError(err)
		// Remplace les \n par des <br> pour sauter des lignes en html
		tmp[0] = strings.Replace(content, string('\r'), "", -1)
		tmp[0] = strings.Replace(content, string('\n'), "<br>", -1)

		// Ajoute le chemin de la photo qui a été choisit par l'utilisateur
		rows, err := database_comment.Query("SELECT photo FROM users WHERE username = ?", data_Info["username"])
		CheckError(err)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&tmp[4])
			CheckError(err)
		}
		err = rows.Err()
		CheckError(err)

		rows_posts, err := database_comment.Query("SELECT id, title, body, author, date FROM posts WHERE id = ?", tmp[1])
		CheckError(err)
		defer rows_posts.Close()
		for rows_posts.Next() {
			err := rows_posts.Scan(&tmp[6], &tmp[7], &tmp[8], &tmp[9], &tmp[10])
			CheckError(err)
		}
		comments = append(comments, tmp)
	}

	err = rows_comment.Err()
	CheckError(err)

	data_Info["commentsPosted"] = comments
}

// get all posts liked by Someone
func LikedPosts(data_Info map[string]interface{}, state string) {
	var categories = GetBruteCategories()
	var all_myLikedPosts [][]interface{}
	var post_liked bool

	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	//range over database
	rows, err := database.Query("SELECT title, body, author, date, id, category, likedBy FROM posts")
	CheckError(err)

	for rows.Next() {
		myLikedPosts := []interface{}{"", "", "", "", "", "", ""}
		rows.Scan(&myLikedPosts[0], &myLikedPosts[1], &myLikedPosts[2], &myLikedPosts[3], &id, &myLikedPosts[6], &likedBy)

		// find all the posts which the user liked
		like_splitted := strings.Split(likedBy, " ")
		for i := 0; i < len(like_splitted); i++ {
			if state == "indexlike" {
				if like_splitted[i] == data_Info["user"].(string) {
					post_liked = true
					break
				} else {
					post_liked = false
				}
			} else if state == "feedlike" {
				if like_splitted[i] == data_Info["username"] {
					post_liked = true
					break
				} else {
					post_liked = false
				}
			}

		}

		if post_liked {
			// Remplace les \n par des <br> pour sauter des lignes en html
			myLikedPosts[1] = strings.Replace(myLikedPosts[1].(string), string('\r'), "", -1)
			myLikedPosts[1] = strings.Replace(myLikedPosts[1].(string), string('\n'), "<br>", -1)
			myLikedPosts[5] = strconv.Itoa(id)
			if myLikedPosts[6] != nil {
				temp := []interface{}{} // string
				for _, e := range strings.Split(myLikedPosts[6].(string), ",") {
					j, _ := strconv.Atoi(string(e))
					if categories[j] != "" {
						temp = append(temp, categories[j])
					}
				}
				myLikedPosts = append(myLikedPosts, temp)
			} else {
				myLikedPosts[6] = []string{}
				myLikedPosts = append(myLikedPosts, []string{})
			}
			all_myLikedPosts = append(all_myLikedPosts, myLikedPosts)
		}
	}
	rows.Close()

	if post_liked {
		// Ajoute le chemin de la photo qui a été choisit par l'utilisateur
		for i := 0; i < len(all_myLikedPosts); i++ {
			rows, err := database.Query("SELECT photo FROM users WHERE username = ?", all_myLikedPosts[i][2])
			CheckError(err)
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&photo)
				CheckError(err)
			}
			err = rows.Err()
			CheckError(err)
			all_myLikedPosts[i][4] = photo
		}
	}

	data_Info["all_myLikedPosts"] = all_myLikedPosts
}

// get all posts created by someone
func Createdposts(data_Info map[string]interface{}, state string) {
	var all_myPosts [][]interface{}
	var photo string

	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()
	//range over database
	var rows *sql.Rows
	if state == "userpage" {
		rows, err = database.Query("SELECT title, body, author, date, id, category FROM posts WHERE author = ?", data_Info["username"])
		CheckError(err)
	} else if state == "indexuser" {
		rows, err = database.Query("SELECT title, body, author, date, id, category FROM posts WHERE author = ?", data_Info["user"].(string))
		CheckError(err)
	}
	categories := GetBruteCategories()
	for rows.Next() {
		myPosts := []interface{}{"", "", "", "", "", "", ""}
		rows.Scan(&myPosts[0], &myPosts[1], &myPosts[2], &myPosts[3], &id, &myPosts[6])
		// si le RegExp correspond à la DB
		// Remplace les \n par des <br> pour sauter des lignes en html
		myPosts[1] = strings.Replace(myPosts[1].(string), string('\r'), "", -1)
		myPosts[1] = strings.Replace(myPosts[1].(string), string('\n'), "<br>", -1)
		myPosts[5] = strconv.Itoa(id)
		if myPosts[6] != nil {
			temp := []interface{}{} // string
			for _, e := range strings.Split(myPosts[6].(string), ",") {
				j, err := strconv.Atoi(string(e))
				if err == nil && categories[j] != "" {
					temp = append(temp, categories[j])
				}
			}
			myPosts = append(myPosts, temp)
		} else {
			myPosts[6] = []string{}
			myPosts = append(myPosts, []string{})
		}
		all_myPosts = append(all_myPosts, myPosts)
	}

	rows.Close()
	// Ajoute le chemin de la photo qui a été choisit par l'utilisateur
	for i := 0; i < len(all_myPosts); i++ {
		rows, err := database.Query("SELECT photo FROM users WHERE username = ?", all_myPosts[i][2])
		CheckError(err)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&photo)
			CheckError(err)
		}
		err = rows.Err()
		CheckError(err)
		all_myPosts[i][4] = photo
	}

	data_Info["all_myPosts"] = all_myPosts
}

func SearchUserToLog(data_logIn map[string]interface{}, user_login string, password_login string, state string, data map[string]interface{}) bool {
	var password string
	var create_cookie = false

	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	// Parcourir la BDD
	rows, err := database.Query("SELECT username, password FROM users")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&username, &password)
		// Si l'input username est trouvé
		if user_login == username {
			// Compare l'input password avec celui de la BDD
			if state == "logFromExternalWebsite" || ComparePasswords(password, []byte(password_login)) {
				create_cookie = true
				data["user"] = username
				data["cookieExist"] = true
				break
			} else {
				data_logIn["wrongPassword"] = true
			}
		} else if user_login != "" {
			data_logIn["wrongUsername"] = true

		}
	}
	return create_cookie
}

func GetRoleForAllUsers(data_allUsers map[string]interface{}, role string) {
	var someone [2]string
	var everyone [][2]string

	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()

	// Ajouter tous les admins
	rows, err := database.Query("SELECT username, email FROM users WHERE role = ?", role)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&someone[0], &someone[1])
		everyone = append(everyone, someone)
	}
	data_allUsers["all"+role] = everyone
}
