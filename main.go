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
	http.HandleFunc("/login", login)

	http.ListenAndServe(":8080", nil)
}

var db = map[string][]byte{}

func index(w http.ResponseWriter, r *http.Request) {
	errMsg := r.FormValue("msg")
	html := `<!DOCTYPE html>
	<html>
	<body>
	
	<h2>HTML Forms</h2>
	
	<form action="/register" method="POST">
	  <label for="fname">Username:</label><br>
	  <input type="text" id="username" name="username" ><br>

	  <label for="password">Password:</label><br>
	  <input type="password" id="password" name="password" value=""><br><br>
	  <input type="submit" value="register">
	</form> 
	<br><br><br><br><br><br>
	<form action="/login" method="POST">
	  <label for="fname">Username:</label><br>
	  <input type="text" id="username" name="username" ><br>

	  <label for="password">Password:</label><br>
	  <input type="password" id="password" name="password" value=""><br><br>
	  <input type="submit" value="login">
	</form> 
	
	
	<p> Error: %s </p>
	</body>
	</html>
	`
	fmt.Fprintf(w, html, errMsg)
}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		msg := url.QueryEscape("Method isn't Post")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" {
		msg := url.QueryEscape("username is required")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	if password == "" {
		msg := url.QueryEscape("password is required")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		msg := url.QueryEscape("internal server error")
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	db[username] = hash
	fmt.Println(db)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		msg := url.QueryEscape("Method isn't Post")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" {
		msg := url.QueryEscape("username is required")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	if password == "" {
		msg := url.QueryEscape("password is required")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	if _, ok := db[username]; !ok {
		msg := url.QueryEscape("email password didn't match")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(db[username]), []byte(password))
	if err != nil {
		msg := url.QueryEscape("email password didn't match")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	msg := url.QueryEscape("Logged in as: " + username)
	http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
}
