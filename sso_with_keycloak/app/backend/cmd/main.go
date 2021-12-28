//
// 現時点では以下のコードをほぼそのまま利用している。
// [OpenID Connectを使ったアプリケーションのテストのためにKeycloakを使ってみる](https://qiita.com/shibukawa/items/fd78d1ca6c23ce2fa8df)
//
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

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
			RedirectURL:  os.Getenv("APP_URL") + "/api/callback",
		}
	})
	return oauth2Config, provider
}

func main() {
	fmt.Fprintln(os.Stdout, "===== Golang main =====")

	// 認証で保護したいページ。ログインしていなければKeycloakのOpenID Connect認証ページに飛ばす
	http.HandleFunc("/api", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stdout, "===== Golang /api =====")

		// クッキーがない時はリダイレクト
		if _, err := req.Cookie("Authorization"); err != nil {
			config, _ := getConfig(req)
			url := config.AuthCodeURL(os.Getenv("STATE_STRING"))

			http.Redirect(w, req, url, http.StatusFound)
			return
		}
		io.WriteString(w, "login success")
	})

	// OpenID Connectの認証が終わった時に呼ばれるハンドラ
	// もろもろトークンを取り出したりした後に、クッキーを設定して元のページに飛ばす
	http.HandleFunc("/api/callback", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stdout, "===== Golang /api/callback =====")

		config, provider := getConfig(req)
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
		oidcConfig := &oidc.Config{
			ClientID: os.Getenv("CLIENT_ID"),
		}
		verifier := provider.Verifier(oidcConfig)
		idToken, err := verifier.Verify(context.Background(), rawIDToken)
		if err != nil {
			http.Error(w, "id token verify error : "+err.Error(), http.StatusInternalServerError)
			return
		}
		// IDトークンのクレームをとりあえずダンプ
		// アプリで必要なものはセッションストレージに入れておくと良いでしょう
		idTokenClaims := map[string]interface{}{}
		if err := idToken.Claims(&idTokenClaims); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("%#v", idTokenClaims)
		http.SetCookie(w, &http.Cookie{
			Name:  "Authorization",
			Value: "Bearer " + rawIDToken, // 行儀が悪いので真似しないねで
			Path:  "/api",
		})
		http.Redirect(w, req, "/api", http.StatusFound)
	})
	log.Println(http.ListenAndServe(":8080", nil))
}
