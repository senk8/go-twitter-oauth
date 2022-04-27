package gotwtr_oauth

import (
	"context"
	"encoding/json"
	"github.com/sivchari/gotwtr"
	"log"
	"net/http"
)

type Client struct {
	client  *gotwtr.Client
	config  *Config
	session *PKCESession
	ch      chan string
}

type OAuth2 interface {
	Client() *gotwtr.Client
	GenerateAccessToken(ctx context.Context) (string, error)
	authHandler(w http.ResponseWriter, r *http.Request)
	callbackHandler(w http.ResponseWriter, r *http.Request)
}

var _ OAuth2 = (*Client)(nil)

const (
	LISTEN_ADDR = "127.0.0.1:3000"
)

func New(config *Config) *Client {
	return config.createClient()
}

func (c *Client) Client() *gotwtr.Client {
	return c.client
}

func (c *Client) GenerateAccessToken(ctx context.Context) (string, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", c.authHandler)
	mux.HandleFunc("/callback", c.callbackHandler)
	srv := &http.Server{
		Addr:    LISTEN_ADDR,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	accessToken := <-c.ch
	if err := srv.Shutdown(ctx); err != nil {
		return "", err
	}

	return accessToken, nil
}

func (c *Client) authHandler(w http.ResponseWriter, r *http.Request) {
	c.session = CreatePKCESession()
	authURL := c.config.buildAuthURL(c.session)
	w.Header().Set("Location", authURL)
	w.WriteHeader(http.StatusFound)
	return
}

func (c *Client) callbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if c.session.State != state {
		log.Println("invalid state")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	req := c.config.createTokenRequest(code, c.session.CodeVerifier)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var tokenResponse TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		log.Fatal(err)
	}
	c.ch <- tokenResponse.AccessToken
	return
}
