package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"

	. "./config"
	. "./cookies"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"

	_ "github.com/mattn/go-sqlite3"
)

// Refuser une demande d'un user pour devenir modo créé un fichier db-journal (peut etre la page dashboard)
// Lors de la suppression d'un user, retirer sa demande de devenir admin + sur ses commentaires de posts, retirer le lien cliquable + son nom dans la liste de likedBy
// Utiliser l'ID utilisateur comme identifiant (998168070934676) sinon 2 comptes avec le meme nom peuvent se connecter

var id, like int
var username, email, age, fewWords, address, photo, state, title, body, author, date, content, likedBy, nbComments string
var create_cookie, userFound = false, false
var categories = []string{"gaming", "informatique", "sport", "culture", "politique", "loisir", "sciences", "sexualite", "finance"}
var data = make(map[string]interface{})

func main() {
	data["user"] = ""

	CreateDB()

	fmt.Println("Please connect to http://localhost:8000")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.HandleFunc("/", index)
	http.HandleFunc("/SignUp", SignUp)
	http.HandleFunc("/logIn", logIn)
	http.HandleFunc("/logOut", logOut)
	http.HandleFunc("/welcome", welcome)
	http.HandleFunc("/allUsers", allUsers)
	http.HandleFunc("/user", user)
	http.HandleFunc("/newPost", newPost)
	http.HandleFunc("/post", post)
	http.HandleFunc("/delPost", delPost)
	http.HandleFunc("/delComment", delComment)
	http.HandleFunc("/myPosts", myPosts)
	http.HandleFunc("/myLikedPosts", myLikedPosts)
	http.HandleFunc("/PendingPosts", PendingPosts)
	http.HandleFunc("/dashboard", Dashboard)
	http.HandleFunc("/loginFB", HandleFacebookLogin)
	http.HandleFunc("/oauth2callback", HandleFacebookCallback)

	err := http.ListenAndServe(":8000", nil) // Set listen port
	checkError(err)
}

// Generate the main page when first loading the site
func index(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_index := make(map[string]interface{})
	for k, v := range data {
		data_index[k] = v
	}
	data_index["cookieExist"] = false
	GetCookie(w, data_index, r) // ./cookies/getCookies.go

	if data["user"] != nil {
		CheckNotif(w, r, data_index)
		GetRole(data_index, false, "")
	} else {
		data_index["cookieExist"] = false
	}

	state = "index"
	DisplayPosts(r, data_index, state) // ./config/post.go

	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "index", data_index)
}

