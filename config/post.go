package config

import (
	"database/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	. "./err"
)

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
	var categories = []string{"gaming", "informatique", "sport", "culture", "politique", "loisir", "sciences", "sexualite", "finance"}
	var like, nbComments, photo string
	var id int
	var post [][]interface{}

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
	var query string
	if state == "pendingPosts" {
		query = "SELECT title, body, author, date, id, category, like, nbComments FROM pendingPosts"
	} else if state == "index" {
		query = "SELECT title, body, author, date, id, category, like, nbComments FROM posts"
	}
	rows, _ := database.Query(query)
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
		CheckError(err)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&photo)
			CheckError(err)
		}
		err = rows.Err()
		CheckError(err)
		post[i][4] = photo
	}

	data_info["allposts"] = post
	data_info["categories"] = categories
}

func Display_post_comment(post_id string, data_post map[string]interface{}, body string) [8]string {
	var post [8]string
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	//range over database
	rows, _ := database.Query("SELECT id, title, body, author, date, like FROM posts WHERE id = ?", post_id)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&post[0], &post[1], &body, &post[3], &post[4], &post[6])
		CheckError(err)
	}
	// Remplace les \n par des <br> pour sauter des lignes en html
	post[2] = strings.Replace(body, string('\r'), "", -1)
	post[2] = strings.Replace(body, string('\n'), "<br>", -1)
	// Ajoute le chemin de la photo qui a été choisit par l'utilisateur
	rows, err := database.Query("SELECT photo FROM users WHERE username = ?", post[3])
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
