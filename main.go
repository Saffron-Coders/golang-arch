package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// key is github id, value is user Id
var githubConnection map[string]string

// {"data":{"viewer":{"id":"...."}}}
type githubResponse struct {
	Data struct {
		Viewer struct {
			Id string `json:"id"`
		} `json:"viewer"`
	} `json:"data"`
}

var githubOauthConfig = &oauth2.Config{
	ClientID:     "52a1d5b647d23713a8b8",
	ClientSecret: "098838fa720f56d6245061c387bc1d9f3da1636a",
	Endpoint:     github.Endpoint,
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/oauth/github", startGithubOauth)
	http.HandleFunc("/oauth2/receive", completeGithubOauth)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Oauth2</title>
	</head>
	<body>
			<form action="/oauth/github" method="POST">
				<input type="submit", value="Login with github">
			</form>
	</body>
	</html>`)
}

func startGithubOauth(w http.ResponseWriter, r *http.Request) {
	redirectURL := githubOauthConfig.AuthCodeURL("0000")
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func completeGithubOauth(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	state := r.FormValue("state")

	if state != "0000" {
		http.Error(w, "State is incorrect", http.StatusBadRequest)
		return
	}

	token, err := githubOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Couldn't login", http.StatusInternalServerError)
		return
	}

	ts := githubOauthConfig.TokenSource(r.Context(), token)
	client := oauth2.NewClient(r.Context(), ts)

	requestBody := strings.NewReader(`{"query": "query {viewer {id}}"}`)

	resp, err := client.Post("https://api.github.com/graphql", "application/json", requestBody)
	if err != nil {
		http.Error(w, "Couldn't get user", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var gr githubResponse

	err = json.NewDecoder(resp.Body).Decode(&gr)
	if err != nil {
		http.Error(w, "Github invalid response", http.StatusInternalServerError)
	}

	githubID := gr.Data.Viewer.Id
	userId, ok := githubConnection[githubID]
	if !ok {
		// New User Create account
		// Maybe return may be not, it depends
	}

	// Login to account userId using JWT
}