// Generate the log In page
func logIn(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_logIn := make(map[string]interface{})
	for k, v := range data {
		data_logIn[k] = v
	}
	data_logIn["cookieExist"] = false
	GetCookie(w, data_logIn, r)

	if data_logIn["cookieExist"] == true {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	// get user input to log in
	user_login := r.FormValue("user-login")
	password_login := r.FormValue("password-login")
	if FBuser.Name != "" {
		create_cookie = true
		data["user"] = FBuser.Name
		data["cookieExist"] = true
		AddUser(FBuser.Name, "", "", data_logIn)
		SearchUserToLog(data_logIn, user_login, "")
	} else {
		create_cookie = false
		SearchUserToLog(data_logIn, user_login, password_login)
	}

	if create_cookie {
		// Créé un cookie si user bien authentifié
		CreateCookie(w, r)
		data_logIn["wrongPassword"] = false
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	t := template.New("logIn-template")
	t = template.Must(t.ParseFiles("./html/LogIn.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "LogIn", data_logIn)
}

// Generate the sign in page
func SignUp(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_SignUp := make(map[string]interface{})
	for k, v := range data {
		data_SignUp[k] = v
	}

	// Check if a cookie exist
	data_SignUp["cookieExist"] = false
	GetCookie(w, data_SignUp, r)

	// if a cookie already exist, redirect
	if data_SignUp["cookieExist"] == true {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	// intiate values of the map
	data_SignUp["accountCreated"] = false
	data_SignUp["username_used"] = false
	data_SignUp["email_used"] = false

	// Get the inputs
	input_username := r.FormValue("username")
	input_email := r.FormValue("email")
	input_password := r.FormValue("password")

	// Crypt the password
	hash := HashPassword(input_password) // ./config/HashPassword.go

	// Call the function to add users
	if input_email != "" || input_username != "" || input_password != "" {
		AddUser(input_username, input_email, hash, data_SignUp) // ./config/AddUser.go
	}

	t := template.New("SignUp-template")
	t = template.Must(t.ParseFiles("./html/SignUp.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "SignUp", data_SignUp)
}

func SearchUserToLog(data_logIn map[string]interface{}, user_login string, password_login string) {
	var password string

	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	checkError(err)
	defer database.Close()

	// Parcourir la BDD
	rows, err := database.Query("SELECT username, password FROM users")
	checkError(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&username, &password)
		// Si l'input username est trouvé
		if user_login == username {
			// Compare l'input password avec celui de la BDD
			if ComparePasswords(password, []byte(password_login)) {
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
}

// Generate the Welcome page
func welcome(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_welcome := make(map[string]interface{})
	for k, v := range data {
		data_welcome[k] = v
	}

	GetCookie(w, data_welcome, r)
	if data["user"] != nil {
		CheckNotif(w, r, data_welcome)
		GetRole(data_welcome, false, "")
	}

	t := template.New("welcome-template")
	t = template.Must(t.ParseFiles("./html/welcome.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "welcome", data_welcome)
}

// Generate the user page
func user(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_user := make(map[string]interface{})
	for k, v := range data {
		data_user[k] = v
	}
	data_user["cookieExist"] = false
	data_user["username"] = ""
	GetCookie(w, data_user, r)

	if data["user"] != nil {
		CheckNotif(w, r, data_user)                 // ./config/user.go
		on_user_page := true                        // Lorsque la page actuelle est le profil d'un utilisateur
		user_page := r.FormValue("user")            //
		GetRole(data_user, on_user_page, user_page) // Enregistre le role de l'utilisateur
		// Change un role lorsqu'un Admin submit le formulaire
		ChangeRole(w, r) // ./config/modifDB.go
	}

	// Récupère les infos de l'utilisateur
	GetInfoUser(w, r, data_user) // ./config/modifDB.go
	Feed(data_user)
	// Check if the user logged is on his personnal page
	if data_user["username"] == data_user["user"] && data_user["cookieExist"] != false {
		data_user["sameUser"] = true
	} else {
		data_user["sameUser"] = false
	}

	// Change the photo / info of users
	UpdateInfoUsers(w, r, id) // ./config/modifDB.go

	delete_account := r.FormValue("del-account")
	if delete_account != "" {
		DelAccount(delete_account) // ./config/modifDB.GO
		if data_user["sameUser"] == true {
			DeleteCookie(w)
		}
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	// Add the user that ask for being moderator
	userAskingForModo := r.FormValue("user")
	if userAskingForModo != "" {
		ResquestForModo(userAskingForModo)
	}

	// Add the user when a moderator reports him
	nameUser := r.FormValue("nameUser")
	nameReporter := r.FormValue("nameReporter")
	reasonReport := r.FormValue("reasonReport")
	if nameUser != "" && nameReporter != "" && reasonReport != "" {
		Report(nameUser, nameReporter, reasonReport)
	}

	t := template.New("user-template")
	t = template.Must(t.ParseFiles("./html/user.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "user", data_user)
}

// Generate allUsers page
func allUsers(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	var aUser [2]string
	var all_users [][2]string
	var all_moderator [][2]string
	var all_admin [][2]string
	data_allUsers := make(map[string]interface{})
	for k, v := range data {
		data_allUsers[k] = v
	}
	GetCookie(w, data_allUsers, r)
	if data["user"] != nil {
		CheckNotif(w, r, data_allUsers)
		GetRole(data_allUsers, false, "")
	}

	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	checkError(err)
	defer database.Close()

	var role string
	// Ajouter tous les admins
	rows, err := database.Query("SELECT username, email, photo, role FROM users WHERE role = ?", "admin")
	checkError(err)
	for rows.Next() {
		rows.Scan(&username, &email, &photo, &role)
		aUser[0] = username
		aUser[1] = email
		all_admin = append(all_admin, aUser)
	}
	rows.Close()
	data_allUsers["allAdmin"] = all_admin
	// Ajouter tous les moderateur
	rows, err = database.Query("SELECT username, email, photo, role FROM users WHERE role = ?", "moderator")
	checkError(err)
	for rows.Next() {
		rows.Scan(&username, &email, &photo, &role)
		aUser[0] = username
		aUser[1] = email
		all_moderator = append(all_moderator, aUser)
	}
	rows.Close()
	data_allUsers["allModerator"] = all_moderator
	// Ajouter tous les users
	rows, err = database.Query("SELECT username, email, photo, role FROM users WHERE role = ?", "user")
	checkError(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&username, &email, &photo, &role)
		aUser[0] = username
		aUser[1] = email
		all_users = append(all_users, aUser)
	}
	data_allUsers["allUsers"] = all_users

	t := template.New("allUsers-template")
	t = template.Must(t.ParseFiles("./html/allUsers.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "allUsers", data_allUsers)
}

// Delete the cookie to disconnect the user
func logOut(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	DeleteCookie(w)

	// reset the login from FB
	FBuser.Name = ""
	http.Redirect(w, r, "/index", http.StatusSeeOther)
	delete(data, "user")
}

// Generate the page to create new Post (accessible only to logged in users)
func newPost(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_newPost := make(map[string]interface{})
	for k, v := range data {
		data_newPost[k] = v
	}
	GetCookie(w, data_newPost, r)
	// Redirection pour ceux qui ne sont pas connecté
	if data_newPost["cookieExist"] == false {
		http.Redirect(w, r, "/logIn", http.StatusSeeOther)
	}
	if data["user"] != nil {
		CheckNotif(w, r, data_newPost)
		GetRole(data_newPost, false, "")
	}

	// Input de la page
	title := r.FormValue("title")
	body := r.FormValue("body")

	category := []string{r.FormValue("gaming"), r.FormValue("informatique"), r.FormValue("sport"), r.FormValue("culture"), r.FormValue("politique"), r.FormValue("loisir"), r.FormValue("sciences"), r.FormValue("sexualite"), r.FormValue("finance")}

	if title != "" && body != "" {
		// Capture la date de submit
		dt := time.Now()
		// appel de la fonction pour créer le post
		AddNewPost(title, body, dt.Format("02-01-2006 15:04:05"), data_newPost, category)
		UploadHandler(w, r)
	}
	data_newPost["categorie"] = categories

	t := template.New("newPost-template")
	t = template.Must(t.ParseFiles("./html/newPost.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "newPost", data_newPost)
}

// Generate page des posts avec ses commentaires
func post(w http.ResponseWriter, r *http.Request) {
	// var dislike int
	var dislikedBy string
	post_id := r.FormValue("id")

	// initiate the data that will be send to html
	data_post := make(map[string]interface{})
	for k, v := range data {
		data_post[k] = v
	}
	data_post["cookieExist"] = false
	data_post["already_liked"] = false
	data_post["already_disliked"] = false
	GetCookie(w, data_post, r)
	if data["user"] != nil {
		CheckNotif(w, r, data_post)
		GetRole(data_post, false, "")
	}

	// Affiche les posts et commentaires
	var post = Display_post_comment(post_id, data_post, body)
	Display_comments(data_post, post_id)

	// Ajoute un commentaire
	add_comment := r.FormValue("add_comment")
	if add_comment != "" {
		Adding_comment(w, r, add_comment, &post, data_post["user"].(string))
		http.Redirect(w, r, "/post?id="+post_id, http.StatusSeeOther)
	}

	// Système de Like
	change_nmb_like := r.FormValue("Like")
	likedBy = Like(change_nmb_like, data_post, post_id, w, r)

	if data_post["user"] == nil {
		data_post["user"] = ""
	}
	// Système de disLike
	change_nmb_dislike := r.FormValue("DisLike")
	dislikedBy = DisLike(change_nmb_dislike, data_post, post_id, w, r)

	if data_post["user"] == nil {
		data_post["user"] = ""
	}

	// Modification de post
	modif_post := r.FormValue("modifPost")
	id_mainPost := r.FormValue("id-mainPost")
	modif_comment := r.FormValue("modifComment")
	id_comment := r.FormValue("id-comment")
	if modif_post != "" && id_mainPost != "" || modif_comment != "" && id_comment != "" {
		ModifPostAndComment(modif_post, id_mainPost, modif_comment, id_comment)
		http.Redirect(w, r, "/post?id="+post_id, http.StatusSeeOther)
	}

	// ajoute les personnes qui ont liké le post principal
	likedBy = strings.ReplaceAll(likedBy, " ", "<br>")
	data_post["mainPost_likedBy"] = likedBy
	// ajoute les personnes qui ont disliké le post principal
	dislikedBy = strings.ReplaceAll(dislikedBy, " ", "<br>")
	data_post["mainPost_dislikedBy"] = dislikedBy

	t := template.New("post-template")
	t = template.Must(t.ParseFiles("./html/post.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "post", data_post)
}

func delPost(w http.ResponseWriter, r *http.Request) {
	delete_post := r.FormValue("delPost")
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	checkError(err)
	defer database.Close()

	// delete the image
	state := "mainPost"
	Delete_image(state, delete_post)

	// DELETE the comments of the main post
	tx, err := database.Begin()
	checkError(err)
	stmt, err := tx.Prepare("DELETE FROM comments WHERE idMainPost = ?")
	checkError(err)
	_, err = stmt.Exec(delete_post)
	checkError(err)

	// DELETE the main POST
	stmt, err = tx.Prepare("DELETE FROM posts WHERE id = ?")
	checkError(err)
	_, err = stmt.Exec(delete_post)
	checkError(err)
	tx.Commit()

	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

// supprime un commentaire
func delComment(w http.ResponseWriter, r *http.Request) {
	delete_comment := r.FormValue("delComment")
	idMainPost := r.FormValue("id-mainPost")
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	checkError(err)
	defer database.Close()

	// DELETE the comments of the main post
	tx, err := database.Begin()
	checkError(err)
	stmt, err := tx.Prepare("DELETE FROM comments WHERE id = ?")
	checkError(err)
	_, err = stmt.Exec(delete_comment)
	checkError(err)
	tx.Commit()
	http.Redirect(w, r, "/post?id="+idMainPost, http.StatusSeeOther)
}

// Page qui affiche les posts créé par l'utilisateur connecté
func myPosts(w http.ResponseWriter, r *http.Request) {

	// initiate the data that will be send to html
	data_myPosts := make(map[string]interface{})
	for k, v := range data {
		data_myPosts[k] = v
	}
	data_myPosts["cookieExist"] = false
	GetCookie(w, data_myPosts, r)
	if data["user"] != nil {
		CheckNotif(w, r, data_myPosts)
		GetRole(data_myPosts, false, "")
	}
	Createdposts(data_myPosts, "indexuser")

	t := template.New("myPosts-template")
	t = template.Must(t.ParseFiles("./html/myPosts.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "myPosts", data_myPosts)
}

// Page qui affiche les posts liké par l'utilisateur connecté
func myLikedPosts(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_myLikedPosts := make(map[string]interface{})
	for k, v := range data {
		data_myLikedPosts[k] = v
	}
	data_myLikedPosts["cookieExist"] = false
	GetCookie(w, data_myLikedPosts, r)
	if data["user"] != nil {
		CheckNotif(w, r, data_myLikedPosts)
		GetRole(data_myLikedPosts, false, "")
	}

	LikedPosts(data_myLikedPosts, "indexlike")

	t := template.New("myLikedPosts-template")
	t = template.Must(t.ParseFiles("./html/myLikedPosts.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "myLikedPosts", data_myLikedPosts)
}

// verifie les erreurs
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// page des posts en attente, un admin ou moderateur doit les accepter pour les mettres sur la page principal
func PendingPosts(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_pendingPosts := make(map[string]interface{})
	for k, v := range data {
		data_pendingPosts[k] = v
	}

	GetCookie(w, data_pendingPosts, r)
	// Redirection pour ceux qui ne sont pas connecté
	if data["user"] != nil {
		CheckNotif(w, r, data_pendingPosts)
		GetRole(data_pendingPosts, false, "")
	}
	if data_pendingPosts["cookieExist"] == false || data_pendingPosts["role"] == "user" {
		http.Redirect(w, r, "/logIn", http.StatusSeeOther)
	}

	post_accepted := r.FormValue("post-accepted")
	id_pendingPost := r.FormValue("id-pendingPost")
	// supprime ou déplace le post s'il est accepté ou non
	if post_accepted != "" && id_pendingPost != "" {
		PostAcceptedOrNot(post_accepted, id_pendingPost) // ./config/modifDB.go
	}

	// affiche le post
	state = "pendingPosts"
	DisplayPosts(r, data_pendingPosts, state) // ./config/post.go

	t := template.New("pendingPosts-template")
	t = template.Must(t.ParseFiles("./html/pendingPosts.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "pendingPosts", data_pendingPosts)
}

// get all posts liked by Someone
func LikedPosts(data_Info map[string]interface{}, state string) {

	var all_myLikedPosts [][]interface{}
	var post_liked bool

	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	checkError(err)
	defer database.Close()

	//range over database
	rows, err := database.Query("SELECT title, body, author, date, id, category, likedBy FROM posts")
	checkError(err)

	for rows.Next() {
		myLikedPosts := []interface{}{"", "", "", "", "", "", ""}
		rows.Scan(&myLikedPosts[0], &myLikedPosts[1], &myLikedPosts[2], &myLikedPosts[3], &id, &myLikedPosts[6], &likedBy)

		// find all the posts which the user liked
		like_splitted := strings.Split(likedBy, " ")
		for i := 0; i < len(like_splitted); i++ {
			if state == "indexlike" {
				if like_splitted[i] == data["user"].(string) {
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
				for _, e := range myLikedPosts[6].(string) {
					j, _ := strconv.Atoi(string(e))
					temp = append(temp, categories[j])
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
			checkError(err)
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&photo)
				checkError(err)
			}
			err = rows.Err()
			checkError(err)
			all_myLikedPosts[i][4] = photo
		}
	}

	data_Info["all_myLikedPosts"] = all_myLikedPosts
}

// get all posts created by someone
func Createdposts(data_Info map[string]interface{}, state string) {
	var all_myPosts [][]interface{}

	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	checkError(err)
	defer database.Close()
	//range over database
	var rows *sql.Rows
	if state == "userpage" {
		rows, err = database.Query("SELECT title, body, author, date, id, category FROM posts WHERE author = ?", data_Info["username"])
		checkError(err)
	} else if state == "indexuser" {
		rows, err = database.Query("SELECT title, body, author, date, id, category FROM posts WHERE author = ?", data["user"].(string))
		checkError(err)
	}

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
			for _, e := range myPosts[6].(string) {
				j, _ := strconv.Atoi(string(e))
				temp = append(temp, categories[j])
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
		checkError(err)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&photo)
			checkError(err)
		}
		err = rows.Err()
		checkError(err)
		all_myPosts[i][4] = photo
	}

	data_Info["all_myPosts"] = all_myPosts
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
	checkError(err)
	defer database_comment.Close()
	//range over database
	rows_comment, err := database_comment.Query("SELECT content, idMainPost, date, id FROM comments WHERE author = ?", data_Info["username"])
	checkError(err)
	defer rows_comment.Close()
	for rows_comment.Next() {

		err := rows_comment.Scan(&content, &tmp[1], &tmp[2], &tmp[5])
		checkError(err)
		// Remplace les \n par des <br> pour sauter des lignes en html
		tmp[0] = strings.Replace(content, string('\r'), "", -1)
		tmp[0] = strings.Replace(content, string('\n'), "<br>", -1)

		// Ajoute le chemin de la photo qui a été choisit par l'utilisateur
		rows, err := database_comment.Query("SELECT photo FROM users WHERE username = ?", data_Info["username"])
		checkError(err)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&tmp[4])
			checkError(err)
		}
		err = rows.Err()
		checkError(err)

		rows_posts, err := database_comment.Query("SELECT id, title, body, author, date FROM posts WHERE id = ?", tmp[1])
		checkError(err)
		defer rows_posts.Close()
		for rows_posts.Next() {
			err := rows_posts.Scan(&tmp[6], &tmp[7], &tmp[8], &tmp[9], &tmp[10])
			checkError(err)
		}
		comments = append(comments, tmp)
	}

	err = rows_comment.Err()
	checkError(err)

	data_Info["commentsPosted"] = comments
}

// page to see whose asking to be moderator and the reports
func Dashboard(w http.ResponseWriter, r *http.Request) {
	data_dashboard := make(map[string]interface{})
	for k, v := range data {
		data_dashboard[k] = v
	}
	GetCookie(w, data_dashboard, r)
	if data_dashboard["user"] != nil {
		CheckNotif(w, r, data_dashboard)
		GetRole(data_dashboard, false, "")
	} else {
		data_dashboard["cookieExist"] = false
	}
	RefuserDemande(w, r)
	AccepterDemande(w, r)
	DisplayAdminModo(&data_dashboard)
	DisplayPendingForModo(&data_dashboard)
	SelectReport(data_dashboard)

	// get input of the admin that accept or not the report
	answerReport := r.FormValue("answerReport")
	nameReported := r.FormValue("nameReported")
	reportAccepted := r.FormValue("reportAccepted") // /!\ checkbox

	if answerReport != "" && nameReported != "" {
		if reportAccepted == "" {
			reportAccepted = "0"
		}
		DeleteUserFromReport(answerReport, nameReported, reportAccepted)
	}

	t := template.New("dashboard-template")
	t = template.Must(t.ParseFiles("./html/dashboard.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "dashboard", data_dashboard)
}

// add the report to the database
func Report(nameUser string, nameReporter string, reasonReport string) {
	// Open the database
	var nameReported string
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	checkError(err)
	defer database.Close()
	//verify if the user is already in the table
	rows_double, err := database.Query("SELECT nameUser FROM report")
	checkError(err)
	defer rows_double.Close()
	for rows_double.Next() {
		err := rows_double.Scan(&nameReported)
		if err != nil {
			log.Fatal(err)
		}
		if nameReported == nameUser {
			return
		}
	}

	// add the user to the table report
	tx, err := database.Begin()
	checkError(err)
	stmt, err := tx.Prepare("INSERT INTO report (nameUser, reasonReport, answerReport, reported, nameReporter) VALUES (?,?,'','',?)")
	checkError(err)
	_, err = stmt.Exec(nameUser, reasonReport, nameReporter)
	checkError(err)

	tx.Commit()
}

// get all the report pending
func SelectReport(data_dashboard map[string]interface{}) {
	var slicereport [3]string
	var all_Report [][3]string
	// Open the database
	database_report, err := sql.Open("sqlite3", "./db-sqlite.db")
	checkError(err)
	defer database_report.Close()
	rows_report, err := database_report.Query("SELECT nameReporter, nameUser, reasonReport FROM report WHERE reported = ?", "")
	checkError(err)
	defer rows_report.Close()
	for rows_report.Next() {
		err := rows_report.Scan(&slicereport[0], &slicereport[1], &slicereport[2])
		if err != nil {
			log.Fatal(err)
		}
		all_Report = append(all_Report, slicereport)

	}
	data_dashboard["report"] = all_Report
}

func DeleteUserFromReport(answerReport string, nameReported string, reportAccepted string) {
	// Open the database
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	checkError(err)
	defer database.Close()

	// if the report is accepted from the admnin
	if reportAccepted == "1" {
		// delete everything from this user
		DelAccount(nameReported)
		// Update the db report
		tx, err := database.Begin()
		checkError(err)
		stmt, err := tx.Prepare("UPDATE report SET answerReport = ?, reported = ? WHERE nameUser = ?")
		checkError(err)
		_, err = stmt.Exec(answerReport, "reportAccepted", nameReported)
		checkError(err)
		tx.Commit()
	}
}

var (
	oauthConf = &oauth2.Config{
		ClientID:     "166859168782629",
		ClientSecret: "d7fc1ba6430a4c122b41caf0282fea44",
		RedirectURL:  "http://localhost:8000/oauth2callback",
		Scopes:       []string{"public_profile"},
		Endpoint:     facebook.Endpoint,
	}
	oauthStateString = "eeeeeeee"
)

type Info struct {
	Name string
	Id   int
}

var info Info
var FBuser = &info

func HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	Url, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Println("invalid oauth state, expected , got \n", oauthStateString, state)
		http.Redirect(w, r, "/logIn", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("oauthConf.Exchange() failed with \n", err)
		http.Redirect(w, r, "/logIn", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		fmt.Println("Get: \n", err)
		http.Redirect(w, r, "/logIn", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll: \n", err)
		http.Redirect(w, r, "/logIn", http.StatusTemporaryRedirect)
		return
	}
	json.Unmarshal(response, &info)
	info.Name = strings.ReplaceAll(info.Name, " ", "_")
	http.Redirect(w, r, "/logIn", http.StatusTemporaryRedirect)
}
