package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var data []string

const LocalhostPort = ":8000"

func main() {
	fmt.Println("Please connect to\u001b[31m localhost", LocalhostPort, "\u001b[0m")

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.HandleFunc("/", index)
	http.HandleFunc("/signIn", signIn)
	err := http.ListenAndServe(LocalhostPort, nil) // Set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//generate the main page when first loading the site
func index(w http.ResponseWriter, r *http.Request) {
	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "index", data)
}

//generate the main page when first loading the site
func signIn(w http.ResponseWriter, r *http.Request) {
	t := template.New("signIn-template")
	t = template.Must(t.ParseFiles("./html/signIn.html", "./html/header&footer.html"))
	t.ExecuteTemplate(w, "signIn", data)
}
