package cookies

import (
	"net/http"
)

// Verify if a cookie exist, add a boolean Value to the map
func GetCookie(data_info map[string]interface{}, r *http.Request) {
	// get cookie
	_, err := r.Cookie("session")
	if err == nil && data_info["user"] != "" {
		data_info["cookieExist"] = true
	} else {
		data_info["cookieExist"] = false
		data_info["username"] = ""
		data_info["role"] = ""
		data_info["user"] = ""
	}
}
