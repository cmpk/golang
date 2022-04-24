package authentication

import (
	"backend/singleton"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const OIDC_CALLBACK_PATH = "/api/set-token"

func CheckTokenAndLogin(w http.ResponseWriter, req *http.Request) {
	// トークンがセットされていない場合は、セッションからログイン状態を確認する
	var err error

	sessionManager := singleton.GetSessionManager()
	session := sessionManager.SessionStart(w, req)
	rawIDToken, err := session.Get("token")
	if err != nil {
		authenticate(w, req) // Keycloak にリダイレクト
		return
	}

	_, err = checkToken(w, req, rawIDToken)
	if err != nil {
		authenticate(w, req) // Keycloak にリダイレクト
		return
	}
}

func CheckTokenMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Print("===== START : CheckTokenMiddleware")

		if req.Header["Authorization"] != nil {
			//TODO フロントエンドで認証している場合はトークンを検証する
			token := req.Header["Authorization"]
			log.Printf("token = %s", token)
		} else {
			// Backend 側で Keycloak 認証を行う
			CheckTokenAndLogin(w, req)
		}

		handler.ServeHTTP(w, req)
		log.Print("===== END : CheckTokenMiddleware")
	})
}

var LogoutHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	//TODO Keycloak からのログアウト
	sessionManager := singleton.GetSessionManager()
	sessionManager.SessionDestroy(w, req)

	http.Redirect(w, req, "/api", http.StatusFound)
})

var SetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	var err error

	config, _ := getConfig()
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

	// アプリケーションで利用する情報をセッションに保存
	idTokenClaims := map[string]interface{}{}
	if err = idToken.Claims(&idTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("idTokenClaims = %#v", idTokenClaims)
	sessionManager := singleton.GetSessionManager()
	session := sessionManager.SessionStart(w, req)
	session.Set("uid", fmt.Sprint(idTokenClaims["sub"]))
	session.Set("token", rawIDToken)

	previous, _ := session.Get("previous")
	log.Printf("===== previous = " + previous)
	log.Printf("===== sessionId = " + session.SessionID())
	session.Delete("previous")

	http.Redirect(w, req, previous, http.StatusFound)
})

func authenticate(w http.ResponseWriter, req *http.Request) {
	config, _ := getConfig()
	url := config.AuthCodeURL(os.Getenv("STATE_STRING"))

	// Keycloak ログイン後、リクエストURIにリダイレクトするために、URIをセッションに保存
	sessionManager := singleton.GetSessionManager()
	session := sessionManager.SessionStart(w, req)
	session.Set("previous", req.RequestURI)

	//TODO
	previous, _ := session.Get("previous")
	log.Printf("===== authenticate previous = " + previous)
	log.Printf("===== authenticate sessionId = " + session.SessionID())

	http.Redirect(w, req, url, http.StatusFound)
}

func checkToken(w http.ResponseWriter, req *http.Request, token string) (*oidc.IDToken, error) {
	var err error

	_, provider := getConfig()
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

func getConfig() (*oauth2.Config, *oidc.Provider) {
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
