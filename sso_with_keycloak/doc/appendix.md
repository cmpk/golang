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
    $ cd sso/docker/nginx/

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
- [開発環境をhttps化するmkcertの仕組み](https://qiita.com/k_kind/items/b87777efa3d29dcc4467)
- [開発環境用にワイルドカード証明書を発行する](https://blog.tnantoka.com/posts/151/)

## MySQL 起動時に自動で複数の DB を作成したい

ググって見つけた docker-compose による Keycloak の起動を解説しているいくつかの記事では、docker-compose.yml に環境変数を指定して Keycloak 用の DB を作成していた。  
今回はアプリケーション用の DB も必要だが、Keycloak 用 DB だけ docker-compose.yaml に記載するのは分かりづらいと考えた。  
そのため、MySQL の Docker イメージで用意されているデータの初期化の機能を利用して、DB の初期化を実施した。  

参考：[Docker で MySQL 起動時にデータの初期化を行う](https://qiita.com/moaikids/items/f7c0db2c98425094ef10)

## Nginx の設定ファイルからドメイン名を外出ししたい

環境によって変わるドメイン名は Nginx の設定ファイルにハードコーディングせず、外出ししたい。  
Nginx の Docker イメージには、既にこれを実現する機能がある模様。

参考：[【Docker】Nginxのconfで環境変数を使う](https://qiita.com/jungissei/items/2d6b40320b520f52b502)

## Goアプリケーションの変更がホットリロードされるようにしたい

開発時は変更を即確認できるよう、ホットリロードさせたい。  
Go言語のホットリロードのライブラリの1つ、Air を使用して実現可能。  

公式：[cosmtrek/air](https://github.com/cosmtrek/air)
参考：[Go+gin+Air環境をDockerで構築](https://zenn.dev/hrs/articles/go-gin-air-docker)
