package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	. "./config"
	. "./cookies"

	_ "github.com/mattn/go-sqlite3"
)

var id, like int
var username, password, email, age, fewWords, address, photo, state, title, body, author, date, content, likedBy, nbComments string
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

	err := http.ListenAndServe(":8000", nil) // Set listen port
	checkError(err)
}

// Generate the main page when first loading the site
func index(w http.ResponseWriter, r *http.Request) {
	var post [][]interface{}

	// initiate the data that will be send to html
	data_index := make(map[string]interface{})
	for k, v := range data {
		data_index[k] = v
	}
	data_index["cookieExist"] = false
	GetCookie(data_index, r)

	// filtre de categorie
	selected_categories := ""
	for i := range categories {
		tmp := r.FormValue(categories[i])
		if tmp != "" {
			selected_categories += strconv.Itoa(i)
		}
	}
	// RegExp
	var filter = regexp.MustCompile(selected_categories)

	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	//range over database
	rows, _ := database.Query("SELECT title, body, author, date, id, category, like, nbComments FROM posts")
	defer rows.Close()

	for rows.Next() {
		aPost := []interface{}{"", "", "", "", "", "", ""}
		rows.Scan(&aPost[0], &aPost[1], &aPost[2], &aPost[3], &id, &aPost[6], &like, &nbComments)
		// si le RegExp correspond à la DB
		if filter.MatchString(aPost[6].(string)) {
			// Remplace les \n par des <br> pour sauter des lignes en html
			aPost[1] = strings.Replace(aPost[1].(string), string('\r'), "", -1)
			aPost[1] = strings.Replace(aPost[1].(string), string('\n'), "<br>", -1)
			aPost[5] = strconv.Itoa(id)
			if aPost[6] != nil {
				temp := []interface{}{} // string
				for _, e := range aPost[6].(string) {
					j, _ := strconv.Atoi(string(e))
					temp = append(temp, categories[j])
				}
				aPost = append(aPost, temp)
			} else {
				aPost[6] = []string{}
				aPost = append(aPost, []string{})
			}
			aPost = append(aPost, like)
			aPost = append(aPost, nbComments)
			post = append(post, aPost)
		}

	}
	// Ajoute le chemin de la photo qui a été choisit par l'utilisateur
	for i := 0; i < len(post); i++ {
		rows, err := database.Query("SELECT photo FROM users WHERE username = ?", post[i][2])
		checkError(err)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&photo)
			checkError(err)
		}
		err = rows.Err()
		checkError(err)
		post[i][4] = photo
	}

	data_index["allposts"] = post
	data_index["categories"] = categories

	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "index", data_index)
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
	GetCookie(data_SignUp, r)

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

