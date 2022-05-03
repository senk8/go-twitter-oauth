package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "github.com/sivchari/gotwtr"
    "github.com/senk8/go_twitter_oauth/oauth2"
)

const (
    authzEndpoint = "https://twitter.com/i/oauth2/authorize"
    tokenEndpoint = "https://api.twitter.com/2/oauth2/token"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal(err)
    }
    config := &oauth2.Config{
        ClientID:      os.Getenv("CLIENT_ID"),
        ClientSecret:  os.Getenv("CLIENT_SECRET"),
        RedirectURI:   os.Getenv("REDIRECT_URI"),
        Scopes:        []string{"tweet.read", "users.read"},
        AuthzEndpoint: authzEndpoint,
        TokenEndPoint: tokenEndpoint,
    }
    ctx := context.Background()
    oauth := oauth2.New(config)
    tokenResponse, err := oauth.ExecFlow(ctx)
    if err != nil {
        log.Fatal(err)
    }
    client := gotwtr.New(tokenResponse.AccessToken)
    tr, err := client.RetrieveSingleTweet(ctx, "1493108554618015752")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(*tr.Tweet)
}
