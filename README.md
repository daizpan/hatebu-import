# hatebu-import

## Setting

https://www.hatena.ne.jp/my/config/auth

「OAuth ではてなのサービスにアクセスする」にアクセスし
`OAuth Consumer Key` `OAuth Consumer Secret` を取得する

```
export HATENA_OAUTH_KEY=XXXXX
export HATENA_OAUTH_SECRET=XXXXX
```
## Usage

### Export bookmark
https://b.hatena.ne.jp/-/my/config/data_management

### Import bookmark
```
go run main.go import -f bookmark.html
```

1. Go to https://www.hatena.com/oauth/authorize?oauth_token=xxxxx
2. Authorize the application
3. Enter verification code