// Generate the log In page
func logIn(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_logIn := make(map[string]interface{})
	for k, v := range data {
		data_logIn[k] = v
	}
	data_logIn["cookieExist"] = false
	GetCookie(data_logIn, r)

	if data_logIn["cookieExist"] == true {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	// get user input to log in
	user_login := r.FormValue("user-login")
	password_login := r.FormValue("password-login")

	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	create_cookie = false
	// Parcourir la BDD
	rows, _ := database.Query("SELECT username, password FROM users")
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
			}
		}
	}

	if create_cookie {
		// Créé un cookie si user bien authentifié
		CreateCookie(w, r)
		data["cookieExist"] = true
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	t := template.New("logIn-template")
	t = template.Must(t.ParseFiles("./html/LogIn.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "LogIn", data_logIn)
}

// Generate the Welcome page
func welcome(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	data_welcome := make(map[string]interface{})
	for k, v := range data {
		data_welcome[k] = v
	}

	GetCookie(data_welcome, r)

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
	GetCookie(data_user, r)

	user := r.FormValue("user")

	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// tx, _ := database.Begin()
	// Parcourir la BDD
	rows, _ := database.Query("SELECT id, username, email, fewWords, address, age, photo FROM users")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id, &username, &email, &fewWords, &address, &age, &photo)
		// Si l'input username est trouvé
		if user == username {
			userFound = true
			data_user["username"] = username
			data_user["email"] = email
			data_user["fewWords"] = fewWords
			data_user["address"] = address
			data_user["age"] = age
			data_user["photo"] = photo
			break
		}
	}
	rows.Close()

	// Check if the user logged is on his personnal page
	if data_user["username"] == data_user["user"] && data_user["cookieExist"] != false {
		data_user["sameUser"] = true
	} else {
		data_user["sameUser"] = false
	}

	// get input age/address/FewWords of the user
	add_few_words := r.FormValue("addFewWords")
	add_age := r.FormValue("age")
	add_address := r.FormValue("address")
	change_photo := r.FormValue("photo")
	// add the input to his data in the DB
	if add_few_words != "" {
		state = "fewWords"
		UpdateInfoUser(database, add_few_words, state, id)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}
	if add_address != "" {
		state = "address"
		UpdateInfoUser(database, add_address, state, id)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}
	if add_age != "" {
		state = "age"
		UpdateInfoUser(database, add_age, state, id)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}
	if change_photo != "" {
		state = "photo"
		UpdateInfoUser(database, change_photo, state, id)
		http.Redirect(w, r, "/user?user="+username, http.StatusSeeOther)
	}

	t := template.New("user-template")
	t = template.Must(t.ParseFiles("./html/user.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "user", data_user)
}

// Generate allUsers page
func allUsers(w http.ResponseWriter, r *http.Request) {
	// initiate the data that will be send to html
	var aUser [3]string
	var all_users [][3]string
	data_allUsers := make(map[string]interface{})
	for k, v := range data {
		data_allUsers[k] = v
	}
	GetCookie(data_allUsers, r)

	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// Parcourir la BDD
	rows, _ := database.Query("SELECT username, email, age FROM users")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&username, &email, &age)
		aUser[0] = username
		aUser[1] = email
		aUser[2] = age
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
	GetCookie(data_newPost, r)
	// Redirection pour ceux qui ne sont pas connecté
	if data_newPost["cookieExist"] == false {
		http.Redirect(w, r, "/logIn", http.StatusSeeOther)
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
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}
	data_newPost["categorie"] = categories

	t := template.New("newPost-template")
	t = template.Must(t.ParseFiles("./html/newPost.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "newPost", data_newPost)
}

// Generate page des posts avec ses commentaires
func post(w http.ResponseWriter, r *http.Request) {
	post_id := r.FormValue("id")

	// initiate the data that will be send to html
	data_post := make(map[string]interface{})
	for k, v := range data {
		data_post[k] = v
	}
	data_post["cookieExist"] = false
	data_post["already_liked"] = false
	GetCookie(data_post, r)

	// Affiche les posts et commentaires
	var post = Display_post(post_id, data_post, body)
	Display_comments(data_post, post_id)

	// Ajoute un commentaire
	add_comment := r.FormValue("add_comment")
	if add_comment != "" {
		Adding_comment(add_comment, &post, data_post["user"].(string))
		http.Redirect(w, r, "/post?id="+post_id, http.StatusSeeOther)
	}

	// Système de Like
	change_nmb_like := r.FormValue("Like")
	likedBy = Like(change_nmb_like, data_post, post_id, w, r)

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

	t := template.New("post-template")
	t = template.Must(t.ParseFiles("./html/post.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "post", data_post)
}

// supprime un post et ses commentaires
func delPost(w http.ResponseWriter, r *http.Request) {
	delete_post := r.FormValue("delPost")

	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

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
	fmt.Println(delete_comment)
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
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
	var all_myPosts [][]interface{}

	// initiate the data that will be send to html
	data_myPosts := make(map[string]interface{})
	for k, v := range data {
		data_myPosts[k] = v
	}
	data_myPosts["cookieExist"] = false
	GetCookie(data_myPosts, r)

	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	//range over database
	rows, _ := database.Query("SELECT title, body, author, date, id, category FROM posts WHERE author = ?", data["user"].(string))
	defer rows.Close()

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

	data_myPosts["all_myPosts"] = all_myPosts

	t := template.New("myPosts-template")
	t = template.Must(t.ParseFiles("./html/myPosts.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "myPosts", data_myPosts)
}

// Page qui affiche les posts liké par l'utilisateur connecté
func myLikedPosts(w http.ResponseWriter, r *http.Request) {
	var all_myLikedPosts [][]interface{}
	var post_liked bool

	// initiate the data that will be send to html
	data_myLikedPosts := make(map[string]interface{})
	for k, v := range data {
		data_myLikedPosts[k] = v
	}
	data_myLikedPosts["cookieExist"] = false
	GetCookie(data_myLikedPosts, r)

	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	//range over database
	rows, _ := database.Query("SELECT title, body, author, date, id, category, likedBy FROM posts")
	defer rows.Close()

	for rows.Next() {
		myLikedPosts := []interface{}{"", "", "", "", "", "", ""}
		rows.Scan(&myLikedPosts[0], &myLikedPosts[1], &myLikedPosts[2], &myLikedPosts[3], &id, &myLikedPosts[6], &likedBy)

		// find all the posts which the user liked
		like_splitted := strings.Split(likedBy, " ")
		for i := 0; i < len(like_splitted); i++ {
			if like_splitted[i] == data["user"].(string) {
				post_liked = true
				break
			} else {
				post_liked = false
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

	data_myLikedPosts["all_myLikedPosts"] = all_myLikedPosts

	t := template.New("myLikedPosts-template")
	t = template.Must(t.ParseFiles("./html/myLikedPosts.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "myLikedPosts", data_myLikedPosts)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
