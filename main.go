package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	. "./config"
	. "./cookies"
	. "./oauth"

	_ "github.com/mattn/go-sqlite3"
)

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
	http.HandleFunc("/delPost", DelPost)       // ./config/modifDB
	http.HandleFunc("/delComment", DelComment) // ./config/modifDB
	http.HandleFunc("/myPosts", myPosts)
	http.HandleFunc("/myLikedPosts", myLikedPosts)
	http.HandleFunc("/pendingPosts", PendingPosts)
	http.HandleFunc("/dashboard", Dashboard)

	// Login with facebook
	http.HandleFunc("/loginFB", HandleFacebookLogin) // ./oauth/facebook
	http.HandleFunc("/FBloginCallBack", HandleFacebookCallback)

	// Login with google
	http.HandleFunc("/loginGoogle", HandleGoogleLogin) // ./oauth/google
	http.HandleFunc("/GoogleCallBack", HandleGoogleCallback)

	err := http.ListenAndServe(":8000", nil)
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

	DisplayPosts(r, data_index, "index") // ./config/post.go

	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "index", data_index)
}

// Generate the log In page
func logIn(w http.ResponseWriter, r *http.Request) {
	create_cookie := false
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
	if FB_google_user.Name != "" {
		AddUser(FB_google_user.Name, FB_google_user.Email, "", data_logIn)
		create_cookie = SearchUserToLog(data_logIn, FB_google_user.Name, "", "logFromExternalWebsite", data)
		data["user"] = FB_google_user.Name
		data["cookieExist"] = true
	} else {
		create_cookie = SearchUserToLog(data_logIn, user_login, password_login, "logFromInternalDB", data)
	}

	if create_cookie { // Créé un cookie si user bien authentifié
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
	UpdateInfoUsers(w, r) // ./config/modifDB.go

	delete_account := r.FormValue("del-account")
	if delete_account != "" {
		DelAccount(delete_account) // ./config/modifDB.GO
		if data_user["sameUser"] == true {
			DeleteCookie(w)
		}
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	// Add the user that ask for being moderator
	userAskingForModo := r.FormValue("askingToBeModo")
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

	// Ajouter tous les admins
	rows, err := database.Query("SELECT username, email FROM users WHERE role = ?", "admin")
	checkError(err)
	for rows.Next() {
		rows.Scan(&aUser[0], &aUser[1])
		all_admin = append(all_admin, aUser)
	}
	rows.Close()
	data_allUsers["allAdmin"] = all_admin
	// Ajouter tous les moderateur
	rows, err = database.Query("SELECT username, email FROM users WHERE role = ?", "moderator")
	checkError(err)
	for rows.Next() {
		rows.Scan(&aUser[0], &aUser[1])
		all_moderator = append(all_moderator, aUser)
	}
	rows.Close()
	data_allUsers["allModerator"] = all_moderator
	// Ajouter tous les users
	rows, err = database.Query("SELECT username, email FROM users WHERE role = ?", "user")
	checkError(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&aUser[0], &aUser[1])
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
	FB_google_user.Name = ""
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
	var likedBy, dislikedBy string
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
	var post = Display_post_comment(post_id, data_post)
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
	DisplayPosts(r, data_pendingPosts, "pendingPosts") // ./config/post.go

	t := template.New("pendingPosts-template")
	t = template.Must(t.ParseFiles("./html/pendingPosts.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "pendingPosts", data_pendingPosts)
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
	SelectReport(data_dashboard) // ./

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
