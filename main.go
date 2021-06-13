package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	config "./config"
	cookie "./cookies"
	oauth "./oauth"

	_ "github.com/mattn/go-sqlite3"
)

var data = make(map[string]interface{})

func main() {
	data["user"] = ""

	config.CreateDB()                // ./config/CreateDB
	config.InitCategoriePrincipale() // ./config/categorieManager

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
	http.HandleFunc("/delPost", config.DelPost)       // ./config/modifDB.go
	http.HandleFunc("/delComment", config.DelComment) // ./config/modifDB.go
	http.HandleFunc("/myPosts", myPosts)
	http.HandleFunc("/myLikedPosts", myLikedPosts)
	http.HandleFunc("/pendingPosts", PendingPosts)
	http.HandleFunc("/dashboard", Dashboard)
	http.HandleFunc("/categories", CategoriesHandleur)

	// Login with facebook
	http.HandleFunc("/loginFB", oauth.HandleFacebookLogin) // ./oauth/facebook.go
	http.HandleFunc("/FBloginCallBack", oauth.HandleFacebookCallback)

	// Login with google
	http.HandleFunc("/loginGoogle", oauth.HandleGoogleLogin) // ./oauth/google.go
	http.HandleFunc("/GoogleCallBack", oauth.HandleGoogleCallback)

	err := http.ListenAndServe(":8000", nil)
	config.CheckError(err)
}

// Generate the main page when first loading the site
func index(w http.ResponseWriter, r *http.Request) {
	// Initialize the data that will be send to html
	data_index := make(map[string]interface{})
	SetDataToSend(w, r, data_index, data, false, "")

	// Print posts to the index page
	config.DisplayPosts(r, data_index, "index") // ./config/post.go

	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "index", data_index)
}

