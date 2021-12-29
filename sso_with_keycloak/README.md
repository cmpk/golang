# Go 言語で Keycloak と SSO

!!! 開発中です !!!

Go 言語で開発した２つのアプリケーションへのログインを、Keycloak による Single Sign On で実現する。

## 前提

- 言語そのものの学習のため、フレームワークは利用しない。

## システム構成

![システム構成図](./doc/system_structure.drawio.svg)

- 今時ユーザーからのアクセスは HTTPS でしょう！ということで、オレオレ証明書を使っている。
- Keycloak とアプリケーションの通信も、HTTPS とドメイン名を利用するため、Nginx を経由させる。
  - これにより、URL の変換処理を省くことができる。  
    （Keycloak に登録する URL や Go 言語から Keycloak に通信する際に認識される URL などの整合性を揃えるのが面倒）

## 動作確認環境

- macOS Monterery 12.0.1
- Google Chrome
- Docker Engine - Community 20.10.10
- Docker Compose version v2.1.1
- Docker コンテナで動作させるもの
  - Nginx
  - MySQL
  - Keycloak
  - Go - React アプリケーション 1
  - Go - React アプリケーション 2
- 証明書は自己署名証明書（オレオレ証明書）

## 開発環境

- macOS Monterery 12.0.1
- Visual Studio Code 1.59  
  開発に使用した拡張機能は `.vscode/extensions.json` に記録している。
  - Formatter
    - NGINX Configuration 0.7.2
    - nginx-format 0.0.6
    - Prettier - Code formatter 9.0.0
    - ESLint v2.2.1
  - 言語サーバー
    - gopls
- Go 言語
- Node.js
- yarn

## ディレクトリ構成

T.B.D

## 開始手順

1. 証明書を準備する。

   証明書を入手し、入手した証明書を `sso_with_keycloak/docker/ssl/`、`sso_with_keycloak/app/backend/ssl/` 配下に配置する。  
   ファイル名は `server.crt` および `server.csr` とすること。

   「[自己証明局と自己証明書を使用したい](./doc/appendix.md#自己証明局と自己証明書を使用したい)」の手順により自己証明書を用いた動作確認が可能である。  
   ※ インターネットに公開する環境では自己証明書を利用しないこと。

1. 環境情報を用意する。

   `sso_with_keycloak/docker/sample.env` をコピーして `sso_with_keycloak/docker/.env` を作成し、必要に応じて内容を書き換える。

1. Docker コンテナを起動する。

   ```bash
   $ cd sso_with_keycloak/docker/

   # ログの確認のためフォアグラウンドで起動
   $ docker-composer up
   ...
   keycloak_1  | 10:16:47,369 INFO  [org.jboss.as] (Controller Boot Thread) WFLYSRV0051: Admin console listening on http://127.0.0.1:9990
   ```

1. Keycloak にアクセスし、Console にログインする。

   - URL : https://{Keycloak のドメイン名}/auth
   - ユーザー名 : admin
   - パスワード : admin

1. 名前が `app1` の Realm を作成する。

1. アプリケーション用 Client を登録する。  
   サイドメニュー「Clients」から以下の値をもつ Client を作成する。

   - Client ID : application
   - Root URL : https://{アプリケーションのドメイン名}/api

1. 作成した Client の「Settings」タブで以下を設定して、保存する。

   - Access Type : confidential

1. アプリケーション用ユーザーを登録する。  
   サイドメニュー「Users」から任意の User を作成する。  
   パスワードは「Credentials」から作成すること。

1. 作成した Client の「Credentials」から確認した Secret の値を `.env` の `CLIENT_SECRET` に記載する。

1. アプリケーションを再起動する。  
   先ほど実行した `docker-compose` を Ctrl+C で停止して、再度起動する。

1. ブラウザからアプリケーションにアクセスし、Keycloak のログイン画面が表示されることを確認する。
   - URL : https://{アプリケーションのドメイン名}/api

ログインするためには、あらかじめ Keycloak 側で User を登録している必要がある。

## Tips

### DB データの作り直し

1. コンテナが起動していないことを確認する。

1. `docker/mysql/data` 配下の .\* 以外のファイルを削除する。

1. コンテナを起動する。

## お世話になったサイト

- [OSS なシングルサインオンサービス Keycloak を docker で立ち上げる](https://qiita.com/myoshimi/items/7e9f1de7373427233880)
- [OpenID Connect を使ったアプリケーションのテストのために Keycloak を使ってみる](https://qiita.com/shibukawa/items/fd78d1ca6c23ce2fa8df)
