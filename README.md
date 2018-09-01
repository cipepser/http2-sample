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

### オレオレ証明書を登録する

上記のように`bad certificate`で怒られてしまったのでオレオレ証明書をroot CAとしてクライアントに登録する。

```go
caCert, err := ioutil.ReadFile("server.crt")
if err != nil {
  log.Fatalf("Reading server certificate: %s", err)
}
caCertPool := x509.NewCertPool()
caCertPool.AppendCertsFromPEM(caCert)

tlsConfig := &tls.Config{
  RootCAs: caCertPool,
}
```

#### 別案

試していないが、以下で無視してしまう方法もあるらしい。

```go
&tls.Config{
    InsecureSkipVerify: true,
}
```

#### 実行

これで実行するとうまくいく。

```sh
❯ go run client.go
Got response 200: HTTP/2.0 Hello
```

※`http.Get(url)`が`client.Get(url)`になっていることにも注意。少しハマった。

サーバ側のログ。

```sh
2018/09/01 14:57:12 Got connection: HTTP/2.0
```

#### HTTP/1.1で接続

```sh
❯ go run client.go -version 1
Got response 200: HTTP/1.1 Hello
```

サーバ側のログ。

```sh
2018/09/01 14:58:14 Got connection: HTTP/1.1
```

## Server Push

HTTP/1.1で接続する。

```Sh
❯ go run client.go -version 1
Got response 200: HTTP/1.1 Hello
```

結果。

```sh
2018/09/02 00:38:28 Got connection: HTTP/1.1
2018/09/02 00:38:28 Handling 1st
2018/09/02 00:38:28 Can\'t push to client
```

HTTP/2で接続する

```sh
❯ go run client.go -version 2
Got response 200: HTTP/2.0 Hello
```

結果。うまくいかない。

```sh
2018/09/02 00:39:16 Got connection: HTTP/2.0
2018/09/02 00:39:16 Handling 1st
2018/09/02 00:39:16 Failed push: feature not supported
```

現時点では、go client上、pushがdiabledになっているらしい。


ブラウザ(chrome)で`https://localhost:8080`へアクセスしてみて、ログを見る。

```sh
2018/09/02 00:47:46 http2: server: error reading preface from client [::1]:55036: read tcp [::1]:8080->[::1]:55036: read: connection reset by peer
2018/09/02 00:47:53 Got connection: HTTP/2.0
2018/09/02 00:47:53 Handling 1st
2018/09/02 00:47:53 Got connection: HTTP/2.0
2018/09/02 00:47:53 Handling 2nd
2018/09/02 00:47:53 Got connection: HTTP/2.0
2018/09/02 00:47:53 Handling 1st
2018/09/02 00:47:53 Got connection: HTTP/2.0
2018/09/02 00:47:53 Handling 2nd
```

## References
* [HTTP/2 Adventure in the Go World](https://posener.github.io/http2/)
* [opensslコマンドで証明書情報を確認したい。](https://jp.globalsign.com/support/faq/07.html)
* [\[Ssl\] Golang - 自己署名証明書付きTLS](https://code.i-harness.com/ja/q/159dbb3)