// Generate the log In page
func logIn(w http.ResponseWriter, r *http.Request) {
	create_cookie := false

	// Initialize the data that will be send to html
	data_logIn := make(map[string]interface{})
	SetDataToSend(w, r, data_logIn, data, false, "")

	// Redirect if the user is already logged in
	if data_logIn["cookieExist"] == true {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	// Get user input to log in
	user_login := r.FormValue("user-login")
	password_login := r.FormValue("password-login")
	if oauth.FB_google_user.Name != "" { // when the user wants to log in from FB or Google
		config.AddUser(oauth.FB_google_user.Name, oauth.FB_google_user.Email, "", data_logIn)                             // ./config/user.go
		create_cookie = config.SearchUserToLog(data_logIn, oauth.FB_google_user.Name, "", "logFromExternalWebsite", data) // ./config/user.go
		data["user"] = oauth.FB_google_user.Name
		data["cookieExist"] = true
	} else { // when he logs from a created account
		create_cookie = config.SearchUserToLog(data_logIn, user_login, password_login, "logFromInternalDB", data) // ./config/user.go
	}

	// Créé un cookie si user bien authentifié
	if create_cookie {
		cookie.CreateCookie(w, r) // ./cookies/createCookie.go
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
	SetDataToSend(w, r, data_SignUp, data, false, "")

	// Redirect if the user is already logged in
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
	hash := config.HashPassword(input_password) // ./config/HashPassword.go

	// Call the function to add users
	if input_email != "" || input_username != "" || input_password != "" {
		config.AddUser(input_username, input_email, hash, data_SignUp) // ./config/AddUser.go
	}

	t := template.New("SignUp-template")
	t = template.Must(t.ParseFiles("./html/signUp.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "SignUp", data_SignUp)
}

// Generate the Welcome page
func welcome(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_welcome := make(map[string]interface{})
	SetDataToSend(w, r, data_welcome, data, false, "")

	t := template.New("welcome-template")
	t = template.Must(t.ParseFiles("./html/welcome.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "welcome", data_welcome)
}

// Generate the user page
func user(w http.ResponseWriter, r *http.Request) {
	user_page := r.FormValue("user")
	// initiate the data that will be send to html
	data_user := make(map[string]interface{})
	SetDataToSend(w, r, data_user, data, true, user_page)

	if data["user"] != nil {
		// Change un role lorsqu'un Admin submit le formulaire
		config.ChangeRole(w, r) // ./config/modifDB.go
	}

	// Récupère les infos de l'utilisateur
	config.GetInfoUser(w, r, data_user) // ./config/modifDB.go
	config.Feed(data_user)              // ./config/user.go

	// Check if the user logged is on his personnal page
	if data_user["username"] == data_user["user"] && data_user["cookieExist"] != false {
		data_user["sameUser"] = true
	} else {
		data_user["sameUser"] = false
	}

	// Change the photo / info of users
	config.UpdateInfoUsers(w, r) // ./config/modifDB.go

	// Supprime le compte d'une personne
	delete_account := r.FormValue("del-account")
	if delete_account != "" {
		config.DelAccount(delete_account)  // ./config/modifDB.go
		if data_user["sameUser"] == true { // si l'utilisateur supprime son propre compte, détruit le cookie et le déconnecte
			cookie.DeleteCookie(w) // ./cookies/deleteCookie.go
		}
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	// Add the user that ask for being moderator
	userAskingForModo := r.FormValue("askingToBeModo")
	if userAskingForModo != "" {
		config.ResquestForModo(userAskingForModo)
	}

	// Add the user when a moderator reports him
	nameUser := r.FormValue("nameUser")
	nameReporter := r.FormValue("nameReporter")
	reasonReport := r.FormValue("reasonReport")
	if nameUser != "" && nameReporter != "" && reasonReport != "" {
		config.Report(nameUser, nameReporter, reasonReport)
	}

	t := template.New("user-template")
	t = template.Must(t.ParseFiles("./html/user.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "user", data_user)
}

// Generate allUsers page
func allUsers(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_allUsers := make(map[string]interface{})
	SetDataToSend(w, r, data_allUsers, data, false, "")

	// Enregistrer tous les admins, moderateurs et users pour les trier sur la page
	config.GetRoleForAllUsers(data_allUsers, "admin")
	config.GetRoleForAllUsers(data_allUsers, "moderator") // ./config/user
	config.GetRoleForAllUsers(data_allUsers, "user")

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

	cookie.DeleteCookie(w)

	// reset the login from FB
	oauth.FB_google_user.Name = ""
	http.Redirect(w, r, "/index", http.StatusSeeOther)
	delete(data, "user")
}

// Generate the page to create new Post (accessible only to logged in users)
func newPost(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_newPost := make(map[string]interface{})
	SetDataToSend(w, r, data_newPost, data, false, "")

	// Redirect users who are not connected
	if data_newPost["cookieExist"] == false {
		http.Redirect(w, r, "/logIn", http.StatusSeeOther)
	}

	// Input de la page
	title := r.FormValue("title")
	body := r.FormValue("body")

	categories := config.GetCategories()
	var category []string
	for i := 0; i < len(categories); i++ {
		query := "categorie" + strconv.Itoa(i)
		test := r.FormValue(query)
		if test != "" {
			category = append(category, test)
		}
	}

	if title != "" && body != "" {
		dt := time.Now()                                                                         // Capture la date de submit
		config.AddNewPost(title, body, dt.Format("02-01-2006 15:04:05"), data_newPost, category) // ./config/post.go
		config.UploadHandler(w, r)                                                               // ./config/images.go
	}
	data_newPost["categorie"] = config.GetCategories()

	t := template.New("newPost-template")
	t = template.Must(t.ParseFiles("./html/newPost.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "newPost", data_newPost)
}

// Generate page des posts avec ses commentaires
func post(w http.ResponseWriter, r *http.Request) {
	var likedBy, dislikedBy string
	post_id := r.FormValue("id")

	// initiate the data that will be send to html
	data_post := make(map[string]interface{})
	SetDataToSend(w, r, data_post, data, false, "")

	data_post["already_liked"] = false
	data_post["already_disliked"] = false

	// Affiche les posts et commentaires
	var post = config.Display_post_comment(post_id, data_post) // ./config/post.go
	config.Display_comments(data_post, post_id)                // ./config/post.go

	// Ajoute un commentaire
	add_comment := r.FormValue("add_comment")
	if add_comment != "" {
		config.Adding_comment(w, r, add_comment, &post, data_post["user"].(string)) // ./config/comment.go
		http.Redirect(w, r, "/post?id="+post_id, http.StatusSeeOther)
	}

	// Système de Like
	change_nmb_like := r.FormValue("Like")
	likedBy = config.Like(change_nmb_like, data_post, post_id, w, r) // ./config/like.go
	// Système de disLike
	change_nmb_dislike := r.FormValue("DisLike")
	dislikedBy = config.DisLike(change_nmb_dislike, data_post, post_id, w, r) // ./config/like.go

	// Modification de post
	modif_post := r.FormValue("modifPost")
	id_mainPost := r.FormValue("id-mainPost")
	modif_comment := r.FormValue("modifComment")
	id_comment := r.FormValue("id-comment")
	if modif_post != "" && id_mainPost != "" || modif_comment != "" && id_comment != "" {
		config.ModifPostAndComment(modif_post, id_mainPost, modif_comment, id_comment) // ./config/modifDB
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
	SetDataToSend(w, r, data_myPosts, data, false, "")
	config.Createdposts(data_myPosts, "indexuser") // ./config/user

	t := template.New("myPosts-template")
	t = template.Must(t.ParseFiles("./html/myPosts.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "myPosts", data_myPosts)
}

// Page qui affiche les posts liké par l'utilisateur connecté
func myLikedPosts(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_myLikedPosts := make(map[string]interface{})
	SetDataToSend(w, r, data_myLikedPosts, data, false, "")

	config.LikedPosts(data_myLikedPosts, "indexlike")

	t := template.New("myLikedPosts-template")
	t = template.Must(t.ParseFiles("./html/myLikedPosts.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "myLikedPosts", data_myLikedPosts)
}

// page des posts en attente, un admin ou moderateur doit les accepter pour les mettres sur la page principal
func PendingPosts(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_pendingPosts := make(map[string]interface{})
	SetDataToSend(w, r, data_pendingPosts, data, false, "")

	// Redirect users that are not connected or not admin or not moderator
	if data_pendingPosts["cookieExist"] == false || data_pendingPosts["role"] == "user" {
		http.Redirect(w, r, "/logIn", http.StatusSeeOther)
	}

	// Affiche les posts
	config.DisplayPosts(r, data_pendingPosts, "pendingPosts") // ./config/post.go

	// Supprime ou déplace le post s'il est accepté ou non
	post_accepted := r.FormValue("post-accepted")
	id_pendingPost := r.FormValue("id-pendingPost")
	if post_accepted != "" && id_pendingPost != "" {
		config.PostAcceptedOrNot(post_accepted, id_pendingPost) // ./config/modifDB.go
		http.Redirect(w, r, "/pendingPosts", http.StatusSeeOther)
	}

	t := template.New("pendingPosts-template")
	t = template.Must(t.ParseFiles("./html/pendingPosts.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "pendingPosts", data_pendingPosts)
}

// page to see whose asking to be moderator and the reports
func Dashboard(w http.ResponseWriter, r *http.Request) {
	data_dashboard := make(map[string]interface{})
	SetDataToSend(w, r, data_dashboard, data, false, "")

	// Affiche les roles des users, les requetes
	config.DisplayAdminModo(&data_dashboard)      // ./config/modo.go
	config.DisplayPendingForModo(&data_dashboard) // ./config/modo.go
	config.SelectReport(data_dashboard)           // ./ config/modo.go

	// Accepter ou non la demande d'un user pour devenir modo
	config.RefuserDemande(w, r)  // ./config/modo.go
	config.AccepterDemande(w, r) // ./config/modo.go

	// get inputs of the admin that accept or not the report
	answerReport := r.FormValue("answerReport")
	nameReported := r.FormValue("nameReported")
	reportAccepted := r.FormValue("reportAccepted") // /!\ checkbox
	if answerReport != "" && nameReported != "" {
		if reportAccepted == "" {
			reportAccepted = "0"
		}
		config.DeleteUserFromReport(answerReport, nameReported, reportAccepted) // ./config/report.go
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	t := template.New("dashboard-template")
	t = template.Must(t.ParseFiles("./html/dashboard.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "dashboard", data_dashboard)
}

// Copie les données principales, verifie si un cookie existe, les notifications ...
func SetDataToSend(w http.ResponseWriter, r *http.Request, data_info map[string]interface{}, data map[string]interface{}, on_user_page bool, user_page string) {
	// Copy the main map to get all important info
	for k, v := range data {
		data_info[k] = v
	}
	data_info["cookieExist"] = false
	data_info["username"] = ""
	cookie.GetCookie(w, data_info, r) // ./cookies/getCookies.go

	if data["user"] != nil {
		config.CheckNotif(w, r, data_info)
		config.GetRole(data_info, on_user_page, user_page)
	} else {
		data_info["cookieExist"] = false
	}
}

func CategoriesHandleur(w http.ResponseWriter, r *http.Request) {
	dataCategorie := make(map[string]interface{})
	SetDataToSend(w, r, dataCategorie, data, false, "")

	// Add new Categorie
	nom_categorie := r.FormValue("newCategorie")
	color_cat := r.FormValue("color-categorie")
	if nom_categorie != "" && color_cat != "" {
		config.NewCategorie(nom_categorie, color_cat)
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}

	// delete categorie
	del_categorie := r.FormValue("del-cat")
	if del_categorie != "" {
		config.Del_Categorie(del_categorie) // ./config/categorieManager.go
	}

	dataCategorie["categories"] = config.DisplayCategories()

	t := template.New("categoriesManager-template")
	t = template.Must(t.ParseFiles("./html/categoriesManager.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "categoriesManager", dataCategorie)
}
