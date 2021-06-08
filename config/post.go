package config

import (
	"database/sql"
	"net/http"
	"regexp"
	"strconv"
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
	Arrimage     []interface{}
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
	_, err = stmt.Exec(title, body, data_newPost["user"], dt, strings.Join(category, ""))
	CheckError(err)

	tx.Commit()
}

// Affiche les posts pour la page INDEX et pendingPosts
func DisplayPosts(r *http.Request, data_info map[string]interface{}, state string) {
	var categories []string = GetCategories()
	var id int
	var imgstr string
	var arrimg []string
	var all_image [][]string
	var arrimage []interface{}

	// filtre de categorie
	selected_categories := ""
	for i := range categories {
		tmp := r.FormValue(categories[i])
		if tmp != "" {
			selected_categories += tmp
		}
	}
	// RegExp
	var filter = regexp.MustCompile(selected_categories)

	database, err := sql.Open("sqlite3", "./db-sqlite.db")
	CheckError(err)
	defer database.Close()
	//range over database
	var query string
	if state == "pendingPosts" {
		query = "SELECT title, body, author, date, id, category, like, nbComments FROM pendingPosts"
	} else if state == "index" {
		query = "SELECT title, body, author, date, id, category, like, nbComments FROM posts"
	}
	rows, err := database.Query(query)
	CheckError(err)
	defer rows.Close()
	var post []Post
	for rows.Next() {
		var aPost Post
		rows.Scan(&aPost.Title, &aPost.Body, &aPost.Author, &aPost.Date, &aPost.Id, &aPost.Category, &aPost.Like, &aPost.NbComments)
		// si le RegExp correspond à la DB
		if filter.MatchString(aPost.Category) {
			// Remplace les \n par des <br> pour sauter des lignes en html
			aPost.Body = strings.Replace(strings.Replace(aPost.Body, string('\r'), "", -1), string('\n'), "<br>", -1)

			if aPost.Category != "" {
				for _, e := range aPost.Category {
					j, _ := strconv.Atoi(string(e))
					if j >= len(categories) {
						continue
					}
					aPost.Category_Tab = append(aPost.Category_Tab, categories[j])
				}
				//
				if len(aPost.Category_Tab) == 0 {
					continue
				}
			}
			post = append(post, aPost)
		}

		// Ajoute les images pour la page pendingPosts seulement
		if state == "pendingPosts" {
			rows, err := database.Query("SELECT image FROM pendingPosts WHERE id = ?", id)
			CheckError(err)
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&imgstr)
				CheckError(err)
			}
			arrimg = strings.Split(imgstr, ",")
			arrimg = arrimg[1:]
			all_image = append(all_image, arrimg)
			for i := 0; i < len(all_image); i++ {
				if i == len(all_image)-1 {
					arrimage = append(arrimage, all_image[i])

				}
			}
			CheckError(rows.Err())
		}
		if len(arrimage) > 0 {
			aPost.Arrimage = arrimage[len(arrimage)-1:]
		}
	}

	// Ajoute le chemin de la photo qui a été choisit par l'utilisateur
	for i := 0; i < len(post); i++ {
		rows, err := database.Query("SELECT photo FROM users WHERE username = ?", post[i].Photo)
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
