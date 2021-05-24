package cookies

import (
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Create a cookie
func CreateCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		id, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:     "session",   // nom du cookie
			Value:    id.String(), // uuid pour le cookie
			HttpOnly: true,        // protection pour que le cookie ne soit pas visible par le JS
			Path:     "/",         // cookie valable de puis la racine du serveur
			MaxAge:   60 * 5,      // cookie valable 5 minutes
		}
		http.SetCookie(w, cookie)
	}
	fmt.Println(cookie)
}
