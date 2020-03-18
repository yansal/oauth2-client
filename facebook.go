package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func newFacebookCallbackHandler() (*facebookCallbackHandler, error) {
	return &facebookCallbackHandler{
		clientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
		clientsecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		redirectURI:  os.Getenv("FACEBOOK_REDIRECT_URI"),
	}, nil
}

type facebookCallbackHandler struct {
	clientID, clientsecret string
	redirectURI            string
}

func (h *facebookCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.serveHTTP(w, r)
	if err == nil {
		return
	}
	log.Print(err)
}

func (h *facebookCallbackHandler) serveHTTP(w http.ResponseWriter, r *http.Request) error {
	body := make(url.Values)
	body.Set("client_id", h.clientID)
	body.Set("client_secret", h.clientsecret)
	body.Set("redirect_uri", h.redirectURI)
	body.Set("code", r.FormValue("code"))
	req, err := http.NewRequest(
		http.MethodPost,
		"https://graph.facebook.com/v6.0/oauth/access_token",
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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("%s", data)

	var v struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	me, err := url.Parse("https://graph.facebook.com/me")
	if err != nil {
		return err
	}
	q := me.Query()
	q.Set("access_token", v.AccessToken)
	q.Set("fields", "id,name,email")
	me.RawQuery = q.Encode()

	http.Redirect(w, r, me.String(), http.StatusFound)
	return nil
}
