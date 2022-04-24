# Appendix

開発中に利用した技術情報を備忘録として記載する。

## 自己証明局と自己証明書を使用したい

自己証明局と自己証明書は、以下の手順で用意した。

1. MacOS に自己認証局として [mkcert](https://github.com/FiloSottile/mkcert) を導入する。
   この操作により、MacOS のキーチェーンアクセス「システム」に `mkcert ***@***.local` が登録される。

   ```bash
   $ brew install mkcert

   # 自己認証局を作成して、信頼する認証局として OS に登録する。
   $ mkcert -install

   # ルート証明書が保存された位置を確認する。
   $ mkcert -CAROOT
   /Users/***/Library/Application Support/mkcert

   $ ls -l "/Users/***/Library/Application Support/mkcert"
   total 16
   -r--------  1 ***  staff  2484 12  5 10:57 rootCA-key.pem  # 秘密鍵　
   -rw-r--r--  1 ***  staff  1700 12  5 10:57 rootCA.pem      # SSL証明書（公開鍵）
   ```

1. 自己証明書を作成する。

   ```bash
   $ cd sso_with_keycloak/docker/nginx/

   # CSR（証明書署名要求）を作成する
   $ openssl req -new -key 【作成済みの秘密鍵】 -out ./server.csr -sha256
   ...
   Country Name (2 letter code) []:JP
   State or Province Name (full name) []:Tokyo
   Locality Name (eg, city) []:
   Organization Name (eg, company) []:
   Organizational Unit Name (eg, section) []:
   Common Name (eg, fully qualified host name) []: *.【ドメイン名】# ワイルドカード証明書にする
   Email Address []:

   Please enter the following 'extra' attributes
   to be sent with your certificate request
   A challenge password []:

   # SAN 設定を用意する
   $ echo "subjectAltName = DNS:app1.【ドメイン名】, DNS:auth.【ドメイン名】, IP:127.0.0.1" > san.txt

   # CRT（SSLサーバ証明書）を作成する
   $ openssl x509 -days 3650 -req -signkey 【作成済みの秘密鍵】 -in ./server.csr -out ./server.crt -extfile san.txt
   ```

1. MacOS のキーチェーンアクセスを開き、作成した server.crt をキーチェーンに追加する。  
   キーチェーンの種類は「ログイン」にする。

1. 登録された証明書を表示して、信頼 > SSL の値を「常に信頼」にする。

1. 登録したドメイン名でローカルからアクセスできるよう、hosts にエントリを追加する。

   ```bash
   $ vi /private/etc/hosts
   ...
   #### Added by me
   127.0.0.1       app1.【ドメイン名】
   127.0.0.1       auth.【ドメイン名】

   ```

参考：

- [開発環境を https 化する mkcert の仕組み](https://qiita.com/k_kind/items/b87777efa3d29dcc4467)
- [開発環境用にワイルドカード証明書を発行する](https://blog.tnantoka.com/posts/151/)

## MySQL 起動時に自動で複数の DB を作成したい

ググって見つけた docker-compose による Keycloak の起動を解説しているいくつかの記事では、docker-compose.yml に環境変数を指定して Keycloak 用の DB を作成していた。  
今回はアプリケーション用の DB も必要だが、Keycloak 用 DB だけ docker-compose.yaml に記載するのは分かりづらいと考えた。  
そのため、MySQL の Docker イメージで用意されているデータの初期化の機能を利用して、DB の初期化を実施した。

参考：[Docker で MySQL 起動時にデータの初期化を行う](https://qiita.com/moaikids/items/f7c0db2c98425094ef10)

## Nginx の設定ファイルからドメイン名を外出ししたい

環境によって変わるドメイン名は Nginx の設定ファイルにハードコーディングせず、外出ししたい。  
Nginx の Docker イメージには、既にこれを実現する機能がある模様。

参考：[【Docker】Nginx の conf で環境変数を使う](https://qiita.com/jungissei/items/2d6b40320b520f52b502)

## Go アプリケーションの変更がホットリロードされるようにしたい

開発時は変更を即確認できるよう、ホットリロードさせたい。  
Go 言語のホットリロードのライブラリの 1 つ、Air を使用して実現可能。

公式：[cosmtrek/air](https://github.com/cosmtrek/air)
参考：[Go+gin+Air 環境を Docker で構築](https://zenn.dev/hrs/articles/go-gin-air-docker)

## Go Rest API で、すべての API が必ずとおる動作を定義したい

ログイン状態が維持されているかすべての API で確認したい。  
ログイン状態を確認するためのミドルウェアを作成し、ルーティング部分で必ずミドルウェアを経由するよう実装して実現可能。

参考：[Go 言語で作る MiddleWare(by REST API)](https://selfnote.work/20200319/programming/golang-with-middleware/)

## React 開発環境を作成したい

以下を参考に、React の開発環境を準備した。

- [React の環境構築（セットアップ） | 独自に環境を構築](https://www.webdesignleaves.com/pr/jquery/react_basic_01.html)
- [webpack の基本的な使い方](https://www.webdesignleaves.com/pr/jquery/webpack_basic_01.html)

React の勉強だけなら `create-react-app` を利用した環境構築で足りそうだったが、React を動作させる環境の勉強も兼ねて、モジュールバンドラも自分で導入する方法を選択した。

開発環境の作成時に実行した手順は、以下のとおり。

1. プロジェクトフォルダを作成する。

   ```bash
   $ mkdir sso_with_keycloak/app/frontend
   $ cd sso_with_keycloak/app/frontend/
   ```

1. モジュールをインストールする。

   ```bash
   $ npm init -y
   $ npm install --save webpack webpack-cli  # Webpack インストール
   $ npm install --save-dev @babel/core @babel/preset-env @babel/preset-react babel-loader  # Babel インストール
   $ npm install --save react react-dom  # React インストール
   $ npm install --save webpack-dev-server # webpack-dev-server インストール
   $ npm install --save keycloak-js  # Keycloak 連携モジュール インストール
   $ npm install --save dotenv  # .env 読み込み用
   $ npm ls -depth=0  # インストール内容確認
   ├── @babel/core@7.16.5
   ├── @babel/preset-env@7.16.5
   ├── @babel/preset-react@7.16.5
   ├── babel-loader@8.2.3
   ├── keycloak-js@18.0.0
   ├── react-dom@17.0.2
   ├── react@17.0.2
   ├── webpack-cli@4.9.1
   ├── webpack-dev-server@4.7.1
   └── webpack@5.65.0
   ```

1. 参考サイトに従いディレクトリ構成, html, js, webpack.config.js を作成する。

   - [React の環境構築（セットアップ） | 独自に環境を構築](https://www.webdesignleaves.com/pr/jquery/react_basic_01.html)

1. 以下を用意して、docker-compose から起動する。  
   webpack-dev-server はバージョンにより起動コマンドを含めて設定方法がかなり異なる。  
   参考サイトそのままでは最新バージョンで動作しないため、導入したバージョンに合わせて設定方法を変更する必要がある。
   - React アプリケーション用 Dockerfile
   - React アプリケーションにプロキシさせる Nginx（app1.conf.template）
   - React アプリケーションを起動するための docker-compose.yaml
