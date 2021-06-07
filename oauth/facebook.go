package authentification

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var oauthConf = &oauth2.Config{
	ClientID:     "166859168782629",
	ClientSecret: "d7fc1ba6430a4c122b41caf0282fea44",
	RedirectURL:  "http://localhost:8000/FBloginCallBack",
	Scopes:       []string{"public_profile", "email"},
	Endpoint:     facebook.Endpoint,
}

type Info struct {
	Name  string `json:"name"`
	Id    int    `json:"id"`
	Email string `json:"email"`
}

var info Info
var FB_google_user = &info

func HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	Url, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("oauthConf.Exchange() failed with \n", err)
		http.Redirect(w, r, "/logIn", http.StatusTemporaryRedirect)
		return
	}
	resp, err := http.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		fmt.Println("Get: \n", err)
		http.Redirect(w, r, "/logIn", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll: \n", err)
		http.Redirect(w, r, "/logIn", http.StatusTemporaryRedirect)
		return
	}
	json.Unmarshal(response, &info)
	info.Name = strings.ReplaceAll(info.Name, " ", "_")

	http.Redirect(w, r, "/logIn", http.StatusTemporaryRedirect)
}
