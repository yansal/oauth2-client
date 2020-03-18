package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	log.SetFlags(0)

	facebookCallback, err := newFacebookCallbackHandler()
	if err != nil {
		log.Fatal(err)
	}
	githubCallback, err := newGitHubCallbackHandler()
	if err != nil {
		log.Fatal(err)
	}

	root, err := newRootHandler()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/facebook/callback", facebookCallback)
	mux.Handle("/github/callback", githubCallback)
	mux.Handle("/", root)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), mux))
}

func newRootHandler() (*rootHandler, error) {
	facebookAuthURL, err := url.Parse("https://www.facebook.com/v6.0/dialog/oauth")
	if err != nil {
		return nil, err
	}
	q := facebookAuthURL.Query()
	q.Set("client_id", os.Getenv("FACEBOOK_CLIENT_ID"))
	q.Set("redirect_uri", os.Getenv("FACEBOOK_REDIRECT_URI"))
	q.Set("scope", "email")
	facebookAuthURL.RawQuery = q.Encode()

	githubAuthURL, err := url.Parse("https://github.com/login/oauth/authorize")
	if err != nil {
		return nil, err
	}
	q = githubAuthURL.Query()
	q.Set("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	githubAuthURL.RawQuery = q.Encode()

	tmpl, err := template.ParseGlob("tmpl/*")
	if err != nil {
		return nil, err
	}
	return &rootHandler{
		facebookAuthURL: facebookAuthURL.String(),
		githubAuthURL:   githubAuthURL.String(),
		tmpl:            tmpl,
	}, nil
}

type rootHandler struct {
	facebookAuthURL string
	githubAuthURL   string
	tmpl            *template.Template
}

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.Execute(w, struct {
		FacebookAuthURL string
		GitHubAuthURL   string
	}{
		FacebookAuthURL: h.facebookAuthURL,
		GitHubAuthURL:   h.githubAuthURL,
	}); err != nil {
		log.Print(err)
	}
}
