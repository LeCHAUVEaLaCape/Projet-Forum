package config

import (
	"database/sql"
	"net/http"
	"strings"
)

type Post struct {
	Id           string `json:"id"`
	Title        string
	Body         string
	Like         string
	Author       string
	Date         string
	Category     string
	Category_Tab []string
	LikedBy      string
	NbComments   string
	Photo        string
	Arrimage     []string
}

func AddNewPost(title string, body string, dt string, data_newPost map[string]interface{}, category []string) {
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	// add the inputs to the database
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare("INSERT INTO pendingPosts (title, body, author, date, category, like, likedBy, nbComments, dislikedBy, image) VALUES (?, ?, ?, ?, ?, 0, '', 0, '', '')")
	CheckError(err)
	_, err = stmt.Exec(title, body, data_newPost["user"], dt, strings.Join(category, ","))
	CheckError(err)

	tx.Commit()
}

// Affiche les posts pour la page INDEX et pendingPosts
func DisplayPosts(r *http.Request, data_info map[string]interface{}, state string) {
	categories := GetCategories()

	var imgstr string
	var arrimg []string

	// filtre de categorie

	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()
	//range over database
	var query string
	if state == "pendingPosts" {
		query = "SELECT title, body, author, date, id, category, like, nbComments, image FROM pendingPosts"
	} else if state == "index" {
		query = "SELECT title, body, author, date, id, category, like, nbComments, image FROM posts"
	}
	rows, err := database.Query(query)
	CheckError(err)
	defer rows.Close()
	var post []Post
	for rows.Next() {
		var aPost Post
		rows.Scan(&aPost.Title, &aPost.Body, &aPost.Author, &aPost.Date, &aPost.Id, &aPost.Category, &aPost.Like, &aPost.NbComments, &imgstr)
		// Remplace les \n par des <br> pour sauter des lignes en html
		aPost.Body = strings.Replace(strings.Replace(aPost.Body, string('\r'), "", -1), string('\n'), "<br>", -1)

		// Split pour avoir toutes les catégories du post
		categorie_post := strings.Split((aPost.Category), ",")

		filtered := true
		// When user wants to filter the posts
		if r.FormValue("categorie") != "" {
			for i := 0; i < len(categorie_post); i++ { // range over each categories of all posts
				if categorie_post[i] == r.FormValue("categorie") {
					break // when one category is the one the user is looking for
				} else if categorie_post[i] != r.FormValue("categorie") && i == len(categorie_post)-1 {
					filtered = false // if none of the post's categories match the user was looking for
				}
			}
		}

		// ajoute toutes les catégories dans un array
		if aPost.Category != "" {
			for i := 0; i < len(categories); i++ {
				for j := 0; j < len(categorie_post); j++ {
					// Compare les catégories et ajoute celles qui correspondent
					if categories[i][0] == categorie_post[j] {
						aPost.Category_Tab = append(aPost.Category_Tab, categorie_post[j])
					}
				}
			}
		}

		// Ajoute les images pour la page pendingPosts seulement
		arrimg = strings.Split(imgstr, ",")
		arrimg = arrimg[1:]
		aPost.Arrimage = arrimg

		if filtered {
			post = append(post, aPost)
		}
	}

	for i := 0; i < len(post); i++ {
		rows, err := database.Query("SELECT photo FROM users WHERE username = ?", post[i].Author)
		CheckError(err)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&post[i].Photo)
			CheckError(err)
		}
		err = rows.Err()
		CheckError(err)
	}

	data_info["allposts"] = post
	data_info["categories"] = categories
}

// Page posts
func Display_post_comment(post_id string, data_post map[string]interface{}) [8]string {
	var post [8]string
	var body string
	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()
	//range over database
	rows, err := database.Query("SELECT id, title, body, author, date, like, dislike FROM posts WHERE id = ?", post_id)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&post[0], &post[1], &body, &post[3], &post[4], &post[6], &post[7])
		CheckError(err)
	}
	// Remplace les \n par des <br> pour sauter des lignes en html
	post[2] = strings.Replace(body, string('\r'), "", -1)
	post[2] = strings.Replace(body, string('\n'), "<br>", -1)
	// Ajoute le chemin de la photo qui a été choisit par l'utilisateur
	rows, err = database.Query("SELECT photo FROM users WHERE username = ?", post[3])
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&post[5])
		CheckError(err)
	}
	err = rows.Err()
	CheckError(err)

	var imgstr string
	var arrimg []string
	var all_image [][]string
	rows, err = database.Query("SELECT image FROM posts WHERE id = ?", post_id)
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&imgstr)
		CheckError(err)
	}
	arrimg = strings.Split(imgstr, ",")
	arrimg = arrimg[1:]
	all_image = append(all_image, arrimg)
	err = rows.Err()
	CheckError(err)
	data_post["main_post"] = post
	data_post["imagepost"] = all_image
	return post
}
