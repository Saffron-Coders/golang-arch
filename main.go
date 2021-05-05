package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"net/http"
)

func getCode(msg string) string {
	hmac.New(sha256.New, []byte(""))
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

	email := r.FormValue("email")
	if email == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	c := http.Cookie{
		Name:  "session",
		Value: "",
	}
}

func foo(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>HMAC Example</title>
	</head>
	<body>
		<form action="/submit" method="POST">
			<input type="text" name="email"/>
			<input type="submit" />
		</form>
	</body>
	</html>`

	io.WriteString(w, html)
}
