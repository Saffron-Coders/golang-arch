package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	db  = map[string][]byte{}
	key = []byte("my secret key 007 james bond rule the world from my mom's basement")
)

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)

	http.ListenAndServe(":8080", nil)
}

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

func createToken(sId string) string {

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(sId))

	// To hex
	// signedMac := fmt.Sprintf("%x", mac.Sum(nil))

	// To base64
	signedMac := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// SignedSeesionId as base64 | created from sid
	return signedMac + "|" + sId
}

func parseToken(ss string) (string, error) {
	xs := strings.SplitN(ss, "|", 2)
	if len(xs) != 2 {
		return "", fmt.Errorf("stop hacking me")
	}
	// signedSessionId as BASE64 | created from sid
	b64 := xs[0]
	bx, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("couldn't parse token decode string: %w", err)
	}

	// signedSessionId as BASE64 | created from sid
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(xs[1]))

	ok := hmac.Equal(bx, mac.Sum(nil))
	if !ok {
		return "", fmt.Errorf("couldn't pasreToken not equal signed sid and session id: %w", err)
	}

	return xs[1], nil

}
