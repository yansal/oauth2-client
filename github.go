package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func newGitHubCallbackHandler() (*githubCallbackHandler, error) {
	return &githubCallbackHandler{
		clientID:     os.Getenv("GITHUB_CLIENT_ID"),
		clientsecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}, nil
}

type githubCallbackHandler struct {
	clientID, clientsecret string
}

func (h *githubCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.serveHTTP(w, r)
	if err == nil {
		return
	}
	log.Print(err)
}

func (h *githubCallbackHandler) serveHTTP(w http.ResponseWriter, r *http.Request) error {
	body := make(url.Values)
	body.Set("client_id", h.clientID)
	body.Set("client_secret", h.clientsecret)
	body.Set("code", r.FormValue("code"))
	req, err := http.NewRequest(
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		strings.NewReader(body.Encode()),
	)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("%s", b)
	return nil
}
