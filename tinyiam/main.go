package main

import (
	"fmt"
	"html/template"
	"net/http"

	"log"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gitlab"
)

func main() {

	key := "Secret-session-key" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30        // 30 days
	isProd := false             // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	// goth.UseProviders(
	// 	google.New("our-google-client-id", "our-google-client-secret", "http://localhost:3000/auth/google/callback", "email", "profile"),
	// )

	gitlab.AuthURL = "http://gitlab.dev-git.podinate.com:8081/oauth/authorize"
	gitlab.TokenURL = "http://gitlab.dev-git.podinate.com:8081/oauth/token"
	gitlab.ProfileURL = "http://gitlab.dev-git.podinate.com:8081/api/v3/user"

	goth.UseProviders(
		gitlab.New("acd228dbd320a2be27f152d01944519df127f0d4e739d24b24171a54f1aeb6bc", "e834e1f3261806bc7da403e426374a36a06b9b4eb5625b49f95de7642a37ba35", "http://localhost:3002/auth/gitlab/callback", "read_user", "read_repository", "read_registry", "profile", "email"),
	)

	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			fmt.Println("Error", err)
			return
		}
		fmt.Printf("Success %+v\n", user)
		t, err := template.ParseFiles("templates/success.html")
		if err != nil {
			fmt.Println("Error parsing template", err)
			return
		}
		err = t.Execute(res, user)
		if err != nil {
			fmt.Println("Error executing template", err)
			return
		}
		fmt.Println("What the fuck is happening")
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})
	log.Println("listening on localhost:3002")
	log.Fatal(http.ListenAndServe(":3002", p))
}
