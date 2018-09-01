# http2-sample

[HTTP/2 Adventure in the Go World](https://posener.github.io/http2/)を写経したメモ。

## 秘密鍵と証明書を作る

HTTP/2ではTLSが必要なので以下で秘密鍵と証明書を作る。

```sh
❯ openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt
```

以下で設定した。

```
server.crtCountry Name (2 letter code) []:JP
State or Province Name (full name) []:Tokyo
Locality Name (eg, city) []:Tokyo
Organization Name (eg, company)[]:Individual
Organizational Unit Name (eg, section) []:Individual
Common Name (eg, fully qualified host name) []:localhost
Email Address []:no address
```

一応確認しておく。

```
❯ openssl x509 -text -noout -in server.crt
Certificate:
（略）
        Issuer: C=JP, ST=Tokyo, L=Tokyo, O=:Individual, OU=Individual, CN=localhost/emailAddress=no address
```


## `server.go`

記事の通り実装。

## `client.go`

記事の通り実装。
(単純に`http.Get`するだけのバージョン)

### 実行

クライアント側

```sh
❯ go run client.go
2018/09/01 13:58:05 Get https://localhost:8080: x509: certificate signed by unknown authority
exit status 1
```

サーバ側

```sh
❯ go run server.go
2018/09/01 13:57:37 Serving on https://0.0.0.0:8080
2018/09/01 13:58:05 http: TLS handshake error from [::1]:52629: remote error: tls: bad certificate
```


## References
* [HTTP/2 Adventure in the Go World](https://posener.github.io/http2/)
* [opensslコマンドで証明書情報を確認したい。](https://jp.globalsign.com/support/faq/07.html)