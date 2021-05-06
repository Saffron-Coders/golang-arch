package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func getJWT(msg string) (string, error) {
	myKey := "i love thursday when it rains 8723 inches"

	type myClaims struct {
		jwt.StandardClaims
		email string
	}

	claims := myClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		email: msg,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	ss, err := token.SignedString(myKey)
	if err != nil {
		return "", fmt.Errorf("couldn't signString: %w", err)
	}

	return ss, nil
}

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/submit", bar)

	http.ListenAndServe(":8080", nil)
}

func bar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("emailThing")
	if email == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ss, err := getJWT(email)
	if err != nil {
		http.Error(w, "couldn't get jwt", http.StatusInternalServerError)
	}

	c := http.Cookie{
		Name:  "session",
		Value: ss,
	}

	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func foo(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		c = &http.Cookie{}
	}

	// isEqual := true

	message := "Not logged in"
	if isEqual {
		message = "Logged in"
	}

	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>HMAC Example</title>
	</head>
	<body>
	<p>Cookie Value: ` + c.Value + `</p>
	<p>Cookie Value: ` + message + `</p>
		<form action="/submit" method="POST">
			<input type="text" name="emailThing"/>
			<input type="submit" />
		</form>
	</body>
	</html>`

	io.WriteString(w, html)
}
