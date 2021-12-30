package authentication

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const OIDC_CALLBACK_PATH = "/api/set-token"

func CheckTokenHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var err error

		rawIDTokenCookie, err := req.Cookie("token")
		if err != nil {
			authenticate(w, req) // Keycloak にリダイレクト
			return
		}

		_, err = checkToken(w, req, rawIDTokenCookie.Value)
		if err != nil {
			authenticate(w, req) // Keycloak にリダイレクト
			return
		}

		handler.ServeHTTP(w, req)
	})
}

var SetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	var err error

	config, _ := getConfig(req)
	if err = req.ParseForm(); err != nil {
		http.Error(w, "parse form error", http.StatusInternalServerError)
		return
	}

	accessToken, err := config.Exchange(context.Background(), req.Form.Get("code"))
	if err != nil {
		http.Error(w, "Can't get access token", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := accessToken.Extra("id_token").(string)
	if !ok {
		http.Error(w, "missing token", http.StatusInternalServerError)
		return
	}

	idToken, err := checkToken(w, req, rawIDToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// IDトークンのクレームをとりあえずダンプ
	// アプリで必要なものはセッションストレージに入れておくと良いでしょう
	idTokenClaims := map[string]interface{}{}
	if err = idToken.Claims(&idTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("idTokenClaims = %#v", idTokenClaims)
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: rawIDToken, // 行儀が悪いので真似しないねで
		Path:  "/api",
	})
	http.Redirect(w, req, "/api", http.StatusFound)
})

func authenticate(w http.ResponseWriter, req *http.Request) {
	config, _ := getConfig(req)
	url := config.AuthCodeURL(os.Getenv("STATE_STRING"))

	http.Redirect(w, req, url, http.StatusFound)
}

func checkToken(w http.ResponseWriter, req *http.Request, token string) (*oidc.IDToken, error) {
	var err error

	_, provider := getConfig(req)
	if err = req.ParseForm(); err != nil {
		return nil, errors.New("parse form error")
	}
	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("CLIENT_ID"),
	}
	verifier := provider.Verifier(oidcConfig)
	idToken, err := verifier.Verify(context.Background(), token)
	if err != nil {
		log.Printf("id token verify error : " + err.Error())
		return nil, err
	}

	return idToken, nil
}

var once sync.Once
var provider *oidc.Provider
var oauth2Config *oauth2.Config

func getConfig(req *http.Request) (*oauth2.Config, *oidc.Provider) {
	once.Do(func() {
		// See : https://github.com/coreos/go-oidc
		var err error

		provider, err = oidc.NewProvider(context.Background(), os.Getenv("AUTH_URL")+"/"+os.Getenv("APP_NAME"))
		if err != nil {
			panic(err)
		}
		oauth2Config = &oauth2.Config{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID},
			RedirectURL:  os.Getenv("APP_URL") + OIDC_CALLBACK_PATH,
		}
	})
	return oauth2Config, provider
}
