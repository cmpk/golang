//
// 現時点では以下のコードをほぼそのまま利用している。
// [OpenID Connectを使ったアプリケーションのテストのためにKeycloakを使ってみる](https://qiita.com/shibukawa/items/fd78d1ca6c23ce2fa8df)
//
package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"backend/authentication"
)

var TestHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "This is test.")
})

func main() {
	r := http.NewServeMux()

	// 認証で保護したいページ。ログインしていなければKeycloakのOpenID Connect認証ページに飛ばす
	r.Handle("/api", authentication.CheckTokenHandler(TestHandler))

	// OpenID Connectの認証が終わった時に呼ばれるハンドラ
	// もろもろトークンを取り出したりした後に、クッキーを設定して元のページに飛ばす
	r.HandleFunc(authentication.OIDC_CALLBACK_PATH, authentication.SetTokenHandler)

	// log.Println(http.ListenAndServe(":8080", r))  // リクエストをロギングしなくてもよい場合
	log.Println(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r)))
}
