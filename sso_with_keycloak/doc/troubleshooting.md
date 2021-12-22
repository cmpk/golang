# トラブルシューティング

開発中に遭遇したエラーや警告とその解消方法を備忘録として記載する。  

## port is already allocated

443 ポートが既に使用されているとして、nginx が起動しない。

```bash
$ docker-compose up
...
Error response from daemon: driver failed programming external connectivity on endpoint docker-nginx-1 (ead043eb2be705925b7791fe483d7f4901a25c1a4937a87a584b4a1b2d0cdb34): Bind for 0.0.0.0:443 failed: port is already allocated
```

443 ポートを利用しているプロセスは無さそうである。  

```bash
$ lsof -i:443 | grep LISTENING
# 何も表示されない
```

MacOS ごと再起動して、再度 `docker-compose up` を実行したら、nginx が起動した。  
原因不明だが、設定ファイルのミスや解決のための検証の過程でコンテナを強制終了させており、ポートが解放しきれていなかったのかもしれない。

## Error: Row size too large

Keycloak 用のテーブルについて、カラムサイズが大きすぎてテーブルを作成できない。  

```bash
$ docker-compose up
...
docker-keycloak-1  | 04:47:01,781 ERROR [org.keycloak.connections.jpa.updater.liquibase.conn.DefaultLiquibaseConnectionProvider] (ServerService Thread Pool -- 65) Change Set META-INF/jpa-changelog-1.9.1.xml::1.9.1::keycloak failed.  Error: Row size too large. The maximum row size for the used table type, not counting BLOBs, is 65535. This includes storage overhead, check the manual. You have to change some columns to TEXT or BLOBs [Failed SQL: ALTER TABLE keycloak.REALM MODIFY CERTIFICATE VARCHAR(4000)]
```

DB の CHARSET にとりあえず utf8mb4 を指定していたが、utf8でなければならない模様。  
`sso/docker/mysql/init/1_ddl.sql` で Keycloak 用 DB 作成時に utf8 を指定することで、エラー解消した。  

