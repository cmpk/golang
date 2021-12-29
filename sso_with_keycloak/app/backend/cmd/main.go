//
// 現時点では以下のコードをほぼそのまま利用している。
// [OpenID Connectを使ったアプリケーションのテストのためにKeycloakを使ってみる](https://qiita.com/shibukawa/items/fd78d1ca6c23ce2fa8df)
//
package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/handlers"
	"golang.org/x/oauth2"
)

var once sync.Once

var provider *oidc.Provider
var oauth2Config *oauth2.Config

const OIDC_CALLBACK_PATH = "/api/set-token"

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

func authenticate(w http.ResponseWriter, req *http.Request) {
	config, _ := getConfig(req)
	url := config.AuthCodeURL(os.Getenv("STATE_STRING"))

	http.Redirect(w, req, url, http.StatusFound)
}

func checkToken(w http.ResponseWriter, req *http.Request, token string) (*oidc.IDToken, error) {
	log.Printf("===== checkToken")
	_, provider := getConfig(req)
	if err := req.ParseForm(); err != nil {
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

func checkTokenHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
	config, _ := getConfig(req)
	if err := req.ParseForm(); err != nil {
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
	if err := idToken.Claims(&idTokenClaims); err != nil {
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

var TestHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "This is test.")
})

func main() {
	r := http.NewServeMux()

	// 認証で保護したいページ。ログインしていなければKeycloakのOpenID Connect認証ページに飛ばす
	r.Handle("/api", checkTokenHandler(TestHandler))

	// OpenID Connectの認証が終わった時に呼ばれるハンドラ
	// もろもろトークンを取り出したりした後に、クッキーを設定して元のページに飛ばす
	r.HandleFunc(OIDC_CALLBACK_PATH, SetTokenHandler)

	// log.Println(http.ListenAndServe(":8080", r))  // リクエストをロギングしなくてもよい場合
	log.Println(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r)))
}
