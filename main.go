package main

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)

	http.ListenAndServe(":8080", nil)
}

var db map[string][]byte

func index(w http.ResponseWriter, r *http.Request) {
	errMsg := r.FormValue("errormsg")
	html := `<!DOCTYPE html>
	<html>
	<body>
	
	<h2>HTML Forms</h2>
	
	<form action="/register" method="POST">
	  <label for="fname">Username:</label><br>
	  <input type="text" id="username" name="username" ><br>

	  <label for="password">Password:</label><br>
	  <input type="password" id="password" name="password" value=""><br><br>
	  <input type="submit" value="Submit">
	</form> 
	
	<p>If you click the "Submit" button, the form-data will be sent to the db</p>
	<p> Error: %s </p>
	</body>
	</html>
	`
	fmt.Fprintf(w, html, errMsg)
}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		errorMsg := url.QueryEscape("Method isn't Post")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" {
		errorMsg := url.QueryEscape("username is required")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}
	if password == "" {
		errorMsg := url.QueryEscape("password is required")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errorMsg := url.QueryEscape("internal server error")
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	db[username] = hash
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
