package config

import (
	"database/sql"
	"net/http"
)

type Categorie struct {
	Nom   string
	Actif bool
}

//gaming informatique sport culture polititque loisir sciences sexualite finance
func InitCategoriePrincipale() {
	categories := []string{"gaming", "informatique", "sport", "culture", "politique", "loisir", "sciences", "sexualite", "finance"}
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	//
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
	CheckError(rows.Err())
	for i := range categories {
		// add the category
		tx, err := database.Begin()
		CheckError(err)
		stmt, err := tx.Prepare("INSERT INTO categories (categorie, actif) VALUES (?, 0)")
		CheckError(err)
		_, err = stmt.Exec(categories[i])
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
	rows, err := database.Query("SELECT categorie, actif FROM categories")
	if err != nil {
		return resultat
	}
	for rows.Next() {
		var tmp Categorie
		rows.Scan(&tmp.Nom, &tmp.Actif)
		resultat = append(resultat, tmp)
	}
	rows.Close()
	return resultat
}
func GetCategories() []string {
	categorie := DisplayCategories()
	var res []string
	for i := range categorie {
		if categorie[i].Actif {
			res = append(res, categorie[i].Nom)
		}
	}
	return res
}
func NewCategorie(w http.ResponseWriter, r *http.Request) {
	nom_categorie := r.FormValue("newCategorie")
	if nom_categorie == "" {
		return
	}
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	// range over the database and check if there is double categorie
	rows, err := database.Query("SELECT categorie FROM categories")
	CheckError(err)
	for rows.Next() {
		tmp := ""
		rows.Scan(&tmp)
		//stop the function if a double is found
		if nom_categorie == tmp {
			rows.Close()
			return
		}
	}
	rows.Close()
	// add the category
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare("INSERT INTO categories (categorie, actif) VALUES (?, 0)")
	CheckError(err)
	_, err = stmt.Exec(nom_categorie)
	CheckError(err)
	tx.Commit()
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}
func ActiverCategorie(w http.ResponseWriter, r *http.Request) {
	tmp := func() string {
		categorie := DisplayCategories()
		for i := range categorie {
			categorie_selec := r.FormValue(categorie[i].Nom + "_false")
			if categorie_selec != "" {
				return categorie_selec
			}
		}
		return ""
	}

	categorie_selec := tmp()
	if categorie_selec == "" {
		return
	}
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare("UPDATE categories SET actif = ? WHERE categorie = ?")
	CheckError(err)
	_, err = stmt.Exec(1, categorie_selec)
	CheckError(err)
	tx.Commit()
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)

}
func DesactiverCategorie(w http.ResponseWriter, r *http.Request) {
	tmp := func() string {
		categorie := DisplayCategories()
		for i := range categorie {
			categorie_selec := r.FormValue(categorie[i].Nom + "_true")
			if categorie_selec != "" {
				return categorie_selec
			}
		}
		return ""
	}
	categorie_selec := tmp()
	if categorie_selec == "" {
		return
	}
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare("UPDATE categories SET actif = ? WHERE categorie = ?")
	CheckError(err)
	_, err = stmt.Exec(0, categorie_selec)
	CheckError(err)
	tx.Commit()
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)

}
func RenommerCategorie(w http.ResponseWriter, r *http.Request) {
	// name , new name
	tmp := func() (string, string) {
		categorie := DisplayCategories()
		for i := range categorie {
			categorie_selec := r.FormValue(categorie[i].Nom + "_DEL")
			if categorie_selec != "" {
				return categorie_selec, categorie[i].Nom
			}
		}
		return "", ""
	}
	categorie_selec, categorie := tmp()
	if categorie_selec == "" {
		return
	}
	// Open the database
	database, _ := sql.Open("sqlite3", "./db-sqlite.db")
	defer database.Close()
	// UPDATE the categorie from de table categories
	tx, err := database.Begin()
	CheckError(err)
	stmt, err := tx.Prepare("UPDATE categories SET categorie = ? WHERE categorie = ?")
	CheckError(err)
	_, err = stmt.Exec(categorie_selec, categorie)
	CheckError(err)
	tx.Commit()
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)

}
