package cookies

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Verify if a cookie exist, add a boolean Value to the map
func GetCookie(w http.ResponseWriter, data_info map[string]interface{}, r *http.Request) {
	// get cookie
	cookie, err := r.Cookie("session")
	if err == nil && data_info["user"] != "" {
		data_info["cookieExist"] = true
		// Recréé un cookie pour réinitialiser le temps d'expéritation
		id, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:     "session",   // nom du cookie
			Value:    id.String(), // uuid pour le cookie
			HttpOnly: true,        // protection pour que le cookie ne soit pas visible par le JS
			Path:     "/",         // cookie valable de puis la racine du serveur
			MaxAge:   60 * 5,      // cookie valable 5 minutes
		}
		http.SetCookie(w, cookie)
	} else {
		// Vide tous les données concernant le précédent utilisateur connecté
		data_info["cookieExist"] = false
		data_info["username"] = ""
		data_info["role"] = ""
		data_info["user"] = ""

	}
}