参考  
- [Server Installation | Relational Database Setup | Unicode Considerations for Databases](https://www.keycloak.org/docs/latest/server_installation/#unicode-considerations-for-databases)
- [Setting up Keycloak Standalone with MySQL Database](https://medium.com/@pratik.dandavate/setting-up-keycloak-standalone-with-mysql-database-7ebb614cc229)

## [Warning] Setting lower_case_table_names=2

テーブル名やデータベース名に使用された大文字小文字の解釈方法を明示的に指定していないために発生した警告。  

```bash
$ docker-compose up
...
docker-db-1        | 2021-11-13T04:07:59.573059Z 0 [Warning] Setting lower_case_table_names=2 because file system for /var/lib/mysql/ is case insensitive
```

ググってみると、`2` ではプラットフォームによっては動作しないことがある模様。    
予期せぬエラーを回避するため、`1`（保存時は小文字変換して解釈時は大文字小文字を区別しない）を `sso/docker/mysql/my.cnf` に指定。  
警告は発生しなくなった。

## A deprecated TLS version

デフォルトでは非推奨の TLS バージョンが有効になっている模様。  

```bash
$ docker-compose up
...
docker-db-1        | 2021-11-13T07:38:23.930952Z 0 [Warning] A deprecated TLS version TLSv1 is enabled. Please use TLSv1.2 or higher.
docker-db-1        | 2021-11-13T07:38:23.930955Z 0 [Warning] A deprecated TLS version TLSv1.1 is enabled. Please use TLSv1.2 or higher.
```

`sso/docker/mysql/my.cnf` に構築時点で非推奨でないバージョンを `tls_version` を指定して、警告は発生しなくなった。

## go.mod exists but should not

### MacOS

GOPATH に `go.mod` を保存した状態で `go mod tidy` を実行すると、エラーが発生する。

```bash
$ cd sso/app/backend/
$ export GOPATH=`pwd`
$ go mod tidy        
$GOPATH/go.mod exists but should not
```

GOPATH に別のディレクトリを指定すると、コマンドが実行して、その別のディレクトリにモジュールがインストールされた。

### Dockerbuild

Dockerfile に GOPATH を指定しない場合に、エラーが発生。  
[ベースイメージ](https://github.com/docker-library/golang/blob/d5ee0588aaa4a7be9bba3d1cb4b1abe0323b6442/1.17/alpine3.14/Dockerfile)で GOPATH と WORKDIR が同じパスになっているために発生している模様。  

GOPATH を空にするコード `ENV GOPATH=` を Dockerfile に記載して解消された。

## ブラウザからアクセスするとリダイレクトがループする

Nginx の設定が以下の状態で https://【ドメイン名】/auth/ にアクセスするとリダイレクトループに陥った。

```default.conf
server {
  listen 443 ssl default_server;
  ...
  location /auth/ {
    proxy_pass http://keycloak:8080/;
  }
  ...
}
```

書き方のお作法が間違っていたらしく、`proxy_pass` に指定した URL の末尾のスラッシュ(/)を削除して解消した。  

参考：[リダイレクトループにつながるNginxリバースプロキシ](https://www.webdevqa.jp.net/ja/nginx/%E3%83%AA%E3%83%80%E3%82%A4%E3%83%AC%E3%82%AF%E3%83%88%E3%83%AB%E3%83%BC%E3%83%97%E3%81%AB%E3%81%A4%E3%81%AA%E3%81%8C%E3%82%8Bnginx%E3%83%AA%E3%83%90%E3%83%BC%E3%82%B9%E3%83%97%E3%83%AD%E3%82%AD%E3%82%B7/960309889/)

## クライアントに内部URLが暴露される

https://【ドメイン名】/app1 にアクセスして、Golang アプリケーションから Keycloak にリダイレクトされると、URL が http://keycloak:8080（内部URL）になってしまい、クライアントがその URL を解決できない。

以下のように、`proxy_redirect` によりプロキシー先から帰ってきたレスポンスの Location ヘッダーと Refresh ヘッダーを書き換えることで解消した。

```default.conf
location /app1/ {
  proxy_pass     http://backend:8080;
  proxy_redirect http://keycloak:8080/ https://${DOMAIN_NAME}/;
}
```

なお、最終的に本対応は実施していない（∵ 内部通信も HTTPS およびドメイン名で実施することにしたため）。

参考：[Module ngx_http_proxy_module | Directives | proxy_redirect](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_redirect)

## Get "https://【ドメイン名】/auth/realms/app1/.well-known/openid-configuration": dial tcp 127.0.0.1:443: connect: connection refused

公開しているドメイン名を使用して OpenID Connect プロバイダーを登録しようとした場合に、エラーが発生する。  

```golang
provider, err = oidc.NewProvider(context.Background(), "https://【ドメイン名】/auth/realms/app1")
if err != nil {
	panic(err)
}
```

内部からは対象ドメイン名が解決できない状態のために発生しているのでと考えている（確証なし）。  
内部アドレス `http://backend-1:8080/auth/realms/app1` に変更して、解消した。  

Redirect URLは、KeyCloak に登録した値と Go言語で指定する値が一致していれば、公開しているドメイン名でも内部ドメインでも名前解決できた。 

なお、最終的には docker-compose.yaml に以下を記載して、内部通信もいったん Nginx を経由させる処置をとった（∵ 内部通信も HTTPS およびドメイン名で実施することにしたため）。

```yaml
services:
  nginx:
    ...
    networks:
      default:
        aliases:
          - "${KEYCLOAK_DOMAIN_NAME}"
          - "${APP_DOMAIN_NAME}"
```

参考：[【docker】好きなドメイン名(別名)でコンテナ間通信したい【compose】](https://qiita.com/KeisukeKudo/items/0d11717faeb81e42ddf6)

## 自己証明書ではなく別の証明書で認識される

Google Chrome でサイトを表示すると、以下の証明書で認識されてしまう。

- サブジェクト名 通称：*.mytrafficmanagement.com
- 発行者 組織：Let's Encrypt
- 発行者 通称：R3

自己証明書は `*.hoge.huga.com` を指定して作成していたが、`huga.com` の証明書（中間証明書）を用意していなかったために発生した模様。  
自己証明書を `*.hoge-huga.com` とすることで解消。  

## 自己証明書をローカルに登録しても Google Chrome で `NET::ERR_CERT_COMMON_NAME_INVALID` が表示される

ワイルドカード証明書として自己証明書を作成する際に、アクセスするドメイン名と完全に一致する SAN（Subject Alternative Name）を指定していない場合に発生する。

CRT（SSLサーバ証明書）作成時に、以下のように全てのドメイン名を SAN として指定することで解消。

```
subjectAltName = DNS:app1.【ドメイン名】, DNS:auth.【ドメイン名】, IP:127.0.0.1
```

参考：[Google Chrome で自組織のCAで署名したSSL証明書のサイトにアクセスすると NET::ERR_CERT_COMMON_NAME_INVALID エラーメッセージが表示される (Windows Tips)](https://www.ipentec.com/document/windows-chrime-error-net-err-cert-common-name-invalid-using-ssl-certificate-signed-with-local-ca)

## Get "https://【ドメイン名】/auth/realms/app1/.well-known/openid-configuration": x509: certificate signed by unknown authority

認証局が署名していないドメインを利用して OpenID Connect プロバイダーを登録しようとした場合に、エラーが発生する。  

```golang
provider, err = oidc.NewProvider(context.Background(), "https://【ドメイン名】/auth/realms/app1")
if err != nil {
	panic(err)
}
```

自己認証局を立てて自己認証した証明書を作成した後、Dockerfile に以下を記載して、作成した証明書を Go 言語が動作する OS に認識させて解消した。

```
COPY ./ssl/server.crt /usr/local/share/ca-certificates/server.crt
RUN apt update && \
    apt install -y ca-certificates && \
    update-ca-certificates
```

参考：
- [【図解付き】開発用オレオレ認証局SSL通信(+dockerコンテナ対応) : 2021](https://qiita.com/kaku3/items/e06a02ae1068de5c0663)
- [x509: certificate signed by unknown authority の対応](https://qiita.com/reikkk/items/e81fe384ad83a8e8b845)
