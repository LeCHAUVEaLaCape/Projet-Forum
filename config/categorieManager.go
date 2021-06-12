package config

import (
	"database/sql"
)

type Categorie struct {
	Nom   string
	Color string
	Id    string
}

//gaming informatique sport culture polititque loisir sciences sexualite finance
func InitCategoriePrincipale() {
	categories := []string{"gaming", "informatique", "sport", "culture", "politique", "loisir", "sciences", "sexualite", "finance"}
	color_categories := []string{"#ff0000", "#6485ff", "#0de20d", "#ffff00", "#9e9e9e", "#00ffff", "#ffa500", "#e44dff", "#8c9638"}

	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// Check if the table is already filled
	rows, err := database.Query("SELECT categorie FROM categories WHERE categorie = ?", "gaming")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		tmp := ""
		err := rows.Scan(&tmp)
		CheckError(err)
		if tmp == "gaming" {
			return
		}
	}

	for i := range categories {
		// add the category
		tx, err := database.Begin()
		CheckError(err)
		stmt, err := tx.Prepare("INSERT INTO categories (categorie, color) VALUES (?, ?)")
		CheckError(err)
		_, err = stmt.Exec(categories[i], color_categories[i])
		CheckError(err)
		tx.Commit()
	}
}

func DisplayCategories() []Categorie {
	var resultat []Categorie
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// range over the database and check if there is double categorie
	rows, err := database.Query("SELECT categorie, color, id FROM categories")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		var tmp Categorie
		rows.Scan(&tmp.Nom, &tmp.Color, &tmp.Id)
		resultat = append(resultat, tmp)
	}
	return resultat
}

func GetCategories() [][3]string {
	categorie := DisplayCategories()
	var categories [3]string
	var all_categories [][3]string
	for i := range categorie {
		cat := categorie[i].Nom
		col := categorie[i].Color
		categories[0] = cat
		categories[1] = col
		categories[2] = col
		all_categories = append(all_categories, categories)
	}
	return all_categories
}

func NewCategorie(nom_categorie string, color_cat string) {

	if nom_categorie == "" {
		return
	}
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	// range over the database and check if there is double categorie
	rows, err := database.Query("SELECT categorie FROM categories")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		tmp := ""
		rows.Scan(&tmp)
		//stop the function if a double is found
		if nom_categorie == tmp {
			rows.Close()
			return
		}
	}

	// add the category
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare("INSERT INTO categories (categorie, color) VALUES (?, ?)")
	CheckError(err)
	_, err = stmt.Exec(nom_categorie, color_cat)
	CheckError(err)
	tx.Commit()
}

// return category listed with whitespace when the category isn't "actif"
func GetBruteCategories() []string {
	categorie := DisplayCategories()
	var res []string
	for i := range categorie {
		res = append(res, categorie[i].Nom)
	}
	return res
}

func Del_Categorie(del_categorie string) {
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()

	tx, err := database.Begin()

	stmt, err := tx.Prepare("DELETE FROM categories WHERE id = ?")
	CheckError(err)
	_, err = stmt.Exec(del_categorie)
	CheckError(err)

	tx.Commit()
}
