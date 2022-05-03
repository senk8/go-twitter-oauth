# go-twitter-oauth

goでツイッターAPIを叩く例

# セットアップ

```bash
$ git clone https://github.com/senk8/go-twitter-oauth.git
$ cd go-twitter-oauth
$ go mod tidy
```

# 実行方法

Twitter Developer Portalで 事前の準備を済ませた上で以下を行います。

1. main.goを実行します。

```bash
$ go run main.go
```

2. ブラウザで`http://localhost:3000/authz`にアクセスします。
3. ツイッターの認可画面が出てくるのでOKをします。
4. ブラウザが`http://localhost:3000/callback`にリダイレクトします。
5. APIを叩くのに成功します。

